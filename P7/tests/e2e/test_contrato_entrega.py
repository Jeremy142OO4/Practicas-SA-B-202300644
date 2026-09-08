from pathlib import Path

import pytest


RAIZ = Path(__file__).resolve().parents[3]


@pytest.mark.e2e
def test_flujo_de_entrega_contiene_seis_imagenes_y_despliegue():
    flujo_trabajo = (RAIZ / ".github/workflows/p7-ci-cd.yml").read_text(encoding="utf-8")
    servicios = ["api-gateway", "auth-ms", "catalog-ms", "orders-ms", "payments-ms"]
    for servicio in servicios:
        assert servicio in flujo_trabajo
        assert (RAIZ / "P4" / servicio / "Dockerfile.prod").is_file()
    assert "scheduled-jobs" in flujo_trabajo
    assert (RAIZ / "P5" / "cronjobs" / "Dockerfile.prod").is_file()
    assert "DOCKERHUB_TOKEN" in flujo_trabajo
    assert "helm upgrade --install" in flujo_trabajo
