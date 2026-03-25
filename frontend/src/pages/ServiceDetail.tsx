import { useEffect, useRef, useState } from 'react';
import type { Service, ServiceUpdateForm } from '../models/msm_models'
import MonitorBackdrop from '../components/monitor';
import { useNavigate, useParams } from 'react-router-dom';

const backendUrl = "http://localhost:8080/microservices";

export function ServiceDetail(){

    const navi = useNavigate();
    const { serviceId } = useParams();

    const [editable, setEditable] = useState(():Service => {
        return {
            "container_id":"",
            "created_at":"",
            "description":"vacío",
            "id":"1",
            "code":"",
            "image":"",
            "language":"",
            "name":"",
            "status":""
        }
    });

    const fetchDeets = async () => {
        const deetsReq = await fetch(backendUrl+"/"+serviceId)
        //espero que el backend me mande un Service con código definido
        const service:Service = await deetsReq.json(); 
        setEditable(service);

    }

    //ON PAGELOAD: fetch service info
    useEffect(():void => {
        fetchDeets()
    }, []);

    const [loading, setLoading] = useState(false);
    const [logs, setLogs] = useState<string[]>([]);
    const esRef = useRef<EventSource | null>(null);

    const editService = async () => {
        setLoading(true);
        console.log("Editando microservicio...");
        var success:boolean = false;
        try {
            const body: ServiceUpdateForm = {
                code:editable.code || "",
                description:"",
            };

            const response = await fetch(backendUrl+"/"+editable.id, {
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
            if(success){
                navi("/admin");
            }
        }
    }

    const bootService = async () => {
        //Prender/apagar dependiendo del status. Al prender, enganchar el EventSource
        setLoading(true);
        try{
            if (editable.status === "running"){
                const resp = await fetch(backendUrl+"/stop/"+editable.id, {
                    method:"PATCH"
                });
                if (resp.ok) {
                    fetchDeets();
                } else {
                    throw new Error("Algo salió mal apagando el contenedor:\n");
                }
            } else if (editable.status === "created" || editable.status === "stopped" || editable.status === "failed") {
                //Cuando se encienda el microservicio, toca enganchar el SSE stream a esRef
                if (esRef.current) {
                    esRef.current.close();
                    esRef.current = null;
                }

                const es = new EventSource(backendUrl+"/stream/"+editable.id);
                esRef.current = es;

                es.onmessage = (event) => {
                     console.log("STREAM EVENT: ", event.data);
                     setLogs((prevLogs) => [...prevLogs, event.data]);
                }

                //por si el server manda una señal de que se acabó el stream
                es.addEventListener("done", () => {
                  es.close();
                  esRef.current = null;
                  fetchDeets();
                });
                
                es.onerror = (err) => {
                  console.error("Stream error:", err);
                  setLogs((prevLogs) => [...prevLogs, "ERROR"]);
                  es.close();
                  esRef.current = null;
                  fetchDeets();
                };

            }
        } catch (error) {
            console.error("Error: ", error);
        } finally {
            setLoading(false);
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
                                <textarea
                                    className="code-block"
                                    name="code"
                                    placeholder="¡Copiar y pegar código aquí para actualizar el microservicio!"
                                    value={editable.code}
                                    onChange={(e) => setEditable({...editable, code:e.target.value})}
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
                                    <button className='monitor-button' onClick={() => {editService()}} disabled={loading}>
                                        {loading ? "CARGANDO..." : "MANDAR CAMBIOS"}
                                    </button>
                                    <button className='monitor-button' onClick={() => {bootService()}} disabled={loading}>
                                        {loading ? "..." : "STATUS: "+editable.status}
                                    </button>
                                    <select
                                        className="drop-down-menu"
                                        name="lang-select-menu"
                                        value={editable.language}
                                        onChange={(e) => setEditable({...editable, language:e.target.value})}
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