export default function ServiceEditor(){
    
    return (
        <div className="micro-editor">
        {/* EDITOR DE MICROSERVICIO */}
        <div className="monitor-frame">
            <div className="monitor-bezel">
                <div className="monitor-screen">
                    <div className='monitor-scanlines'>
                        <div className="monitor-content">

                            <h1>
                                \\ EDITAR O CREAR MICROSERVICIO
                            </h1>

                            <p>
                                Seleccionar lenguaje de programación, editar el código, enviar!!!
                            </p>
                            <p>
                                Recomendación: Pegar el código desde algún Sandbox en linea del lenguaje seleccionado.
                            </p>
                            
                            <form className="service-form" action="/submit-service" method="POST">
                                <textarea className="code-input" name="code" placeholder="Placeholder tiene que ser el boilerplate de cada lenguage"/>
                                <div>
                                    <input className="monitor-button" type="submit" value="MANDAR"/>
                                    <select className="drop-down-menu" name="lang-select-menu">
                                        <option value="rust">Rust</option>
                                        <option value="python">Python</option>
                                        <option value="go">Go</option>
                                    </select>
                                </div>
                            </form>

                        </div>
                    </div>
                </div>
            </div>
        </div>
        </div>
    )
}