# Base de datos PostgreSQL local

La base de datos se ejecuta en Docker y se inicializa automáticamente mediante
`init.sql` la primera vez que se crea el volumen con nombre.

## Usuarios de prueba

| Rol | Correo | Contraseña |
| --- | --- | --- |
| ADMIN | `admin@practica.local` | `Admin123!` |
| CLIENTE | `cliente@practica.local` | `Cliente123!` |

Los nombres, correos y contraseñas se almacenan como texto cifrado con
AES-256-GCM. La unicidad del correo se controla mediante un índice ciego
HMAC-SHA256. Las claves de desarrollo que permiten descifrar estos registros se
encuentran definidas en `backend/.env.example` y nunca deben reutilizarse en
producción.

## Comandos

Ejecutar desde la raíz del proyecto:

```bash
docker compose up -d --build
docker compose ps
docker compose down
```

Para crear nuevamente el esquema y los datos iniciales desde cero, se debe
eliminar el volumen de desarrollo e iniciar los servicios otra vez:

```bash
docker compose down -v
docker compose up -d --build
```
