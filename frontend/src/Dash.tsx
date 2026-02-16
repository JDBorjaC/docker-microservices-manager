import { useState } from 'react'
import './styles/monitor.css'

export default function Dash() {
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
        // Logic to edit the service
    };

    const deleteService = (id:number) => {
        setServices(services.filter(service => service.id !== id));
    };

    return (
        <div className="dash-div">

        <div className="monitor-frame">
            <div className="monitor-bezel">
                <div className="monitor-screen">
                    <div className='monitor-scanlines'>
                    <div className="monitor-content">
                        
                        {/* TERMINAL DE MICROSERVICIOS */}
                        {/* Cada microservicio muestra NOMBRE, BTN_ON, BTN_EDIT, BTN_DEL */}
                        <div className='services-list'>
                            {services.map(service => (
                                <div key={service.id} className='service-item'>
                                    <span>{service.name}</span>
                                    <button className='monitor-button' onClick={() => toggleService(service.id)}>{service.isActive ? 'Turn Off' : 'Turn On'}</button>
                                    <button className='monitor-button' onClick={() => editService(service.id)}>Edit</button>
                                    <button className='monitor-button' onClick={() => deleteService(service.id)}>Delete</button>
                                </div>
                            ))}
                        </div>
                        
                        {/* BTN_CREATE MICROSERVICIO */}
                        <div className='create-service-div'>
                            <button className='monitor-button'>
                                CREAR MICROSERVICIO
                            </button>
                        </div>

                    </div>
                    </div>
                </div>
            </div>
        </div>
        </div>
    )
}
