import asyncio
from decimal import Decimal
from types import SimpleNamespace

import pytest

from .apoyo import SimuladorAsincrono, pagos


@pytest.mark.unit
def test_pago_idempotente_retorna_existente():
    pago_existente = {"id_pago": "pay1", "estado": "APROBADO"}
    repositorio = SimuladorAsincrono(obtener_pago_por_orden=pago_existente)
    servicio = pagos.PagoService(repositorio)
    resultado = asyncio.run(
        servicio.procesar_pago(SimpleNamespace(id_orden="o1", id_metodo_pago="m1", monto=Decimal("10")))
    )
    assert resultado == {"pago": pago_existente, "transaccion": None}
