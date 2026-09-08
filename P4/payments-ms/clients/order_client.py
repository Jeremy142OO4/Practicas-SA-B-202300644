import os

import httpx


class OrderClient:

    def __init__(self):
        self.url = os.getenv(
            "ORDERS_URL",
            "http://orders-ms:4003",
        )

    async def actualizar_estado_pago(
        self,
        id_orden,
        estado,
    ):
        async with httpx.AsyncClient(
            timeout=10.0
        ) as cliente:
            respuesta = await cliente.post(
                f"{self.url}/pagos/resultado",
                json={
                    "id_orden": str(id_orden),
                    "estado": estado,
                },
            )
            respuesta.raise_for_status()
            return respuesta.json()
