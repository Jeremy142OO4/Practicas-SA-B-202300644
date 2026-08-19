from uuid import uuid4


class MetodoPagoNoEncontrado(Exception):
    pass


class MetodoPagoInactivo(Exception):
    pass


class PagoService:

    def __init__(self, repositorio):
        self.repositorio = repositorio

    async def listar_metodos_pago(self):
        return await self.repositorio.listar_metodos_pago()

    async def procesar_pago(self, datos):

        metodo = await self.repositorio.obtener_metodo_pago(
            datos.id_metodo_pago
        )

        if metodo is None:
            raise MetodoPagoNoEncontrado(
                "método de pago no encontrado"
            )

        if not metodo["activo"]:
            raise MetodoPagoInactivo(
                "el método de pago está inactivo"
            )

        codigo_transaccion = (
            "TX-" + uuid4().hex[:12].upper()
        )

        estado = "APROBADO"

        mensaje = "Pago procesado correctamente"

        resultado = (
            await self.repositorio.guardar_pago_con_transaccion(
                id_orden=datos.id_orden,
                id_metodo_pago=datos.id_metodo_pago,
                monto=datos.monto,
                codigo_transaccion=codigo_transaccion,
                estado=estado,
                mensaje=mensaje,
            )
        )

        return resultado

    async def obtener_pago(self, id_pago):

        pago = await self.repositorio.obtener_pago(
            id_pago
        )

        if pago is None:
            raise ValueError(
                "pago no encontrado"
            )

        return pago

    async def listar_pagos_orden(self, id_orden):

        return await self.repositorio.listar_pagos_orden(
            id_orden
        )