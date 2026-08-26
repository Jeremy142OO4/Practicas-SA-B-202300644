import asyncio
import json
import os
from datetime import datetime
from uuid import UUID

import aio_pika

from database import conectar, preparar_esquema


async def almacenar_evento(evento):
    event_id = UUID(evento["event_id"])
    conexion = await conectar()
    try:
        await preparar_esquema(conexion)
        async with conexion.transaction():
            for resumen in evento["summaries"]:
                await conexion.execute(
                    """
                    INSERT INTO
                    scheduled_jobs.resumenes_hora (
                        event_id,
                        hora,
                        cantidad
                    )
                    VALUES ($1, $2, $3)
                    ON CONFLICT (event_id, hora)
                    DO NOTHING
                    """,
                    event_id,
                    datetime.fromisoformat(
                        resumen["hora"]
                    ),
                    int(resumen["cantidad"]),
                )
    finally:
        await conexion.close()


async def main():
    conexion_rabbit = await aio_pika.connect_robust(
        os.environ["RABBITMQ_URL"]
    )
    canal = await conexion_rabbit.channel()
    await canal.set_qos(prefetch_count=1)
    cola = await canal.declare_queue(
        os.getenv("SUMMARY_QUEUE", "scheduled.summary"),
        durable=True,
    )

    print("Consumidor de resumenes iniciado")

    async with cola.iterator() as mensajes:
        async for mensaje in mensajes:
            async with mensaje.process(requeue=True):
                evento = json.loads(
                    mensaje.body.decode("utf-8")
                )
                await almacenar_evento(evento)
                print(
                    "Resumen almacenado:",
                    evento["event_id"],
                )


if __name__ == "__main__":
    asyncio.run(main())
