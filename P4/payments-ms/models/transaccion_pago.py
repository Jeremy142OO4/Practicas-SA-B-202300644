from datetime import datetime
from uuid import UUID

from pydantic import BaseModel


class TransaccionPago(BaseModel):
    id_transaccion: UUID
    id_pago: UUID
    codigo_transaccion: str
    estado: str
    mensaje: str
    fecha_transaccion: datetime