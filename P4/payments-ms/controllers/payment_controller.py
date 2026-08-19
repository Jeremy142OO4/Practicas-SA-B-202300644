from uuid import UUID

from fastapi import HTTPException

from models.pago import ProcesarPagoRequest
from services.payment_service import (
    MetodoPagoInactivo,
    MetodoPagoNoEncontrado,
)


class PagoController:

    def __init__(self, servicio):
        self.servicio = servicio

    async def listar_metodos_pago(self):

        return await self.servicio.listar_metodos_pago()

    async def procesar_pago(
        self,
        datos: ProcesarPagoRequest,
    ):

        try:

            resultado = await self.servicio.procesar_pago(
                datos
            )

            return {
                "mensaje": "pago procesado",
                **resultado,
            }

        except MetodoPagoNoEncontrado as error:

            raise HTTPException(
                status_code=404,
                detail=str(error),
            )

        except MetodoPagoInactivo as error:

            raise HTTPException(
                status_code=400,
                detail=str(error),
            )

        except Exception:

            raise HTTPException(
                status_code=500,
                detail="error al procesar el pago",
            )

    async def obtener_pago(
        self,
        id_pago: UUID,
    ):

        try:

            return await self.servicio.obtener_pago(
                id_pago
            )

        except ValueError as error:

            raise HTTPException(
                status_code=404,
                detail=str(error),
            )

    async def listar_pagos_orden(
        self,
        id_orden: UUID,
    ):

        return await self.servicio.listar_pagos_orden(
            id_orden
        )