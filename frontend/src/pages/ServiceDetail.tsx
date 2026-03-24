import { useEffect, useState } from 'react';
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
        const service:Service = await deetsReq.json();  // toca estar atento a cambios del modelo
        setEditable(service);

    }

    //ON PAGELOAD: fetch service info
    useEffect(():void => {
        fetchDeets()
    }, []);

    const [loading, setLoading] = useState(false);

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
            } else if (editable.status === "created" || editable.status === "stopped") {
                const resp = await fetch(backendUrl+"/start/"+editable.id, {
                    method:"PATCH"
                });
                if (resp.ok) {
                    fetchDeets();
                } else {
                    throw new Error("Algo salió mal iniciando el contenedor:\n");
                }
            }
        } catch (error) {
            console.error("Error: ", error);
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
                                    placeholder="¡Copiar y pegar código aquí!"
                                    value={editable.code}
                                    onChange={(e) => setEditable({...editable, code:e.target.value})}
                                    required
                                    disabled={loading}
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