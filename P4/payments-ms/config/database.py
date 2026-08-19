import os

import asyncpg


async def conectar_base_datos():
    pool = await asyncpg.create_pool(
        host=os.getenv("DB_HOST"),
        port=int(os.getenv("DB_PORT", "5432")),
        user=os.getenv("DB_USER"),
        password=os.getenv("DB_PASSWORD"),
        database=os.getenv("DB_NAME"),
        min_size=1,
        max_size=10,
    )

    return pool


async def cerrar_base_datos(pool):
    if pool:
        await pool.close()