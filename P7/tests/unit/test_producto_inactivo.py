import asyncio

import pytest

from .apoyo import SimuladorAsincrono, ordenes


@pytest.mark.unit
def test_rechaza_producto_inactivo():
    catalogo = SimuladorAsincrono(obtener_producto={"activo": False, "stock": 5})
    servicio = ordenes.OrderService(SimuladorAsincrono(), catalogo, SimuladorAsincrono())
    with pytest.raises(ordenes.ProductoSinStock):
        asyncio.run(servicio.agregar_producto("u1", "p1", 1))
