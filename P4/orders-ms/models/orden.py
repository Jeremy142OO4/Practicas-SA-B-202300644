from datetime import datetime
from decimal import Decimal
from uuid import UUID

from pydantic import BaseModel


class Orden(BaseModel):
    id_orden: UUID
    id_carrito: UUID
    id_usuario: UUID
    total: Decimal
    estado: str
    created_at: datetime
    updated_at: datetime