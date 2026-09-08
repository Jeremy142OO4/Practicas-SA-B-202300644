import asyncio
from decimal import Decimal
from types import SimpleNamespace

import pytest

from .apoyo import SimuladorAsincrono, pagos


@pytest.mark.unit
def test_rechaza_metodo_de_pago_inactivo():
    repositorio = SimuladorAsincrono(obtener_pago_por_orden=None, obtener_metodo_pago={"activo": False})
    servicio = pagos.PagoService(repositorio)
    with pytest.raises(pagos.MetodoPagoInactivo):
        asyncio.run(
            servicio.procesar_pago(SimpleNamespace(id_orden="o1", id_metodo_pago="m1", monto=Decimal("10")))
        )
