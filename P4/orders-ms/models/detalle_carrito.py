from decimal import Decimal
from uuid import UUID

from pydantic import BaseModel, Field


class DetalleCarrito(BaseModel):
    id_detalle_carrito: UUID
    id_carrito: UUID
    id_videojuego_plataforma: UUID
    cantidad: int
    precio_unitario: Decimal


class AgregarProductoInput(BaseModel):
    id_usuario: UUID
    id_videojuego_plataforma: UUID
    cantidad: int = Field(gt=0)