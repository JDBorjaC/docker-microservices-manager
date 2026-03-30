import { useState } from 'react';
import type { ServiceCreateForm } from '../models/msm_models'
import { MonitorBackdrop } from '../components/monitor';
import { useNavigate } from 'react-router-dom';
import { ExpandableMessage } from '../components/acc_message';
import { BackButton } from '../components/back_btn';

const backendUrl = "http://localhost:8080/microservices";

export default function ServiceCreator(){

    const navi = useNavigate();

    const [name, setName] = useState("");
    const [code, setCode] = useState("");
    const [language, setLanguage] = useState("flask");

    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const createService = async () => {
        setLoading(true);
        setError(null);

        var success:boolean = false

        try {
            const body: ServiceCreateForm = {
                name:name,
                code:code,
                language:language,
                description:"",
            };
            const response = await fetch(backendUrl, {
                method: "POST",
                body: JSON.stringify(body),
            });
            
            if (!response.ok){ 
                throw new Error(`POST failed: ${response.statusText}`) };
            
            localStorage.removeItem("editService");
            success = true

        } catch (error) {
            console.error("Error creando servicio:", error);
            setError("No se pudo crear el microservicio. Intenta de nuevo.");
        } finally {
            setLoading(false);
            if(success){
                navi("/admin")
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

                            <div className='title-button-div'>
                                <h1> \\ CREAR MICROSERVICIO  </h1>
                                
                                <BackButton navi={navi} retTo='/admin/'/>
                            </div>

                            <p>Seleccionar lenguaje de programación, editar el código, enviar!!!</p>
                            <p>Recomendación: Pegar el código desde algún Sandbox en linea del lenguaje seleccionado.</p>
                            <ExpandableMessage title='ADVERTENCIA'>
                                <p>- Para escribir el nombre del microservicio, asegúrese de usar sólo minúsculas, sin caracteres especiales y sin espacio</p>
                                <p>- Para todos los lenguajes, es necesario colocar el servicio a escuchar en todas las interfaces.</p>
                                <p>- Por otro lado, cada lenguaje tiene que exponer un puerto especifico de acuerdo con la configuración por defecto del framework. Puede usar el código a continuación para cada lenguaje:</p>
                                <ExpandableMessage title='FLASK'>
                                    <pre>
                                        if __name__ == '__main__':<br/>    app.run(host='0.0.0.0', port=5000)
                                    </pre>
                                </ExpandableMessage>
                                <ExpandableMessage title='EXPRESS'>
                                    <pre>{"app.listen(PORT, '0.0.0.0', () => {\n    console.log('Microservicio escuchando correctamente en el puerto ${PORT}');\n});"}</pre>
                                </ExpandableMessage>
                            </ExpandableMessage>

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
                                    <button className='monitor-button' onClick={() => {createService()}} disabled={loading}>
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
                                        <option value="flask">Flask (Python)</option>
                                        <option value="express">Express (JavaScript)</option>
                                    </select>
                                </div>
                            </div>

                            {error && (
                                <div className="error-message"> ⚠ ERROR: {error}
                                </div>
                            )}

                        </div>
                    </div>
                </div>
            </div>
        </MonitorBackdrop>
        </div>
    )
}