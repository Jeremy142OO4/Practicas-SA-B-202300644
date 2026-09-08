from uuid import UUID

from pydantic import BaseModel


class ResultadoPagoRequest(BaseModel):
    id_orden: UUID
    estado: str
