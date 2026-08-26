import asyncio
import os
from datetime import datetime, timedelta, timezone

from database import conectar, preparar_esquema


async def main():
    conexion = await conectar()
    try:
        await preparar_esquema(conexion)
        fecha_hora = datetime.now(
            timezone(timedelta(hours=-6))
        )
        await conexion.execute(
            """
            INSERT INTO scheduled_jobs.ejecuciones (
                fecha_hora,
                zona_horaria,
                carne
            )
            VALUES ($1, 'GMT-6', $2)
            """,
            fecha_hora,
            os.environ["STUDENT_ID"],
        )
        print(
            "Ejecucion registrada:",
            fecha_hora.isoformat(),
            os.environ["STUDENT_ID"],
        )
    finally:
        await conexion.close()


if __name__ == "__main__":
    asyncio.run(main())
