import asyncio

import pytest

from .apoyo import SimuladorAsincrono, ordenes


@pytest.mark.unit
def test_rechaza_cantidad_superior_al_stock():
    catalogo = SimuladorAsincrono(obtener_producto={"activo": True, "stock": 1})
    servicio = ordenes.OrderService(SimuladorAsincrono(), catalogo, SimuladorAsincrono())
    with pytest.raises(ordenes.ProductoSinStock):
        asyncio.run(servicio.agregar_producto("u1", "p1", 2))
