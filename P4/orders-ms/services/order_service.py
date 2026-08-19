from decimal import Decimal


class CarritoNoEncontrado(Exception):
    pass


class CarritoVacio(Exception):
    pass


class ProductoSinStock(Exception):
    pass


class OrdenNoEncontrada(Exception):
    pass


class OrdenNoPendiente(Exception):
    pass


class OrderService:

    def __init__(
        self,
        repositorio,
        catalogo,
        pagos,
    ):
        self.repositorio = repositorio
        self.catalogo = catalogo
        self.pagos = pagos

    async def obtener_carrito(
        self,
        id_usuario,
    ):

        carrito = (
            await self.repositorio
            .obtener_carrito_activo(
                id_usuario
            )
        )

        if carrito is None:

            carrito = (
                await self.repositorio
                .crear_carrito(
                    id_usuario
                )
            )

        detalles = (
            await self.repositorio
            .listar_detalles_carrito(
                carrito["id_carrito"]
            )
        )

        carrito["detalles"] = detalles

        return carrito

    async def agregar_producto(
        self,
        id_usuario,
        id_videojuego_plataforma,
        cantidad,
    ):

        producto = (
            await self.catalogo
            .obtener_producto(
                id_videojuego_plataforma
            )
        )

        if not producto["activo"]:
            raise ProductoSinStock(
                "el producto no está disponible"
            )

        if producto["stock"] < cantidad:
            raise ProductoSinStock(
                "stock insuficiente"
            )

        carrito = (
            await self.obtener_carrito(
                id_usuario
            )
        )

        await self.repositorio.agregar_producto(
            carrito["id_carrito"],
            id_videojuego_plataforma,
            cantidad,
            Decimal(
                str(producto["precio"])
            ),
        )

        return await self.obtener_carrito(
            id_usuario
        )

    async def eliminar_producto(
        self,
        id_usuario,
        id_detalle_carrito,
    ):

        await self.repositorio.eliminar_producto(
            id_detalle_carrito
        )

        return await self.obtener_carrito(
            id_usuario
        )

    async def crear_orden(
        self,
        id_carrito,
    ):

        carrito = (
            await self.repositorio
            .obtener_carrito(
                id_carrito
            )
        )

        if carrito is None:
            raise CarritoNoEncontrado(
                "carrito no encontrado"
            )

        if carrito["estado"] != "ACTIVO":
            raise CarritoNoEncontrado(
                "el carrito ya fue procesado"
            )

        detalles_carrito = (
            await self.repositorio
            .listar_detalles_carrito(
                id_carrito
            )
        )

        if not detalles_carrito:
            raise CarritoVacio(
                "el carrito está vacío"
            )

        detalles_orden = []

        total = Decimal("0")

        for detalle in detalles_carrito:

            producto = (
                await self.catalogo
                .obtener_producto(
                    detalle[
                        "id_videojuego_plataforma"
                    ]
                )
            )

            if (
                not producto["activo"]
                or producto["stock"]
                < detalle["cantidad"]
            ):

                raise ProductoSinStock(
                    f"stock insuficiente para "
                    f"{producto['titulo']}"
                )

            precio = Decimal(
                str(producto["precio"])
            )

            subtotal = (
                precio
                * detalle["cantidad"]
            )

            total += subtotal

            detalles_orden.append(
                {
                    "id_videojuego_plataforma":
                        detalle[
                            "id_videojuego_plataforma"
                        ],

                    "nombre_videojuego":
                        producto["titulo"],

                    "plataforma":
                        producto["plataforma"],

                    "cantidad":
                        detalle["cantidad"],

                    "precio_unitario":
                        precio,

                    "subtotal":
                        subtotal,
                }
            )

        orden = (
            await self.repositorio.crear_orden(
                id_carrito,
                carrito["id_usuario"],
                total,
                detalles_orden,
            )
        )

        return await self.obtener_orden(
            orden["id_orden"]
        )

    async def obtener_orden(
        self,
        id_orden,
    ):

        orden = (
            await self.repositorio
            .obtener_orden(
                id_orden
            )
        )

        if orden is None:
            raise OrdenNoEncontrada(
                "orden no encontrada"
            )

        orden["detalles"] = (
            await self.repositorio
            .listar_detalles_orden(
                id_orden
            )
        )

        return orden

    async def listar_ordenes(
        self,
        id_usuario,
    ):

        ordenes = (
            await self.repositorio
            .listar_ordenes(
                id_usuario
            )
        )

        resultado = []

        for orden in ordenes:

            orden["detalles"] = (
                await self.repositorio
                .listar_detalles_orden(
                    orden["id_orden"]
                )
            )

            resultado.append(orden)

        return resultado

    async def pagar_orden(
        self,
        id_orden,
        id_metodo_pago,
    ):

        orden = await self.obtener_orden(
            id_orden
        )

        if orden["estado"] != "PENDIENTE":
            raise OrdenNoPendiente(
                "la orden no está pendiente"
            )

        resultado_pago = (
            await self.pagos.procesar_pago(
                orden["id_orden"],
                id_metodo_pago,
                orden["total"],
            )
        )

        pago = resultado_pago["pago"]

        if pago["estado"] != "APROBADO":

            await self.repositorio \
                .actualizar_estado_orden(
                    id_orden,
                    "RECHAZADA",
                )

            return await self.obtener_orden(
                id_orden
            )

        for detalle in orden["detalles"]:

            producto = (
                await self.catalogo
                .obtener_producto(
                    detalle[
                        "id_videojuego_plataforma"
                    ]
                )
            )

            nuevo_stock = (
                producto["stock"]
                - detalle["cantidad"]
            )

            if nuevo_stock < 0:
                raise ProductoSinStock(
                    "stock insuficiente"
                )

            await self.catalogo.actualizar_stock(
                detalle[
                    "id_videojuego_plataforma"
                ],
                nuevo_stock,
            )

        await self.repositorio \
            .actualizar_estado_orden(
                id_orden,
                "PAGADA",
            )

        return await self.obtener_orden(
            id_orden
        )