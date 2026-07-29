### Prompt 1
Estoy desarrollando una API REST para una academia ficticia, que gestione las nuevas solicitudes operativas, usando Go, fiber y PostgreSQL en supabase
Necesito de tu ayuda para que me guies en un paso a paso de como diseñar una base de datos relacional con persistencia y poder generar un archivo init.sql
La base de datos se va a encargar de gestionar las nuevas solicitudes operativas que tiene los siguientes campos:

id (UUID)
titulo (string)
area_solicitante (string)
prioridad (un numero entero entre 1 y 5 )
costo_estimado (un numero decimal)
estado (registrada, en_proceso y finalizada)

Necesito que agregues las restricciones necesarias como Primary key, not null, etc. Siguiendo buenas practicas de PostgreSQL.
Necesito que expliques la eleccion de cada tipo de datos y necesito que el codigo SQL sea seguro, limpio y facil de mantener preparado para una API REST 
donde se aplicaran principios SOLID
Necesito que me des cada detalle del paso a paso.

### Respuesta
![alt text](./Imagenes/Respuesta1_1.jpeg)
![alt text](./Imagenes/Respuesta1_2.png)
![alt text](./Imagenes/Respuesta1_3.png)
![alt text](./Imagenes/Respuesta1_4.png)
![alt text](./Imagenes/Respuesta1_5.png)
![alt text](./Imagenes/Respuesta1_6.png)
![alt text](./Imagenes/Respuesta1_7.png)
![alt text](./Imagenes/Respuesta1_8.png)
![alt text](./Imagenes/Respuesta1_9.png)
![alt text](./Imagenes/Respuesta1_10.png)
![alt text](./Imagenes/Respuesta1_11.png)
![alt text](./Imagenes/Respuesta1_12.png)
![alt text](./Imagenes/Respuesta1_13.png)
![alt text](./Imagenes/Respuesta1_14.png)
![alt text](./Imagenes/Respuesta1_15.png)
![alt text](./Imagenes/Respuesta1_16.png)
![alt text](./Imagenes/Respuesta1_17.png)
![alt text](./Imagenes/Respuesta1_18.png)
![alt text](./Imagenes/Respuesta1_19.png)

### Prompt 2
Estoy desarrollando una API REST para una academia ficticia utilizando Go, Fiber y PostgreSQL en Supabase.
Necesito que me guíes paso a paso para diseñar una arquitectura de monolito modular por capas, utilizando MVC + Service + Repository, siguiendo las buenas prácticas de Clean Architecture y aplicando los principios SOLID.

Antes de generar código, explica el propósito y la responsabilidad de cada carpeta y cómo se relacionan entre sí. La arquitectura debe evitar la duplicación de código, utilizar nombres descriptivos, funciones con responsabilidad única y estar preparada para crecer y mantenerse fácilmente.
Además, configura el proyecto para utilizar variables de entorno mediante un archivo .env para las credenciales de Supabase y explica la función de cada variable.
No avances a la implementación de los endpoints hasta que la arquitectura esté correctamente diseñada y explicada.

### Respuesta
![alt text](./Imagenes/Respuesta2_0.png)
![alt text](./Imagenes/Respuesta2_1.png)
![alt text](./Imagenes/Respuesta2_2.png)
![alt text](./Imagenes/Respuesta2_3.png)
![alt text](./Imagenes/Respuesta2_4.png)
![alt text](./Imagenes/Respuesta2_5.png)
![alt text](./Imagenes/Respuesta2_6.png)
![alt text](./Imagenes/Respuesta2_7.png)
![alt text](./Imagenes/Respuesta2_8.png)
![alt text](./Imagenes/Respuesta2_9.png)
![alt text](./Imagenes/Respuesta2_10.png)
![alt text](./Imagenes/Respuesta2_11.png)
![alt text](./Imagenes/Respuesta2_12.png)
![alt text](./Imagenes/Respuesta2_14.png)
![alt text](./Imagenes/Respuesta2_15.png)
![alt text](./Imagenes/Respuesta2_16.png)
![alt text](./Imagenes/Respuesta2_17.png)
![alt text](./Imagenes/Respuesta2_18.png)
![alt text](./Imagenes/Respuesta2_19.png)
![alt text](./Imagenes/Respuesta2_21.png)

