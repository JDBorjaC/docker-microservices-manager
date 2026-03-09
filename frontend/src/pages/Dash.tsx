import { useState } from 'react'
import { useNavigate } from 'react-router-dom';
import '../styles/monitor.css'

export default function Dash() {

    const navi = useNavigate();

    const [services, setServices] = useState([
        { id: 1, name: 'Service 1', isActive: false },
        { id: 2, name: 'Service 2', isActive: false },
        { id: 3, name: 'Service 3', isActive: false }
    ]);

    const toggleService = (id:number) => {
        setServices(services.map(service => 
            service.id === id ? { ...service, isActive: !service.isActive } : service
        ));
    };

    const editService = (id:number) => {
        console.log('Editing ', id)
        navi('/edit')
    };

    const deleteService = (id:number) => {
        setServices(services.filter(service => service.id !== id));
    };


    return (
        <div className="dash-div">
        
        {/* TERMINAL DE MICROSERVICIOS */}
        <div className="monitor-frame">
            <div className="monitor-bezel">
                <div className="monitor-screen">
                    <div className='monitor-scanlines'>
                        <div className="monitor-content">

                            <h1>
                                \\ LISTADO DE MICROSERVICIOS
                            </h1>

                            {/* Cada microservicio muestra NOMBRE, BTN_ON, BTN_EDIT, BTN_DEL */}
                            <div className='monitor-list'>

                                {services.map(service => (
                                    <div key={service.id} className='monitor-item'>
                                        <div className='service-label'>
                                            <h3>{service.name}</h3>
                                            <span>https://www.test.com.co/sumar-microservicio/pyth</span>
                                        </div>
                                        <div className='monitor-item-buttons'>
                                            <button className='monitor-button' onClick={() => toggleService(service.id)}>{service.isActive ? 'Turn Off' : 'Turn On'}</button>
                                            <button className='monitor-button' onClick={() => editService(service.id)}>Edit</button>
                                            <button className='monitor-button' onClick={() => deleteService(service.id)}>Delete</button>
                                        </div>
                                    </div>
                                ))}
                                {/* BTN_CREATE MICROSERVICIO */}
                                <div className='create-service-div'>
                                    <button className='monitor-button' onClick={() => editService(-1337)}>
                                        CREAR MICROSERVICIO
                                    </button>
                                </div>
                                
                            </div>

                        </div>
                    </div>
                </div>
            </div>
        </div>

        {/*--- COMPONENT END ---*/}
        </div>
    )
}
