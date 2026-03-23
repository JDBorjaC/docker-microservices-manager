import { useState } from 'react';
import type { Service, ServiceCreateForm, ServiceUpdateForm } from '../models/msm_models'
import MonitorBackdrop from '../components/monitor';
import { useNavigate } from 'react-router-dom';

const backendUrl = "http://localhost:8080/microservices";

export default function ServiceEditor(){

    const navi = useNavigate();

    const [editable] = useState(():Service => {
        //cargar editService si es que está
        const e = localStorage.getItem("editService");
        console.log(e)
        return e ? JSON.parse(e):{
            container_id:"",
            created_at:"",
            description:"vacío",
            id:"-1337",
            code:"",
            image:"",
            language:"",
            name:"",
            status:""
        }
    });

    const [name, setName] = useState(editable.name || "");
    const [code, setCode] = useState(editable.code || "");
    const [language, setLanguage] = useState(editable.language || "rust");
    const [description, setDescription] = useState(editable.description || "");

    const [loading, setLoading] = useState(false);

    const editOrCreateService = async () => {
        setLoading(true);

        const isEdit = editable && editable.id && editable.id !== "-1337";
        console.log("editing?: ",isEdit)

        try {
            if (isEdit) {
                const body: ServiceUpdateForm = {
                    code:code,
                    description:description,
                };

                const response = await fetch(backendUrl+"/"+editable.id, {
                    method: "PUT",
                    body: JSON.stringify(body),
                });

                if (!response.ok) throw new Error(`PUT failed: ${response.statusText}`);
            } else {
                const body: ServiceCreateForm = {
                    name:name,
                    code:code,
                    language:language,
                    description:description,
                };

                const response = await fetch(backendUrl, {
                    method: "POST",
                    body: JSON.stringify(body),
                });

                if (!response.ok) throw new Error(`POST failed: ${response.statusText}`);
            }

            localStorage.removeItem("editService");
        } catch (error) {
            console.error("Error creando/actualizando servicio:", error);
        } finally {
            setLoading(false);
            navi("/admin");
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

                            <h1> \\ EDITAR O CREAR MICROSERVICIO </h1>

                            <p>Seleccionar lenguaje de programación, editar el código, enviar!!!</p>
                            <p>Recomendación: Pegar el código desde algún Sandbox en linea del lenguaje seleccionado.</p>
                            <p>Advertencia: Para que su microservicio funcione, tiene que definir una función 'microservice()', que actuará como la función principal que será ejecutada. </p>

                            <div className="service-form">
                                <input
                                    className="code-input"
                                    type="text"
                                    name="name"
                                    placeholder="Nombre del microservicio"
                                    value={name}
                                    onChange={(e) => setName(e.target.value)}
                                    required
                                    disabled={loading}
                                />
                                <textarea
                                    className="code-block"
                                    name="code"
                                    placeholder="¡Copiar y pegar código aquí!"
                                    value={code}
                                    onChange={(e) => setCode(e.target.value)}
                                    required
                                    disabled={loading}
                                />
                                <div>
                                    <button className='monitor-button' onClick={() => {editOrCreateService()}} disabled={loading}>
                                        {loading ? "CARGANDO..." : "CONFIRMAR"}
                                    </button>
                                    <select
                                        className="drop-down-menu"
                                        name="lang-select-menu"
                                        value={language}
                                        onChange={(e) => setLanguage(e.target.value)}
                                        required
                                        disabled={loading}
                                    >
                                        <option value="axum">Axum (Rust)</option>
                                        <option value="flask">Flash (Python)</option>
                                        <option value="express">Express (JavaScript)</option>
                                        <option value="gin">Gin (GO)</option>
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