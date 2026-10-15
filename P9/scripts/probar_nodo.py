"""Drenaje selectivo de pedidos; no simula la perdida total de una maquina."""
import datetime
import json
import signal
import subprocess
import time
from urllib.request import Request, urlopen


def command(*args):
    return subprocess.check_output(["kubectl", *args], text=True)


def utc():
    return datetime.datetime.now(datetime.timezone.utc).isoformat()


def main():
    def stop(_signal, _frame):
        raise KeyboardInterrupt

    signal.signal(signal.SIGTERM, stop)
    signal.signal(signal.SIGINT, stop)
    if command("config", "current-context").strip() != "gke_sa-202300644_us-central1-a_p9-202300644":
        raise SystemExit("Contexto incorrecto")
    pods = json.loads(command("get", "pods", "-n", "sa-p9", "-l", "app.kubernetes.io/name=orders-ms", "-o", "json"))["items"]
    ready = [p for p in pods if any(c["type"] == "Ready" and c["status"] == "True" for c in p["status"].get("conditions", []))]
    if len(ready) < 2:
        raise SystemExit("Se necesitan dos replicas Ready antes de probar")
    victim = ready[0]
    node = victim["spec"]["nodeName"]
    original_uid = victim["metadata"]["uid"]
    node_info = json.loads(command("get", "node", node, "-o", "json"))
    if node_info["spec"].get("unschedulable"):
        raise SystemExit("El nodo ya estaba cordonado; no se modifica")
    forward = subprocess.Popen(
        ["kubectl", "port-forward", "-n", "sa-p9", "svc/api-gateway-stable", "26909:4000"],
        stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL,
    )
    drain = None
    cordoned = False
    samples = []
    body = json.dumps({"query": '{ ordenes(idUsuario:"20230064-4000-4000-8000-000000000001") { __typename } }'}).encode()
    try:
        time.sleep(2)
        if forward.poll() is not None:
            raise RuntimeError("No se pudo abrir el puerto 26909")
        request = Request("http://127.0.0.1:26909/orders/graphql", data=body, headers={"Content-Type": "application/json"})
        with urlopen(request, timeout=5) as response:
            if json.load(response).get("errors"):
                raise RuntimeError("La consulta inicial no fue exitosa")
        print("Inicio UTC:", utc(), flush=True)
        print("Nodo:", node, "pod inicial:", victim["metadata"]["name"], flush=True)
        print("Alcance: solo pedidos; los demas pods del nodo no se desalojan", flush=True)
        print(command("cordon", node).strip(), flush=True)
        cordoned = True
        drain = subprocess.Popen([
            "kubectl", "drain", node, "--pod-selector=app.kubernetes.io/name=orders-ms",
            "--ignore-daemonsets", "--timeout=180s",
        ])
        deadline = time.monotonic() + 200
        while time.monotonic() < deadline:
            try:
                with urlopen(request, timeout=5) as response:
                    result = json.load(response)
                    ok = response.status == 200 and not result.get("errors")
                samples.append(bool(ok))
                if not ok:
                    print("Consulta GraphQL fallida UTC:", utc(), flush=True)
            except Exception as exc:
                samples.append(False)
                print("Solicitud fallida UTC:", utc(), type(exc).__name__, str(exc)[:160], flush=True)
            current = json.loads(command("get", "pods", "-n", "sa-p9", "-l", "app.kubernetes.io/name=orders-ms", "-o", "json"))["items"]
            current_ready = [p for p in current if any(c["type"] == "Ready" and c["status"] == "True" for c in p["status"].get("conditions", []))]
            replaced = all(p["metadata"]["uid"] != original_uid for p in current)
            if replaced and len(current_ready) >= 2 and len(samples) >= 30 and drain.poll() == 0:
                print("Nuevos pods:", [(p["metadata"]["name"], p["spec"]["nodeName"]) for p in current_ready], flush=True)
                break
            if drain.poll() not in (None, 0):
                raise RuntimeError("El drenaje fallo")
            time.sleep(1)
        else:
            raise RuntimeError("No se completo la reprogramacion dentro del plazo")
        print("Fin UTC:", utc(), flush=True)
        print("Consultas:", len(samples), "correctas:", sum(samples), "fallidas:", len(samples) - sum(samples), flush=True)
        if not all(samples):
            raise RuntimeError("Hubo solicitudes fallidas; se conserva el resultado real")
    finally:
        if drain and drain.poll() is None:
            drain.terminate()
            drain.wait(timeout=10)
        if cordoned:
            print(command("uncordon", node).strip(), flush=True)
        if forward.poll() is None:
            forward.terminate()
            forward.wait(timeout=10)


if __name__ == "__main__":
    main()