### Prompt 3
Estoy desarrollando una API REST para una academia ficticia utilizando Go, Fiber y PostgreSQL en Supabase.
La arquitectura del proyecto ya está definida de la siguiente manera:
academia-api/
├── cmd/api/
├── internal/
│   ├── config/
│   ├── infrastructure/database/
│   ├── modules/
│   │   └── solicitudes/
│   │       ├── model/
│   │       ├── dto/
│   │       ├── repository/
│   │       ├── service/
│   │       ├── controller/
│   │       └── routes/
│   └── shared/
│       ├── errors/
│       ├── response/
│       └── validator/
├── migrations/
├── .env
├── .env.example
├── .gitignore
├── go.mod
└── README.md

Necesito crear un CRUD completo para la entidad "solicitud_operativa", respetando esta arquitectura y sin crear carpetas adicionales innecesarias.
La entidad contiene:

id (UUID)
titulo (string)
area_solicitante (string)
prioridad (un numero entero entre 1 y 5 )
costo_estimado (un numero decimal)
estado (registrada, en_proceso y finalizada)

Se debe implementar los siguientes 5 endpoints:
GET /solicitudes: obtener todas las solicitudes.
POST /solicitudes: registrar una solicitud.
PUT /solicitudes/:id : actualizar completamente una solicitud.
DELETE /solicitudes/:id : eliminar una solicitud.
PATCH /solicitudes/:id/estado: actualizar únicamente el estado.

El código debe aplicar principios SOLID, respetar la separación entre controller, service y repository, validar las entradas, usar consultas SQL parametrizadas, prevenir inyección SQL, manejar errores correctamente y devolver respuestas JSON con códigos HTTP apropiados.

Utiliza DTO separados para crear, actualizar completamente y actualizar únicamente el estado. Evita duplicación de código, usa nombres descriptivos y funciones con responsabilidad única.

Genera todos los archivos necesarios para que los 5 endpoints funcionen, indicando la ruta de cada archivo antes de mostrar su contenido. Al finalizar, muestra ejemplos JSON para probar cada endpoint en Postman.


### Respuesta
![alt text](./Imagenes/Respuesta3_0
.png)
![alt text](./Imagenes/Respuesta3_1.png)
![alt text](./Imagenes/Respuesta3_2.png)
![alt text](./Imagenes/Respuesta3_3.png)
![alt text](./Imagenes/Respuesta3_4.png)
![alt text](./Imagenes/Respuesta3_5.png)
![alt text](./Imagenes/Respuesta3_6.png)
![alt text](./Imagenes/Respuesta3_7.png)
![alt text](./Imagenes/Respuesta3_8.png)
![alt text](./Imagenes/Respuesta3_9.png)
![alt text](./Imagenes/Respuesta3_10.png)
![alt text](./Imagenes/Respuesta3_11.png)
![alt text](./Imagenes/Respuesta3_12.png)
![alt text](./Imagenes/Respuesta3_13.png)
![alt text](./Imagenes/Respuesta3_14.png)
![alt text](./Imagenes/Respuesta3_15.png)
![alt text](./Imagenes/Respuesta3_16.png)
![alt text](./Imagenes/Respuesta3_17.png)
![alt text](./Imagenes/Respuesta3_18.png)
![alt text](./Imagenes/Respuesta3_19.png)
![alt text](./Imagenes/Respuesta3_20.png)
![alt text](./Imagenes/Respuesta3_21.png)
![alt text](./Imagenes/Respuesta3_22.png)
![alt text](./Imagenes/Respuesta3_23.png)
![alt text](./Imagenes/Respuesta3_24.png)
![alt text](./Imagenes/Respuesta3_25.png)
![alt text](./Imagenes/Respuesta3_26.png)
![alt text](./Imagenes/Respuesta3_27.png)
![alt text](./Imagenes/Respuesta3_28.png)
![alt text](./Imagenes/Respuesta3_29.png)
![alt text](./Imagenes/Respuesta3_30.png)
![alt text](./Imagenes/Respuesta3_31.png)
![alt text](./Imagenes/Respuesta3_32.png)
![alt text](./Imagenes/Respuesta3_33.png)
![alt text](./Imagenes/Respuesta3_34.png)
![alt text](./Imagenes/Respuesta3_35.png)
![alt text](./Imagenes/Respuesta3_36.png)
![alt text](./Imagenes/Respuesta3_37.png)
![alt text](./Imagenes/Respuesta3_38.png)