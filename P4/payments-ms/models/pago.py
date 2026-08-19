from datetime import datetime
from decimal import Decimal
from uuid import UUID

from pydantic import BaseModel, Field


class ProcesarPagoRequest(BaseModel):
    id_orden: UUID
    id_metodo_pago: int
    monto: Decimal = Field(gt=0)


class Pago(BaseModel):
    id_pago: UUID
    id_metodo_pago: int
    id_orden: UUID
    monto: Decimal
    estado: str
    created_at: datetime
    updated_at: datetime