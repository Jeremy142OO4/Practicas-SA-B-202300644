from pydantic import BaseModel


class MetodoPago(BaseModel):
    id_metodo_pago: int
    nombre: str
    descripcion: str | None = None
    activo: bool