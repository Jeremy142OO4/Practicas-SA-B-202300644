import asyncio
import json
import os
from datetime import datetime, timedelta, timezone
from uuid import uuid4

import aio_pika

from database import conectar, preparar_esquema


async def main():
    conexion = await conectar()
    try:
        await preparar_esquema(conexion)
        filas = await conexion.fetch(
            """
            SELECT
                date_trunc(
                    'hour',
                    fecha_hora AT TIME ZONE
                        'America/Guatemala'
                ) AS hora,
                COUNT(*)::INTEGER AS cantidad
            FROM scheduled_jobs.ejecuciones
            GROUP BY hora
            ORDER BY hora
            """
        )
    finally:
        await conexion.close()

    evento = {
        "event_id": str(uuid4()),
        "student_id": os.environ["STUDENT_ID"],
        "generated_at": datetime.now(
            timezone(timedelta(hours=-6))
        ).isoformat(),
        "summaries": [
            {
                "hora": fila["hora"].isoformat(),
                "cantidad": fila["cantidad"],
            }
            for fila in filas
        ],
    }

    conexion_rabbit = await aio_pika.connect_robust(
        os.environ["RABBITMQ_URL"]
    )
    try:
        canal = await conexion_rabbit.channel(
            publisher_confirms=True
        )
        cola = os.getenv(
            "SUMMARY_QUEUE",
            "scheduled.summary",
        )
        await canal.declare_queue(cola, durable=True)
        await canal.default_exchange.publish(
            aio_pika.Message(
                body=json.dumps(evento).encode("utf-8"),
                content_type="application/json",
                delivery_mode=(
                    aio_pika.DeliveryMode.PERSISTENT
                ),
                message_id=evento["event_id"],
            ),
            routing_key=cola,
        )
        print(
            "Resumen publicado:",
            evento["event_id"],
            len(evento["summaries"]),
        )
    finally:
        await conexion_rabbit.close()


if __name__ == "__main__":
    asyncio.run(main())
