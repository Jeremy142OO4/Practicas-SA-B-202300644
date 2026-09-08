import asyncio

import pytest

from .apoyo import SimuladorAsincrono, ordenes


@pytest.mark.unit
def test_rechaza_carrito_vacio():
    repositorio = SimuladorAsincrono(
        obtener_carrito={"id_carrito": "c1", "id_usuario": "u1", "estado": "ACTIVO"},
        listar_detalles_carrito=[],
    )
    servicio = ordenes.OrderService(repositorio, SimuladorAsincrono(), SimuladorAsincrono())
    with pytest.raises(ordenes.CarritoVacio):
        asyncio.run(servicio.crear_orden("c1"))
