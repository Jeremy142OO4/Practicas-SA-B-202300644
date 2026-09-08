import asyncio
from decimal import Decimal

import pytest

from .apoyo import CatalogoEnMemoria, PublicadorPagos, RepositorioEnMemoria, ordenes


@pytest.mark.integration
def test_carrito_catalogo_generan_orden_con_total_correcto():
    repositorio, catalogo = RepositorioEnMemoria(), CatalogoEnMemoria()
    servicio = ordenes.OrderService(repositorio, catalogo, PublicadorPagos())
    resultado = asyncio.run(servicio.crear_orden("c1"))
    assert resultado["total"] == Decimal("61.00")
    assert len(resultado["detalles"]) == 2
