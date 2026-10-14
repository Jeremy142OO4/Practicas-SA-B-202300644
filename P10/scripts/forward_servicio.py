"""Reconecta el port-forward cuando Argo cambia el pod seleccionado."""
import json
import subprocess
import sys
import time

service, port = sys.argv[1:]
current_hash = None
process = None
try:
    while True:
        value = json.loads(subprocess.check_output([
            "kubectl", "get", "svc", service, "-n", "sa-p9", "-o", "json"]))
        candidate_hash = value["spec"]["selector"].get("rollouts-pod-template-hash")
        if process is None or process.poll() is not None or candidate_hash != current_hash:
            if process is not None and process.poll() is None:
                process.terminate()
                process.wait(timeout=5)
            process = subprocess.Popen([
                "kubectl", "port-forward", "-n", "sa-p9", "svc/" + service,
                f"{port}:4000"], stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)
            current_hash = candidate_hash
            print(f"{service}: pod hash {candidate_hash}", flush=True)
        time.sleep(2)
finally:
    if process is not None:
        process.terminate()
