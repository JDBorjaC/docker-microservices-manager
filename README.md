**Diagrama de Arquitectura**
-
![Diagram de arquitectura](https://github.com/user-attachments/assets/0836ee72-a853-4c03-91fb-c9219267a419)

**Ejemplos Funcionales**
-
A continuación se muestran ejemplos funcionales básicos de código para los microservicios, listos para desplegarse:

### Flask (Python) - Hola Mundo
```python
from flask import Flask, request
app = Flask(__name__)

@app.route('/', methods=['GET'])
def hola():
    return "Hola Mundo"

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000)
```

### Flask (Python) - Sumar
```python
from flask import Flask, request
app = Flask(__name__)

@app.route('/sumar', methods=['GET'])
def sumar():
    a = request.args.get('a', default=0, type=int)
    b = request.args.get('b', default=0, type=int)
    
    resultado = a + b
    return f"La suma de {a} y {b} es: {resultado}"

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000)
```

### Express (Node.js) - Hola Mundo
```javascript
const express = require('express');
const app = express();
const PORT = 3000;

app.get('/', (req, res) => {
    res.send("Hola Mundo");
});

app.listen(PORT, '0.0.0.0', () => {
    console.log(`Microservicio escuchando correctamente en el puerto ${PORT}`);
});
```

### Express (Node.js) - Sumar
```javascript
const express = require('express');
const app = express();
const PORT = 3000;

app.get('/sumar', (req, res) => {
    // Obtener parámetros desde la URL: /sumar?a=10&b=20
    const a = parseInt(req.query.a) || 0;
    const b = parseInt(req.query.b) || 0;
    
    const resultado = a + b;
    res.send(`La suma de ${a} y ${b} es: ${resultado}`);
});

app.listen(PORT, '0.0.0.0', () => {
    console.log(`Microservicio escuchando correctamente en el puerto ${PORT}`);
});
```

**Docker Microservices Manager**
-
Docker Microservices Manager es una plataforma de orquestación ligera tipo PaaS (Platform as a Service) diseñada para simplificar la creación, el despliegue y la gestión del ciclo de vida de microservicios. A través de una interfaz de usuario, permite a los desarrolladores escribir lógica en múltiples lenguajes de programación y generar contenedores Docker auto-enrutados en tiempo real de manera aislada, sin tener que configurar infraestructuras complejas.

**Arquitectura Técnica**
-
 Aquí encontrará información sobre los detalles técnicos de todos los módulos que conforman la arquitectura de este proyecto.
 
 **Módulos del Sistema**
 -
La arquitectura del proyecto está dividida en tres módulos principales y una capa dinámica:
 + Frontend (React / Vite)
 + Reverse Proxy (Traefik)
 + Orchestrator API (Go / Gin)
 + Contenedores Dinámicos

**Frontend (React / Vite)**
-
Este módulo representa la interfaz gráfica de usuario. Construida con React y empaquetada con Vite, proporciona el panel de control (Dashboard) donde los usuarios interactúan con el sistema.
**Responsabilidades:**
+ **Creación visual:** Provee un editor de código en el navegador para que el usuario escriba la lógica de su nuevo microservicio.
+ **Gestión de estado real:** Consume los flujos de Server-Sent Events (SSE) del orquestador para actualizar el estado de un contenedor (ej. de Created a Running o Stopped).
+ **Enrutamiento HTTP:** Realiza peticiones asíncronas (fetch) al backend para gestionar la vida de los servicios (Start, Stop, Delete).

**Reverse Proxy (Traefik)**
-
El proxy inverso es el controlador de tráfico del sistema, escuchando en los puertos :80 (HTTP web) y :8090 (Dashboard). Traefik actúa como la capa de red pública.
**Responsabilidades:**
+ **Auto-descubrimiento:** Se conecta al daemon de Docker (/var/run/docker.sock) para "escuchar" la aparición de nuevos contenedores en la red msm-network.
+ **Enrutamiento dinámico:** Cuando un nuevo microservicio nace, Traefik lee sus etiquetas inyectadas y enruta el tráfico que va al prefijo /services/<nombre-del-servicio> directamente hacia ese contenedor, eliminando el prefijo antes de que el tráfico entre al contenedor.
+ **Redirección del Frontend:** Envía las peticiones a la ruta raíz / directamente al servidor de vite en el puerto 5173.
  
**Orchestrator API (Go / Gin)**
-
Este es el servicio más crítico y el "cerebro" central de la operación. Construido en Go usando el framework Gin, expone la API REST que el Frontend consume.
**Responsabilidades y Flujo:**
+ **Manipulación de Docker:** Utiliza el API nativo de Golang para Docker (github.com/docker/docker/client). Cuando recibe código del frontend, genera un Dockerfile en tiempo de ejecución, empaqueta el directorio, y le ordena al Docker Daemon compilar la imagen (BuildImage
) para luego crear un contenedor (ContainerCreate).
  
  **Nota arquitectónica clave:** Crear un contenedor y arrancar un contenedor son **dos operaciones totalmente independientes**. La creación (Create) solo aprovisiona el contenedor en el servidor y lo deja en estado inactivo. Iniciarlo (Start) es lo que pone a correr el código interno.
+ **Inyección de red:** Es el orquestador quien le dice a Docker que agregue las etiquetas mágicas de Traefik (traefik.http.routers...) a los contenedores que crea.
+ **Bucle de Reconciliación (Reconciliation Loop):** Mantiene una goroutine abierta escuchando los eventos nativos del host de Docker en tiempo real. Esto permite sincronizar siempre el estado del servicio con la base de datos interna.
+ **Persistencia SQLite (data.db):** Almacena y sincroniza el estado internal de cada servicio (ID, contenedor, código y estatus actual).

**Endpoints Principales del Orchestrator:**
+ **POST /microservices** - Empaqueta el código, construye la imagen y crea el contenedor de Docker preparado para ejecutarse (estado inactivo).
+ **GET /microservices/status/events** - Abre un canal SSE para transmitir eventos en vivo.
+ **PATCH /microservices/start/:id** - Operación explícita para iniciar la ejecución de un contenedor que fue creado previamente.
+ **PATCH /microservices/stop/:id** - Operación para detener la ejecución de un contenedor sin destruirlo (lo devuelve a estado inactivo).
+ **DELETE /microservices/:id** - Elimina el contenedor, destruye la imagen y borra los registros.

**Contenedores Dinámicos**
 -
 Este no es un módulo estático, sino el resultado de la acción del proyecto. Son los contenedores aislados que cobran vida bajo demanda según el lenguaje elegido por el usuario (Express, Flask).
 
 **Características:**
 + Totalmente aislados de la infraestructura crítica, pero conectados a la red msm-network para ser visibles por Traefik.
 + Ciclo de vida efímero administrado totalmente por el Orchestrator. Se reconstruyen completamente si el usuario edita el código fuente.
   
