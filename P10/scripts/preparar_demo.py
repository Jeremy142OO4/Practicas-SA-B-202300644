"""Consultas de lectura acotadas para obtener muestras reales antes de demostrar P10."""
import argparse
import json
import signal
import subprocess
import time
from urllib.request import Request, urlopen


def main():
    def stop(_signum, _frame):
        raise KeyboardInterrupt

    signal.signal(signal.SIGINT, stop)
    signal.signal(signal.SIGTERM, stop)
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--segundos", type=int, default=600)
    args = parser.parse_args()
    if not 30 <= args.segundos <= 1800:
        parser.error("La duracion debe estar entre 30 y 1800 segundos")
    context = subprocess.check_output(
        ["kubectl", "config", "current-context"], text=True
    ).strip()
    if context != "gke_sa-202300644_us-central1-a_p9-202300644":
        raise SystemExit("Contexto incorrecto: no se enviaron consultas")

    services = {"api-gateway-stable": 26210, "api-gateway-canary": 27210}
    processes = {}
    body = json.dumps({"query": '{ ordenes(idUsuario:"20230064-4000-4000-8000-000000000001") { __typename } }'}).encode()
    deadline = time.monotonic() + args.segundos
    last_report = 0
    successes = retries = 0
    try:
        while time.monotonic() < deadline:
            for service, port in services.items():
                try:
                    resource = json.loads(subprocess.check_output(
                        ["kubectl", "get", "svc", service, "-n", "sa-p9", "-o", "json"]
                    ))
                    pod_hash = resource["spec"]["selector"].get("rollouts-pod-template-hash")
                    previous = processes.get(service)
                    if previous is None or previous[0] != pod_hash or previous[1].poll() is not None:
                        if previous and previous[1].poll() is None:
                            previous[1].terminate()
                            previous[1].wait(timeout=10)
                        process = subprocess.Popen(
                            ["kubectl", "port-forward", "-n", "sa-p9", "svc/" + service, f"{port}:4000"],
                            stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL,
                        )
                        processes[service] = (pod_hash, process)
                        time.sleep(1)
                        if process.poll() is not None:
                            raise RuntimeError(f"No se pudo abrir el puerto {port}")
                        print(f"{service}: {pod_hash}", flush=True)
                    request = Request(
                        f"http://127.0.0.1:{port}/orders/graphql", data=body,
                        headers={"Content-Type": "application/json"},
                    )
                    with urlopen(request, timeout=5) as response:
                        result = json.load(response)
                    if result.get("errors") or "data" not in result:
                        raise RuntimeError("La consulta GraphQL no fue exitosa")
                    successes += 1
                except Exception as exc:
                    retries += 1
                    print(f"Reintento {service}: {exc}", flush=True)
            if time.monotonic() - last_report >= 30:
                print(f"Lecturas correctas: {successes}; reintentos: {retries}", flush=True)
                last_report = time.monotonic()
            time.sleep(3)
    except KeyboardInterrupt:
        print("Fin solicitado por el usuario", flush=True)
    finally:
        for _, process in processes.values():
            if process.poll() is None:
                process.terminate()
                process.wait(timeout=10)
    if successes == 0:
        raise SystemExit("No se obtuvo ninguna solicitud exitosa")


if __name__ == "__main__":
    main()
