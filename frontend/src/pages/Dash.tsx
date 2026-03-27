import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import type { Service } from '../models/msm_models';
import '../styles/monitor.css';
import MonitorBackdrop from '../components/monitor';

const backend_service:string = "http://localhost:8080/"

export function Dash() {

    const navi = useNavigate();

    const [services, setServices] = useState<Service[]>([]);

    const editService = async (id:string):Promise<void> => {
        //la info del service a editar se carga on pageload
        navi("/deets/"+id);
    };

    const deleteService = async (id:string):Promise<void> => {
        //mandar una petición de Borrar al backend, esperar a que me devuelva la lista de servicios restantes
        const backResponse = await fetch(backend_service+"microservices/"+id, { method:'DELETE' });

        if(backResponse.status === 200) {
            const remainingServices:Service[] = services.filter((s) => s.id !== id);
            setServices(remainingServices);
        } else {
            console.error("No se pudo borrar el microservicio:\n", backResponse.body);
        }

    };

    //Hacer un fetch al backend y parsear los contenedores que esten listeados
    const fetchServices = async ():Promise<void> => {
        /* TODO: hacer url una variable de entorno maybe */
        //espero que el backend me mande un arreglo de Service, pero sin el campo de 'code' definido
        const backResponse =  await fetch(backend_service+"microservices");
        const services:Service[] = await backResponse.json();  // toca estar atento a cambios del modelo
        setServices(services); 
    }

    //ON PAGELOAD: fetch services from backend
    useEffect(():void => {
        fetchServices()
    }, [])


    return (
        <div className="dash-div">
        
        {/* TERMINAL DE MICROSERVICIOS */}
        <MonitorBackdrop>
            <div className="monitor-bezel">
                <div className="monitor-screen">
                    <div className='monitor-scanlines'>
                        <div className="monitor-content">

                            <h1>
                                \\ LISTADO DE MICROSERVICIOS
                            </h1>

                            {/* Cada microservicio muestra NOMBRE, BTN_ON, BTN_EDIT, BTN_DEL */}
                            <div className='monitor-list'>

                                {/* Se itera sobre cada servicio recibido del backend para añadirlo a una lista */}
                                {services.map(service => (
                                    <div key={service.id} className='monitor-item'>
                                        <div className='service-label'>
                                            <h3>{service.name}</h3>
                                            <span>{"http://localhost/services/"+service.name}</span>
                                        </div>
                                        <div className='monitor-item-buttons'>
                                            <button className='monitor-button'>{service.status}</button>
                                            <button className='monitor-button' onClick={() => editService(service.id)}>opciones</button>
                                            <button className='monitor-button' onClick={() => deleteService(service.id)}>eliminar</button>
                                        </div>
                                    </div>
                                ))}

                                {/* BTN_CREATE MICROSERVICIO */}
                                <div className='create-service-div'>
                                    <button className='monitor-button' onClick={() => navi("/edit")}>
                                        CREAR MICROSERVICIO
                                    </button>
                                </div>
                                

                            </div>

                        </div>
                    </div>
                </div>
            </div>
        </MonitorBackdrop>
        {/*--- fin de terminal de microservicios ---*/}

        </div>
    )
}

