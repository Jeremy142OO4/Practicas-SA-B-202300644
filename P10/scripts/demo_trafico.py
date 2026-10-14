"""Solicitudes de negocio de baja frecuencia para observar estable y candidata."""
import json
import time
from urllib.request import Request, urlopen

body = json.dumps({"query": 'mutation { agregarProducto(idUsuario:"20230064-4000-4000-8000-000000000001", idVideojuegoPlataforma:"20230064-4000-4000-8000-000000000002", cantidad:1) { idCarrito } }'}).encode()
deadline = time.monotonic() + 900
while time.monotonic() < deadline:
    for port in (27000, 26000):
        try:
            req = Request(f"http://127.0.0.1:{port}/orders/graphql", data=body,
                          headers={"Content-Type": "application/json"})
            with urlopen(req, timeout=5) as response:
                response.read()
        except Exception:
            # Durante un cambio de pod, el port-forward puede reconectarse.
            pass
    time.sleep(1)
