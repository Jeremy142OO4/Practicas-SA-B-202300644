from datetime import datetime
from uuid import UUID

from pydantic import BaseModel


class Carrito(BaseModel):
    id_carrito: UUID
    id_usuario: UUID
    estado: str
    created_at: datetime
    updated_at: datetime