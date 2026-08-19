from decimal import Decimal
from uuid import UUID

from pydantic import BaseModel


class DetalleOrden(BaseModel):
    id_detalle_orden: UUID
    id_orden: UUID
    id_videojuego_plataforma: UUID
    nombre_videojuego: str
    plataforma: str
    cantidad: int
    precio_unitario: Decimal
    subtotal: Decimal