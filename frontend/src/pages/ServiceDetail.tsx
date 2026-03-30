import { useEffect, useState } from 'react';
import type { Service, ServiceUpdateForm } from '../models/msm_models'
import MonitorBackdrop from '../components/monitor';
import { useNavigate, useParams } from 'react-router-dom';

const backendUrl = "http://localhost:8080/microservices";

export function ServiceDetail() {

    const navi = useNavigate();
    const { serviceId } = useParams();

    const [editable, setEditable] = useState((): Service => {
        return {
            "container_id": "",
            "created_at": "",
            "description": "vacío",
            "id": "1",
            "code": "",
            "image": "",
            "language": "",
            "name": "",
            "status": ""
        }
    });

    const fetchDeets = async () => {
        const deetsReq = await fetch(backendUrl + "/" + serviceId)
        //espero que el backend me mande un Service con código definido
        const service: Service = await deetsReq.json();
        setEditable(service);

    }

    //ON PAGELOAD: fetch service info and subscribe to status updates
    useEffect((): () => void => {
        // Referencias mutables para poder limpiarlas
        let es: EventSource | null = null;
        let retryTimeout: ReturnType<typeof setTimeout> | null = null;

        // Estado del backoff entre intentos de reconexion
        let retryDelay = 1000;
        const MAX_DELAY = 30_000;

        // flag que impide reconectar si el componente ya fue desmontado.
        let destroyed = false;

        const connect = () => {
            if (destroyed) return;

            es = new EventSource("http://localhost:8080/microservices/status/events");

            // Resetear backoff
            es.onopen = () => {
                retryDelay = 1000;
                fetchDeets();
            };

            // Actualizar estado del microservicio
            es.addEventListener("status_update", (e: MessageEvent) => {
                const updatedMs: Service = JSON.parse(e.data);
                if (String(updatedMs.id) === String(serviceId)) {
                    if (updatedMs.status === "removed") {
                        alert("Este microservicio ha sido eliminado desde otra sesión o por Docker.");
                        navi("/admin");
                    } else {
                        setEditable(updatedMs);
                    }
                }
            });

            es.onerror = () => {
                // Cierre explicito para tomar control total de la reconexion
                es?.close();
                es = null;

                if (destroyed) return;

                // Jitter del ±20% porque... cosas
                const jitter = retryDelay * 0.2 * (Math.random() * 2 - 1);
                const delay = retryDelay + jitter;

                retryTimeout = setTimeout(() => {
                    // Backoff exponencial: cada fallo duplica el tiempo de espera hasta el tope de 30s.
                    retryDelay = Math.min(retryDelay * 2, MAX_DELAY);
                    connect();
                }, delay);
            };
        };

        connect();

        // Limpieza del componente
        return () => {
            destroyed = true;
            es?.close();
            if (retryTimeout) clearTimeout(retryTimeout);
        };
    }, [serviceId]);

    const [loading, setLoading] = useState(false);
    const [logs, setLogs] = useState<string[]>([]);

    const editService = async () => {
        setLoading(true);
        console.log("Editando microservicio...");
        var success: boolean = false;
        try {
            const body: ServiceUpdateForm = {
                code: editable.code || "",
                description: "",
            };

            const response = await fetch(backendUrl + "/" + editable.id, {
                method: "PUT",
                body: JSON.stringify(body),
            });

            if (!response.ok) throw new Error(`PUT failed: ${response.statusText}`);

            localStorage.removeItem("editService");
            success = true;
        } catch (error) {
            console.error("Error creando/actualizando servicio:", error);
        } finally {
            setLoading(false);
            if (success) {
                navi("/admin");
            }
        }
    }

    const startService = async () => {
        setLoading(true);
        try {
            await fetch(backendUrl + "/start/" + editable.id, { method: "PATCH" });
            // Status will update via SSE
        } catch (error) {
            console.error("Error starting service:", error);
        } finally {
            setLoading(false);
        }
    }

    const fetchLogs = async () => {
        setLoading(true);
        setLogs(["[FETCHING LOGS...]"]);
        try {
            const resp = await fetch(backendUrl + "/logs/" + editable.id);
            if (!resp.ok) {
                const errData = await resp.json().catch(() => ({}));
                setLogs([`[ERROR] ${errData.error || resp.statusText}`]);
                return;
            }
            const data = await resp.text();
            setLogs(data.split('\n'));
        } catch (error) {
            console.error("Error fetching logs:", error);
            setLogs([`[NETWORK ERROR] ${error}`]);
        } finally {
            setLoading(false);
        }
    }


    const shutService = async () => {
        if (editable.status === "running") {
            setLoading(true);
            try {
                const resp = await fetch(backendUrl + "/stop/" + editable.id, {
                    method: "PATCH"
                });
                if (resp.ok) {
                    setLoading(false);
                } else {
                    throw new Error("Algo salió mal apagando el contenedor:\n");
                }
            } catch (err) {
                setLoading(false);
                console.error("Error apagando el Contenedor: ", err);
            } finally {
                fetchDeets();
            }
        }
    }

    return (
        <div className="micro-editor">
            {/* EDITOR DE MICROSERVICIO */}
            <MonitorBackdrop>
                <div className="monitor-bezel">
                    <div className="monitor-screen">
                        <div className='monitor-scanlines'>
                            <div className="monitor-content">

                                <h1> \\ ESTADO DE MICROSERVICIO </h1>

                                <p>En esta página se puede confirmar el estado del Microservicio.</p>
                                <p>Recordatorio: Para que su microservicio funcione, tiene que definir una función 'microservice()', que actuará como la función principal que será ejecutada. </p>

                                <div className="service-form">
                                    <input
                                        className="code-input"
                                        type="text"
                                        name="name"
                                        placeholder="Nombre del microservicio"
                                        value={editable.name}
                                        disabled={true}
                                    />
                                    <input
                                        className="code-input"
                                        type="text"
                                        name="url"
                                        placeholder="Enlace al microservicio"
                                        value={"ENLACE AL MICROSERVICIO: http://localhost/services/" + editable.name}
                                        disabled={true}
                                    />
                                    <textarea
                                        className="code-block"
                                        name="code"
                                        placeholder="¡Copiar y pegar código aquí para actualizar el microservicio!"
                                        value={editable.code}
                                        onChange={(e) => setEditable({ ...editable, code: e.target.value })}
                                        required
                                        disabled={loading}
                                    />
                                    <textarea
                                        className="code-block"
                                        name="logs"
                                        placeholder="Para ver los logs, haga clic en el botón OBTENER LOGS"
                                        value={logs.join("\n")}
                                        disabled={true}
                                    />
                                    <div>
                                        <button className='monitor-button' onClick={() => { editService() }} disabled={loading}>
                                            {loading ? "..." : "MODIFICAR"}
                                        </button>
                                        <button
                                            className={`monitor-button status-${editable.status}`}
                                            onClick={() => editable.status === "running" ? shutService() : startService()}
                                            disabled={loading}
                                        >
                                            {loading ? "..." : (editable.status === "running" ? "STOP" : "START")}
                                        </button>
                                        <button className='monitor-button' onClick={() => fetchLogs()} disabled={loading}>
                                            OBTENER LOGS
                                        </button>
                                        <select
                                            className="drop-down-menu"
                                            name="lang-select-menu"
                                            value={editable.language}
                                            onChange={(e) => setEditable({ ...editable, language: e.target.value })}
                                            required
                                            disabled={true}
                                        >
                                            <option value="axum">axum (Rust)</option>
                                            <option value="flask">flask (Python)</option>
                                            <option value="express">express (JavaScript)</option>
                                            <option value="gin">gin (Go)</option>
                                        </select>
                                    </div>
                                </div>

                            </div>
                        </div>
                    </div>
                </div>
            </MonitorBackdrop>
        </div>
    )
}