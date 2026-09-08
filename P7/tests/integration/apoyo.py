from ..conftest import cargar_modulo


ordenes = cargar_modulo("servicio_ordenes_integracion", "P4/orders-ms/services/order_service.py")


class RepositorioEnMemoria:
    def __init__(self):
        self.carrito = {"id_carrito": "c1", "id_usuario": "u1", "estado": "ACTIVO"}
        self.detalles_carrito = [
            {"id_videojuego_plataforma": "p1", "cantidad": 2},
            {"id_videojuego_plataforma": "p2", "cantidad": 1},
        ]
        self.orden = None

    async def obtener_carrito(self, _):
        return dict(self.carrito)

    async def listar_detalles_carrito(self, _):
        return list(self.detalles_carrito)

    async def crear_orden(self, _, id_usuario, total, detalles):
        self.orden = {
            "id_orden": "o1", "id_usuario": id_usuario, "total": total, "estado": "PENDIENTE"
        }
        self.detalles_orden = detalles
        return self.orden

    async def obtener_orden(self, _):
        return dict(self.orden)

    async def listar_detalles_orden(self, _):
        return list(self.detalles_orden)

    async def actualizar_estado_orden(self, _, estado):
        self.orden["estado"] = estado


class CatalogoEnMemoria:
    def __init__(self):
        self.productos = {
            "p1": {"titulo": "Juego A", "plataforma": "PC", "precio": "25.50", "stock": 5, "activo": True},
            "p2": {"titulo": "Juego B", "plataforma": "PS5", "precio": "10.00", "stock": 3, "activo": True},
        }

    async def obtener_producto(self, id_producto):
        return dict(self.productos[id_producto])

    async def actualizar_stock(self, id_producto, stock):
        self.productos[id_producto]["stock"] = stock


class PublicadorPagos:
    def __init__(self):
        self.mensajes = []

    async def publicar_pago(self, *mensaje):
        self.mensajes.append(mensaje)
