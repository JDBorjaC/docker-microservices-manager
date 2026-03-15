import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import type { Service } from '../models/msm_models';
import '../styles/monitor.css';

export default function Dash() {

    const navi = useNavigate();

    const [services, setServices] = useState([
        { id: "1", name: 'Service 1', status: false },
        { id: "2", name: 'Service 2', status: false },
        { id: "3", name: 'Service 3', status: false }
    ]);

    const toggleService = (id:string):void => {
        setServices(services.map(service => 
            service.id === id ? { ...service, status: !service.status } : service
        ));
    };

    const editService = async (id:string):Promise<void> => {
        //acceder al archivo del microservicio() dentro del contenedor, para cargarlo en la pag siguiente
        //espero que el backend me mande un Service con todos los campos definidos

        //TODO: descomentar cuando el backend exista
        /*const backResponse = await fetch("http://backend:8000/get-by-id/"+id);
        const userCode:Service = await backResponse.json();*/
        
        const userCode:Service = {id:id, name:"name", status:true, code:"blahblah blah blah"}
        
        try {
            localStorage.setItem('editService', JSON.stringify(userCode));
            navi("/edit")
        } catch (err) {
            console.error("Algo salió mal tratando de guardar el servicio a editar:\n", err);
        }
    };

    const deleteService = async (id:string):Promise<void> => {
        //mandar una petición de Borrar al backend, esperar a que me devuelva la lista de servicios restantes
        const backResponse = await fetch("http://backend:8000/delete-service/"+id, { method:'DELETE' });
        const remainingServices:Service[] = await backResponse.json();

        setServices(remainingServices);

    };

    //Hacer un fetch al backend y parsear los contenedores que esten listeados
    const fetchServices = async ():Promise<void> => {
        /* TODO: hacer url una variable de entorno maybe */
        //espero que el backend me mande un arreglo de Service, pero sin el campo de 'code' definido
        const backResponse =  await fetch("http://backend:8000/get-services");
        const services:Service[] = await backResponse.json();  // toca estar atento a cambios del modelo
        setServices(services); //parsear bien el status!!!
    }

    //ON PAGELOAD: fetch services from backend
    useEffect(():void => {
        fetchServices()
    }, [])


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

                                {/* Se itera sobre cada servicio recibido del backend para añadirlo a una lista */}
                                {services.map(service => (
                                    <div key={service.id} className='monitor-item'>
                                        <div className='service-label'>
                                            <h3>{service.name}</h3>
                                            <span>https://www.test.com.co/sumar-microservicio/pyth</span>
                                        </div>
                                        <div className='monitor-item-buttons'>
                                            <button className='monitor-button' onClick={() => toggleService(service.id)}>{service.status ? 'Turn Off' : 'Turn On'}</button>
                                            <button className='monitor-button' onClick={() => editService(service.id)}>Edit</button>
                                            <button className='monitor-button' onClick={() => deleteService(service.id)}>Delete</button>
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
        </div>
        {/*--- fin de terminal de microservicios ---*/}

        </div>
    )
}
