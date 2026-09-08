import asyncio

import pytest

from .apoyo import SimuladorAsincrono, ordenes


@pytest.mark.unit
def test_crea_carrito_si_no_existe():
    repositorio = SimuladorAsincrono(
        obtener_carrito_activo=None,
        crear_carrito={"id_carrito": "c1"},
        listar_detalles_carrito=[],
    )
    servicio = ordenes.OrderService(repositorio, SimuladorAsincrono(), SimuladorAsincrono())
    resultado = asyncio.run(servicio.obtener_carrito("u1"))
    assert resultado == {"id_carrito": "c1", "detalles": []}
