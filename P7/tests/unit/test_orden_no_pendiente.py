import asyncio

import pytest

from .apoyo import SimuladorAsincrono, ordenes


@pytest.mark.unit
def test_rechaza_pago_de_orden_no_pendiente():
    repositorio = SimuladorAsincrono(
        obtener_orden={"id_orden": "o1", "estado": "PAGADA"},
        listar_detalles_orden=[],
    )
    servicio = ordenes.OrderService(repositorio, SimuladorAsincrono(), SimuladorAsincrono())
    with pytest.raises(ordenes.OrdenNoPendiente):
        asyncio.run(servicio.pagar_orden("o1", "m1"))
