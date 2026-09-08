import asyncio
from decimal import Decimal

import pytest

from .apoyo import CatalogoEnMemoria, PublicadorPagos, RepositorioEnMemoria, ordenes


@pytest.mark.integration
def test_orden_publica_pago_y_cambia_estado():
    repositorio, catalogo, publicador = RepositorioEnMemoria(), CatalogoEnMemoria(), PublicadorPagos()
    servicio = ordenes.OrderService(repositorio, catalogo, publicador)
    asyncio.run(servicio.crear_orden("c1"))
    resultado = asyncio.run(servicio.pagar_orden("o1", "m1"))
    assert resultado["estado"] == "PAGO_EN_PROCESO"
    assert publicador.mensajes == [("o1", "m1", Decimal("61.00"))]
