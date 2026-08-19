import strawberry


@strawberry.type
class DetalleCarritoTipo:
    id_detalle_carrito: strawberry.ID
    id_videojuego_plataforma: strawberry.ID
    cantidad: int
    precio_unitario: float


@strawberry.type
class CarritoTipo:
    id_carrito: strawberry.ID
    id_usuario: strawberry.ID
    estado: str
    detalles: list[DetalleCarritoTipo]


@strawberry.type
class DetalleOrdenTipo:
    id_detalle_orden: strawberry.ID
    id_videojuego_plataforma: strawberry.ID
    nombre_videojuego: str
    plataforma: str
    cantidad: int
    precio_unitario: float
    subtotal: float


@strawberry.type
class OrdenTipo:
    id_orden: strawberry.ID
    id_carrito: strawberry.ID
    id_usuario: strawberry.ID
    total: float
    estado: str
    detalles: list[DetalleOrdenTipo]


def convertir_carrito(datos):

    detalles = [
        DetalleCarritoTipo(
            id_detalle_carrito=str(
                detalle[
                    "id_detalle_carrito"
                ]
            ),
            id_videojuego_plataforma=str(
                detalle[
                    "id_videojuego_plataforma"
                ]
            ),
            cantidad=detalle["cantidad"],
            precio_unitario=float(
                detalle["precio_unitario"]
            ),
        )
        for detalle in datos["detalles"]
    ]

    return CarritoTipo(
        id_carrito=str(
            datos["id_carrito"]
        ),
        id_usuario=str(
            datos["id_usuario"]
        ),
        estado=datos["estado"],
        detalles=detalles,
    )


def convertir_orden(datos):

    detalles = [
        DetalleOrdenTipo(
            id_detalle_orden=str(
                detalle[
                    "id_detalle_orden"
                ]
            ),
            id_videojuego_plataforma=str(
                detalle[
                    "id_videojuego_plataforma"
                ]
            ),
            nombre_videojuego=detalle[
                "nombre_videojuego"
            ],
            plataforma=detalle[
                "plataforma"
            ],
            cantidad=detalle[
                "cantidad"
            ],
            precio_unitario=float(
                detalle[
                    "precio_unitario"
                ]
            ),
            subtotal=float(
                detalle["subtotal"]
            ),
        )
        for detalle in datos["detalles"]
    ]

    return OrdenTipo(
        id_orden=str(
            datos["id_orden"]
        ),
        id_carrito=str(
            datos["id_carrito"]
        ),
        id_usuario=str(
            datos["id_usuario"]
        ),
        total=float(
            datos["total"]
        ),
        estado=datos["estado"],
        detalles=detalles,
    )