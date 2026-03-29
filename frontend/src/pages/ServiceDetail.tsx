import { useEffect, useRef, useState } from 'react';
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
        fetchDeets()

        const es = new EventSource("http://localhost:8080/microservices/status/events");
        es.addEventListener("status_update", (e: MessageEvent) => {
            const updatedMs: Service = JSON.parse(e.data);
            if (String(updatedMs.id) === String(serviceId)) {
                setEditable(updatedMs);
            }
        });

        return () => {
            es.close();
            if (esRef.current) {
                esRef.current.close();
            }
        }
    }, [serviceId]);

    const [loading, setLoading] = useState(false);
    const [logs, setLogs] = useState<string[]>([]);
    const esRef = useRef<EventSource | null>(null);

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

    const toggleLogs = () => {
        if (esRef.current) {
            esRef.current.close();
            esRef.current = null;
            setLogs(prev => [...prev, "[LOG STREAM DISCONNECTED]"]);
            return;
        }

        setLogs([]);
        const es = new EventSource(backendUrl + "/stream/" + editable.id);
        esRef.current = es;

        //Se añade un listener por cada tipo de mensaje que se espera del backend
        es.addEventListener("info", (e: MessageEvent) => {
            setLogs(prev => [...prev, "[INFO] " + e.data]);
        });
        es.addEventListener("log", (e: MessageEvent) => {
            setLogs(prev => [...prev, e.data]);
        });
        es.addEventListener("error", (e: MessageEvent) => {
            setLogs(prev => [...prev, "[ERROR] " + e.data]);
        });

        //abortar en caso de done o de error
        es.addEventListener("done", () => {
            es.close();
            esRef.current = null;
            setLogs(prev => [...prev, "[SESSION TERMINATED]"]);
        });
        es.onerror = (e) => {
            console.error("DEBUG SSE Error Event:", e);
            if (es.readyState === EventSource.CLOSED) {
                console.log("SSE Connection closed by server or network.");
            } else if (es.readyState === EventSource.CONNECTING) {
                console.log("SSE Attempting to reconnect...");
            }
            esRef.current = null;
            setLogs(prev => [...prev, `[SESSION ERROR - ReadyState: ${es.readyState}]`]);
        };
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
                                        placeholder="Para encender los Logs, unda click en el boton de STATUS, abajo de este recuadro"
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
                                        <button className='monitor-button' onClick={() => toggleLogs()} disabled={loading}>
                                            {esRef.current ? "DETENER LOGS" : "VER LOGS"}
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