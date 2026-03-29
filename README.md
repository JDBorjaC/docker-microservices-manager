**Docker Microservices Manager**
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
+ **Gestión de estado real:** Consume los flujos de Server-Sent Events (SSE) del orquestador para actualizar instantáneamente la UI cuando un contenedor arranca, se detiene o falla.
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
+ **Inyección de red:** Es el orquestador quien le dice a Docker que agregue las etiquetas mágicas de Traefik (traefik.http.routers...) a los contenedores que crea.
+ **Bucle de Reconciliación (Reconciliation Loop):** Mantiene una goroutine abierta escuchando los eventos nativos del host de Docker en tiempo real. Esto permite sincronizar siempre el estado del servicio con la base de datos interna.
+ **Persistencia SQLite (data.db):** Almacena y sincroniza el estado internal de cada servicio (ID, contenedor, código y estatus actual).

**Endpoints Principales del Orchestrator:**
+ **POST /microservices** - Empaqueta el código, construye la imagen y levanta el contenedor.
+ **GET /microservices/status/events** - Abre un canal SSE para transmitir eventos en vivo.
+ **PATCH /microservices/start/:id | PATCH /microservices/stop/:id** - Gestión de estado del contenedor.
+ **DELETE /microservices/:id** - Elimina el contenedor, destruye la imagen y borra los registros.

**Contenedores Dinámicos**
 -
 Este no es un módulo estático, sino el resultado de la acción del proyecto. Son los contenedores aislados que cobran vida bajo demanda según el lenguaje elegido por el usuario (Express, Flask, Gin, Cargo).
 **Características:**
 + Totalmente aislados de la infraestructura crítica, pero conectados a la red msm-network para ser visibles por Traefik.
 + Ciclo de vida efímero administrado totalmente por el Orchestrator. Se reconstruyen completamente si el usuario edita el código fuente.
   
**Diagrama de Arquitectura**
![Diagram de arquitectura](https://github.com/user-attachments/assets/0836ee72-a853-4c03-91fb-c9219267a419)
