import os

import httpx


class PaymentClient:

    def __init__(self):
        self.url = os.getenv(
            "PAYMENTS_URL",
            "http://payments-ms:4004",
        )

    async def procesar_pago(
        self,
        id_orden,
        id_metodo_pago,
        monto,
    ):

        async with httpx.AsyncClient(
            timeout=10.0
        ) as cliente:

            respuesta = await cliente.post(
                f"{self.url}/pagos/procesar",
                json={
                    "id_orden": str(id_orden),
                    "id_metodo_pago":
                        id_metodo_pago,
                    "monto": float(monto),
                },
            )

            respuesta.raise_for_status()

            return respuesta.json()