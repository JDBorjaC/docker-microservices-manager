import { useEffect, useState } from 'react';
import type { Service } from '../models/msm_models'



export default function ServiceEditor(){

    const [editable, setEditable] = useState(():Service => {
        //cargar editService si es que está
        const e = localStorage.getItem("editService");
        console.log(e)
        return e ? JSON.parse(e):{
            id:'',
            name:'',
            status:false,
            code:''
        }
    });

    //para mandar MS: nombre, desc?, lenguage y codigo

    // Si se carga un servicio, colocar su código en el textArea

    
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
                            <p>
                                Advertencia: Para que su microservicio funcione, tiene que definir una función 'microservice()', que actuará como la función principal que será ejecutada. 
                            </p>
                            
                            <form className="service-form" action="/submit-service" method="POST">
                                {/* El textarea carga de local-storage cada page-load. es decir, refrescar la página supone perder cambios. Para el scope del proyecto no importa mucho.*/}
                                <textarea className="code-input" name="code" placeholder="¡Copiar y pegar código aquí!">{editable.code}</textarea>
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