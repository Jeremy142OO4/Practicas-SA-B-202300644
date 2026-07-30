# UNIVERSIDAD DE SAN CARLOS DE GUATEMALA  
## FACULTAD DE INGENIERÍA  
### SOFTWARE AVANZADO

**CATEDRÁTICO:** Ing. Marco Tulio Aldana Prillwitz
**TUTOR ACADÉMICO:** Kevin Pozuelos Estrada

---

**Jeremy Estuado Orellana Aldana**  
**CARNÉ:** 202300644
**SECCIÓN:** B

**Guatemala, 30 de Julio del 2026**
# Academia API

## Índice

1. [Descripción](#descripcion)
2. [Tecnologías utilizadas](#tecnologías-utilizadas)
3. [Funcionalidades](#funcionalidades)
4. [Estructura del proyecto](#estructura-del-proyecto)
5. [Responsabilidad de las carpetas](#responsabilidad-de-las-carpetas)
6. [Instalación](#instalación)
7. [Configuración de variables de entorno](#configuración-de-variables-de-entorno)
8. [Configuración de la base de datos](#configuración-de-la-base-de-datos)
9. [Ejecución](#ejecución)
10. [Endpoints](#endpoints)
11. [Formato de respuestas](#formato-de-respuestas)
12. [Validaciones](#validaciones)
13. [Seguridad](#seguridad)
14. [Principios SOLID](#principios-solid)
15. [Código limpio](#código-limpio)
16. [Pruebas con Postman](#pruebas-con-postman)
17. [Verificación del proyecto](#verificación-del-proyecto)
18. [Uso de inteligencia artificial](#uso-de-inteligencia-artificial)


### Descripcion
API REST desarrollada en Go para administrar solicitudes operativas de una academia ficticia. Permite registrar, consultar, actualizar, cambiar el estado y eliminar solicitudes. Los datos se almacenan en PostgreSQL mediante Supabase.

### Tecnologías utilizadas
Incluye:
- Go
- Fiber
- PostgreSQL
- Supabase
- pgx
- UUID
- Postman

### Funcionalidades

La API permite realizar las siguientes operaciones:

- Registrar una solicitud operativa.
- Obtener todas las solicitudes registradas.
- Actualizar completamente una solicitud.
- Actualizar únicamente el estado de una solicitud.
- Eliminar una solicitud.

Las solicitudes tienen el siguiente formato:

| Campo            | Tipo    | Descripción                              |
|------------------|---------|------------------------------------------|
| id             | UUID    | Identificador único                      |
| titulo         | String  | Nombre descriptivo de la solicitud       |
| area_solicitante | String | Área que realiza la solicitud           |
| prioridad      | Entero  | Nivel de prioridad entre 1 y 5           |
| costo_estimado | Decimal | Costo aproximado                         |
| estado         | String  | registrada, en_proceso o finalizada |
| created_at     | Fecha   | Fecha de creación                        |
| updated_at     | Fecha   | Fecha de última actualización            |

### Estructura del proyecto


```text
academia-api/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── infrastructure/
│   │   └── database/
│   │       └── postgres.go
│   ├── modules/
│   │   └── solicitudes/
│   │       ├── controller/
│   │       │   └── controller.go
│   │       ├── dto/
│   │       │   ├── create_request.go
│   │       │   ├── response.go
│   │       │   ├── update_estado_request.go
│   │       │   └── update_request.go
│   │       ├── model/
│   │       │   └── solicitud.go
│   │       ├── repository/
│   │       │   ├── repository.go
│   │       │   └── postgres_repository.go
│   │       ├── routes/
│   │       │   └── routes.go
│   │       └── service/
│   │           └── service.go
│   └── shared/
│       ├── errors/
│       │   └── app_error.go
│       ├── response/
│       │   └── response.go
│       └── validator/
│           └── validator.go
├── migrations/
│   └── 001_create_solicitudes_operativas.sql
├── postman/
│   └── academia-api.postman_collection.json
├── .env.example
├── .gitignore
├── go.mod
├── go.sum
├── PROMPTS.md
└── README.md
```

### Responsabilidad de las carpetas

- cmd/api: contiene el punto de entrada de la aplicación.
- config: carga y valida las variables de entorno.
- infrastructure/database: crea y configura la conexión con PostgreSQL.
- model: representa las entidades principales del sistema.
- dto: define los datos aceptados y devueltos por la API.
- repository: contiene las operaciones de persistencia.
- service: implementa las reglas de negocio y validaciones.
- controller: recibe las peticiones HTTP y devuelve las respuestas.
- routes: registra las rutas de Fiber.
- shared: contiene errores, respuestas y validaciones reutilizables.
- migrations: contiene la estructura de la base de datos.
- postman: contiene la colección utilizada para probar la API.

### Instalación

**Clonar el repositorio:**
``` bash
git clone https://github.com/Jeremy142OO4/Practicas-SA-B-202300644.git
```
**Entrar al directorio de la aplicación:**
``` bash
cd Practicas-SA-B-202300644/P1/academia-api
```
**Descargar las dependencias:**
```bash
go mod download
```
### Configuración de variables de entorno

Crear una copia de .env.example con el nombre .env.

**En Git Bash:**
```
cp .env.example .env
```
**En PowerShell:**
```
Copy-Item .env.example .env
```
Configurar las siguientes variables:
```
APP_NAME=academia-api
APP_ENV=development
APP_HOST=0.0.0.0
APP_PORT=8080

DATABASE_URL=postgresql://USER:PASSWORD@HOST:PORT/postgres?sslmode=require
DB_HOST=HOST
DB_PORT=5432
DB_DATABASE=postgres
DB_USER=USER
DB_PASSWORD=PASSWORD

DATABASE_MAX_CONNECTIONS=10
DATABASE_MIN_CONNECTIONS=2
DATABASE_CONNECT_TIMEOUT_SECONDS=10
```
DATABASE_URL debe contener la cadena de conexión proporcionada por Supabase.

**El archivo .env contiene información sensible y está incluido en .gitignore, por lo que no debe subirse al repositorio.**

### Configuración de la base de datos

El archivo SQL necesario para crear la tabla se encuentra en:
```
migrations/001_create_solicitudes_operativas.sql
```
**Para preparar la base de datos:**

1. Abrir el proyecto en Supabase.
2. Entrar al SQL Editor.
3. Copiar el contenido de la migración.
4. Ejecutar el script.

**La migración crea:**

- La tabla solicitudes_operativas.
- Un identificador UUID automático.
- Restricciones para campos obligatorios.
- Prioridad entre 1 y 5.
- Costo estimado no negativo.
- Estados permitidos.
- Índices para optimizar consultas.
- Un trigger para actualizar updated_at.

### Ejecución

Iniciar la API con:
``` bash
go run ./cmd/api
``` 
La aplicación estará disponible en:
``` 
http://localhost:8080
``` 
La URL base de los endpoints es:
```
http://localhost:8080/api/v1
```

### Endpoints
| Método | Endpoint | Descripción |
|--------|----------|-------------|
| GET | /api/v1/solicitudes | Obtener todas las solicitudes |
| POST | /api/v1/solicitudes | Registrar una solicitud |
| PUT | /api/v1/solicitudes/:id | Actualizar completamente una solicitud |
| PATCH | /api/v1/solicitudes/:id/estado | Actualizar únicamente el estado |
| DELETE | /api/v1/solicitudes/:id | Eliminar una solicitud |

**Registrar una solicitud**
```
POST /api/v1/solicitudes/
```
Cuerpo de la petición:
```
{
  "titulo": "Adquisición de nuevo servidor",
  "area_solicitante": "Infraestructura TI",
  "prioridad": 3,
  "costo_estimado": 2500.00
}
```

Al registrar una solicitud, el estado inicial se establece automáticamente como registrada.

Respuesta esperada:
```
{
  "success": true,
  "data": {
    "id": "UUID_GENERADO",
    "titulo": "Adquisición de nuevo servidor",
    "area_solicitante": "Infraestructura TI",
    "prioridad": 3,
    "costo_estimado": 2500,
    "estado": "registrada",
    "created_at": "2026-07-29T00:00:00Z",
    "updated_at": "2026-07-29T00:00:00Z"
  }
}
```
Código HTTP:
```
201 Created
```

**Obtener todas las solicitudes**
```
GET /api/v1/solicitudes/
```
Respuesta esperada:
```
{
  "success": true,
  "data": [
    {
      "id": "UUID_DE_LA_SOLICITUD",
      "titulo": "Adquisición de nuevo servidor",
      "area_solicitante": "Infraestructura TI",
      "prioridad": 3,
      "costo_estimado": 2500,
      "estado": "registrada",
      "created_at": "2026-07-29T00:00:00Z",
      "updated_at": "2026-07-29T00:00:00Z"
    }
  ]
}
```
Código HTTP:
```
200 OK
```
**Actualizar completamente una solicitud**
```
PUT /api/v1/solicitudes/:id
```
Cuerpo de la petición:
```
{
  "titulo": "Adquisición de servidor actualizada",
  "area_solicitante": "Tecnología",
  "prioridad": 4,
  "costo_estimado": 3000.50,
  "estado": "en_proceso"
}
```
Código HTTP esperado:
```
200 OK
```
Si el identificador no existe, la API devuelve:
```
404 Not Found
```
**Actualizar únicamente el estado**
```
PATCH /api/v1/solicitudes/:id/estado
```
Cuerpo de la petición:
```
{
  "estado": "finalizada"
}
```
Los estados aceptados son:

- registrada
- en_proceso
- finalizada

Código HTTP esperado:
```
200 OK
```
**Eliminar una solicitud**
```
DELETE /api/v1/solicitudes/:id
```
Código HTTP esperado:
```
204 No Content
```
Si la solicitud no existe, la API devuelve:
```
404 Not Found
```

### Formato de respuestas

Las respuestas exitosas utilizan la siguiente estructura:
```
{
  "success": true,
  "data": {}
}
```
Los errores utilizan la siguiente estructura:
```
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "descripción del error"
  }
}
```

### Validaciones

La aplicación comprueba que:

- titulo sea obligatorio.
- titulo no supere 150 caracteres.
- area_solicitante sea obligatoria.
- area_solicitante no supere 100 caracteres.
- prioridad esté entre 1 y 5.
- costo_estimado sea igual o mayor que cero.
- El identificador proporcionado sea un UUID válido.
- El estado sea registrada, en_proceso o finalizada.
- El cuerpo enviado tenga un formato JSON válido.

**Los espacios innecesarios al inicio y al final del título y del área solicitante son eliminados antes de guardar los datos.**

### Seguridad

**Consultas SQL parametrizadas**

El repositorio utiliza parámetros como \$1, \$2, \$3, \$4 y \$5 en lugar de concatenar los datos directamente en las consultas.

**Ejemplo real de postgres_repository.go:**
```
const query = `
    INSERT INTO public.solicitudes_operativas (
        titulo,
        area_solicitante,
        prioridad,
        costo_estimado,
        estado
    )
    VALUES ($1, $2, $3, $4, $5)
`

_, err := r.db.Exec(
    ctx,
    query,
    solicitud.Titulo,
    solicitud.AreaSolicitante,
    solicitud.Prioridad,
    solicitud.CostoEstimado,
    solicitud.Estado,
)
```
**Esto permite que PostgreSQL trate los valores como datos y ayuda a prevenir ataques de inyección SQL.**

**Protección de credenciales**

La cadena de conexión se obtiene desde DATABASE_URL. Las credenciales no se escriben directamente en el código fuente.

El archivo .env está excluido mediante .gitignore.

**Manejo de errores**

Los errores internos se registran en el servidor, pero el cliente recibe mensajes controlados. De esta manera no se exponen consultas SQL, contraseñas o detalles internos de la base de datos.

### Principios SOLID
**S — Principio de responsabilidad única**

Este principio indica que cada módulo debe tener una sola responsabilidad o motivo principal para cambiar.

En el proyecto, las responsabilidades están separadas:

- El controller maneja HTTP.
- El service contiene las reglas de negocio.
- El repository trabaja con PostgreSQL.
- El validator contiene validaciones reutilizables.
- El paquete response construye respuestas JSON.

Ejemplo de controller.go:
```
func (controller *Controller) Crear(c fiber.Ctx) error {
    var input dto.CrearSolicitudRequest

    if err := c.Bind().Body(&input); err != nil {
        return response.Error(
            c,
            fiber.StatusBadRequest,
            "INVALID_JSON",
            "el cuerpo JSON no es válido",
        )
    }

    solicitud, err := controller.service.Crear(
        c.Context(),
        input,
    )

    if err != nil {
        return responderError(c, err)
    }

    return response.JSON(
        c,
        fiber.StatusCreated,
        solicitud,
    )
}
```
**El controller recibe la petición y asigna las reglas de negocio al service. No contiene consultas SQL.**

**O — Principio abierto/cerrado**

Este principio establece que el código debe estar abierto a extensiones, pero cerrado a modificaciones innecesarias.

El servicio trabaja con la interfaz Repository:
```
type solicitudService struct {
    repository repository.Repository
}

func NewService(repository repository.Repository) Service {
    return &solicitudService{
        repository: repository,
    }
}
```
Actualmente existe una implementación para PostgreSQL, pero podría agregarse otra implementación en memoria o para otra base de datos sin cambiar la lógica principal del servicio.

La funcionalidad se amplía creando una nueva implementación de la interfaz.

**L — Principio de sustitución de Liskov**

Este principio indica que una implementación debe poder sustituir a la abstracción que representa sin cambiar el comportamiento esperado del programa.

La implementación postgresRepository cumple el contrato definido por Repository:
```
type Repository interface {
    ObtenerTodas(ctx context.Context) ([]model.SolicitudOperativa, error)
    Crear(ctx context.Context, solicitud *model.SolicitudOperativa) error
    Actualizar(ctx context.Context, solicitud *model.SolicitudOperativa) error
    Eliminar(ctx context.Context, id uuid.UUID) error
    ActualizarEstado(
        ctx context.Context,
        id uuid.UUID,
        estado model.Estado,
    ) (*model.SolicitudOperativa, error)
}
```
El servicio no necesita conocer si los datos provienen de PostgreSQL, una implementación de prueba u otro mecanismo de almacenamiento. Cualquier implementación que respete este contrato puede reemplazar a postgresRepository.

**I — Principio de segregación de interfaces**

Este principio establece que una clase o componente no debería depender de operaciones que no necesita.

La interfaz Repository contiene únicamente operaciones relacionadas con las solicitudes operativas:
```
type Repository interface {
    ObtenerTodas(ctx context.Context) ([]model.SolicitudOperativa, error)
    Crear(ctx context.Context, solicitud *model.SolicitudOperativa) error
    Actualizar(ctx context.Context, solicitud *model.SolicitudOperativa) error
    Eliminar(ctx context.Context, id uuid.UUID) error
    ActualizarEstado(
        ctx context.Context,
        id uuid.UUID,
        estado model.Estado,
    ) (*model.SolicitudOperativa, error)
}
```
No incluye operaciones relacionadas con usuarios, autenticación u otros módulos. El contrato se mantiene enfocado en una única entidad y sus casos de uso.

**D — Principio de inversión de dependencias**

Este principio establece que la lógica de negocio debe depender de abstracciones y no directamente de implementaciones concretas.

El servicio depende de repository.Repository, no de postgresRepository ni de pgxpool.Pool.
```
type solicitudService struct {
    repository repository.Repository
}
```
La implementación concreta se conecta en main.go:
```
solicitudRepository :=
    repository.NewPostgresRepository(db)

solicitudService :=
    service.NewService(solicitudRepository)

solicitudController :=
    controller.NewController(solicitudService)
```
Esto reduce el acoplamiento y permite reemplazar el repositorio sin modificar las reglas de negocio.

### Código limpio

El proyecto aplica las siguientes prácticas:

- Nombres descriptivos para funciones, estructuras y variables.
- Separación entre presentación, negocio y persistencia.
- Funciones con responsabilidades concretas.
- DTO diferentes para crear, actualizar y cambiar el estado.
- Validaciones reutilizables.
- Respuestas JSON centralizadas.
- Manejo centralizado de errores.
- Consultas SQL parametrizadas.
- Eliminación de espacios innecesarios en los textos.
- Cierre correcto del pool de conexiones.
- Apagado controlado del servidor.

**Ejemplo de los DTO separados:**
```
type CrearSolicitudRequest struct {
    Titulo          string  `json:"titulo"`
    AreaSolicitante string  `json:"area_solicitante"`
    Prioridad       int16   `json:"prioridad"`
    CostoEstimado   float64 `json:"costo_estimado"`
}

type ActualizarEstadoRequest struct {
    Estado string `json:"estado"`
}
```
**Esta separación evita utilizar una única estructura para operaciones con necesidades diferentes.**

### Pruebas con Postman

La colección se encuentra en:
```
postman/academia-api.postman_collection.json
```
La colección contiene los cinco endpoints y utiliza las siguientes variables:

- baseUrl: dirección de la API.
- solicitudId: identificador de la solicitud creada.

La petición de creación guarda automáticamente el UUID obtenido en solicitudId. Esto permite ejecutar posteriormente las solicitudes de actualización, cambio de estado y eliminación.

El orden recomendado es:

- Crear solicitud.
- Listar solicitudes.
- Actualizar solicitud.
- Actualizar estado.
- Eliminar solicitud.

### Verificación del proyecto

Para comprobar que los paquetes compilan correctamente:
```
go test ./...
```
Para ejecutar el análisis estático:
```
go vet ./...
```
Los endpoints también fueron verificados manualmente y devolvieron los códigos HTTP esperados:

- POST: 201 Created
- GET: 200 OK
- PUT: 200 OK
- PATCH: 200 OK
- DELETE: 204 No Content
- No existe: 404 Not Found

### Uso de inteligencia artificial

Se utilizó inteligencia artificial como apoyo para:

- Diseñar la estructura de la base de datos.
- Definir la arquitectura modular por capas.
- Generar y revisar los cinco endpoints.
- Aplicar validaciones y consultas parametrizadas.
- Revisar el cumplimiento de los principios SOLID.
- Preparar pruebas para los endpoints.

Las respuestas generadas no fueron aceptadas automáticamente. El código fue revisado, ajustado y probado para comprobar su funcionamiento, seguridad y organización.

La evidencia completa de los prompts, respuestas y ajustes realizados se encuentra en:

[PROMPTS.md](./PROMPTS.md)
