class OrderController:

    def __init__(self, servicio):
        self.servicio = servicio

    async def obtener_carrito(
        self,
        id_usuario,
    ):
        return await self.servicio.obtener_carrito(
            id_usuario
        )

    async def agregar_producto(
        self,
        id_usuario,
        id_videojuego_plataforma,
        cantidad,
    ):
        return await self.servicio.agregar_producto(
            id_usuario,
            id_videojuego_plataforma,
            cantidad,
        )

    async def eliminar_producto(
        self,
        id_usuario,
        id_detalle_carrito,
    ):
        return await self.servicio.eliminar_producto(
            id_usuario,
            id_detalle_carrito,
        )

    async def crear_orden(
        self,
        id_carrito,
    ):
        return await self.servicio.crear_orden(
            id_carrito
        )

    async def obtener_orden(
        self,
        id_orden,
    ):
        return await self.servicio.obtener_orden(
            id_orden
        )

    async def listar_ordenes(
        self,
        id_usuario,
    ):
        return await self.servicio.listar_ordenes(
            id_usuario
        )

    async def pagar_orden(
        self,
        id_orden,
        id_metodo_pago,
    ):
        return await self.servicio.pagar_orden(
            id_orden,
            id_metodo_pago,
        )