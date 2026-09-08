import os

import asyncpg
from dotenv import load_dotenv

load_dotenv()


async def conectar():
    url_base_datos = os.environ.get("DATABASE_URL")
    if url_base_datos:
        return await asyncpg.connect(dsn=url_base_datos)
    return await asyncpg.connect(
        host=os.environ["DB_HOST"],
        port=int(os.getenv("DB_PORT", "5432")),
        user=os.environ["DB_USER"],
        password=os.environ["DB_PASSWORD"],
        database=os.environ["DB_NAME"],
    )


async def preparar_esquema(conexion):
    await conexion.execute(
        """
        CREATE SCHEMA IF NOT EXISTS scheduled_jobs;

        CREATE TABLE IF NOT EXISTS
        scheduled_jobs.ejecuciones (
            id_ejecucion BIGSERIAL PRIMARY KEY,
            fecha_hora TIMESTAMPTZ NOT NULL,
            zona_horaria VARCHAR(20) NOT NULL,
            carne VARCHAR(20) NOT NULL
        );

        CREATE TABLE IF NOT EXISTS
        scheduled_jobs.resumenes_hora (
            id_resumen BIGSERIAL PRIMARY KEY,
            event_id UUID NOT NULL,
            hora TIMESTAMP NOT NULL,
            cantidad INTEGER NOT NULL,
            recibido_at TIMESTAMPTZ NOT NULL
                DEFAULT CURRENT_TIMESTAMP,
            CONSTRAINT unique_evento_hora
                UNIQUE (event_id, hora)
        );
        """
    )
