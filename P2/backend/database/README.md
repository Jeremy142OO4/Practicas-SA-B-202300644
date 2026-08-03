# Local PostgreSQL database

The database runs in Docker and is initialized automatically from `init.sql` the first
time its named volume is created.

## Test users

| Role | Email | Password |
| --- | --- | --- |
| ADMIN | `admin@practica.local` | `Admin123!` |
| CLIENTE | `cliente@practica.local` | `Cliente123!` |

Names, emails, and passwords are stored as AES-256-GCM ciphertext. Email uniqueness is
enforced through an HMAC-SHA256 blind index. The development keys that decrypt these
records are defined in `backend/.env.example` and must never be reused in production.

## Commands

From the project root:

```bash
docker compose up -d --build
docker compose ps
docker compose down
```

To recreate the schema and seed data from scratch, remove the development volume and
start the service again:

```bash
docker compose down -v
docker compose up -d --build
```
