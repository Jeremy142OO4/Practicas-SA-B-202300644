"""Registra tiempos y datos reales; no modifica el cluster."""
import json
from pathlib import Path
import subprocess
import time
from datetime import datetime, timezone
from urllib.parse import urlencode
from urllib.request import Request, urlopen

ROOT = Path(__file__).resolve().parents[1]
OUTPUT = ROOT / "evidencias" / "09-incidente.json"
OUTPUT.parent.mkdir(exist_ok=True)
start_ns = time.time_ns()
result = {"observerStartedUTC": datetime.now(timezone.utc).isoformat(), "events": []}
if OUTPUT.exists():
    previous = json.loads(OUTPUT.read_text(encoding="utf-8"))
    if not previous.get("completed"):
        result = previous
        start_ns = int(datetime.fromisoformat(result["observerStartedUTC"]).timestamp() * 1e9)


def now():
    return datetime.now(timezone.utc).isoformat()


def api(port, path, **params):
    url = f"http://127.0.0.1:{port}" + path
    if params:
        url += "?" + urlencode(params)
    with urlopen(Request(url, headers={"Accept": "application/json"}), timeout=10) as response:
        return json.load(response)


def kube(kind, name=""):
    args = ["kubectl", "get", kind]
    if name:
        args.append(name)
    args += ["-n", "sa-p9", "-o", "json"]
    return json.loads(subprocess.check_output(args))


def event(name, value):
    result[name] = value
    result["events"].append({"timeUTC": now(), "event": name})
    OUTPUT.write_text(json.dumps(result, indent=2), encoding="utf-8")
    print(name + ": " + json.dumps(value), flush=True)


diagnostic_started = None
body = json.dumps({"query": 'mutation { agregarProducto(idUsuario:"20230064-4000-4000-8000-000000000001", idVideojuegoPlataforma:"20230064-4000-4000-8000-000000000002", cantidad:1) { idCarrito } }'}).encode()
deadline = time.monotonic() + 600
while time.monotonic() < deadline:
    try:
        rollout = kube("rollout", "api-gateway")
        env = rollout["spec"]["template"]["spec"]["containers"][0].get("env", [])
        if not any(e.get("name") == "P10_DEMO_DELAY_MS" and e.get("value") == "1200" for e in env):
            time.sleep(2)
            continue
        pod_hash = rollout["status"].get("currentPodHash")
        if rollout["status"].get("observedGeneration") != str(rollout["metadata"]["generation"]):
            # El controlador puede devolver por unos segundos el hash anterior.
            if int(rollout["status"].get("observedGeneration", 0)) != rollout["metadata"]["generation"]:
                time.sleep(2)
                continue
        if result.get("candidateHash") != pod_hash and pod_hash:
            event("candidateHash", pod_hash)
        if "faultStartedUTC" not in result:
            logs = api(13100, "/loki/api/v1/query_range", start=str(start_ns),
                       query='{namespace="sa-p9",pod=~"api-gateway-' + pod_hash + '-.*"} |= "demo_delay"',
                       limit=1, direction="forward")
            for stream in logs["data"]["result"]:
                for timestamp, line in stream["values"]:
                    payload = json.loads(line)
                    event("faultStartedUTC", datetime.fromtimestamp(int(timestamp)/1e9, timezone.utc).isoformat())
                    event("traceID", payload["trace_id"])
                    event("faultLog", payload)
        alerts = api(19090, "/api/v1/alerts")["data"]["alerts"]
        for alert in alerts:
            if alert["labels"].get("alertname") == "ApiLatencyErrorBudgetBurn" and alert["state"] == "firing" and "faultStartedUTC" in result and "alertObservedUTC" not in result:
                event("alertObservedUTC", now())
                event("alert", alert)
                diagnostic_started = time.monotonic()
        if diagnostic_started is not None and "causeFoundUTC" not in result:
            metric = api(19090, "/api/v1/query", query='slo:canary_latency:burn_rate_1m{pod=~"api-gateway-' + pod_hash + '-.*"}')["data"]["result"]
            trace = api(23200, "/api/traces/" + result["traceID"])
            logs = api(13100, "/loki/api/v1/query_range", start=str(start_ns),
                       query='{namespace="sa-p9"} |= "' + result["traceID"] + '"', limit=10)
            found = [json.loads(line) for stream in logs["data"]["result"] for _, line in stream["values"]]
            services = sorted({attr["value"]["stringValue"] for batch in trace["batches"]
                               for attr in batch["resource"]["attributes"] if attr["key"] == "service.name"})
            if metric and found and "api-gateway" in services:
                event("burnRateAtDiagnosis", metric)
                event("traceServices", services)
                event("correlatedLogs", found)
                event("causeFoundUTC", now())
                event("diagnosisSeconds", round(time.monotonic() - diagnostic_started, 3))
        if "analysisRun" not in result:
            for run in kube("analysisruns")["items"]:
                if run["metadata"]["name"].startswith("api-gateway-" + pod_hash + "-") and run.get("status", {}).get("phase") == "Failed":
                    event("analysisRun", run["metadata"]["name"])
                    event("analysisFinishedUTC", run["status"].get("completedAt"))
                    event("analysisStatus", run["status"])
        if rollout["status"].get("abort") and "abortObservedUTC" not in result:
            event("abortObservedUTC", now())
            event("stableRS", rollout["status"].get("stableRS"))
        if "abortObservedUTC" in result and "stableResponseUTC" not in result:
            started = time.monotonic()
            with urlopen(Request("http://127.0.0.1:27000/orders/graphql", data=body,
                                 headers={"Content-Type": "application/json"}), timeout=5) as response:
                response.read()
                if response.status == 200 and time.monotonic() - started < 0.5:
                    event("stableResponseUTC", now())
                    event("stableResponseMilliseconds", round((time.monotonic() - started)*1000, 3))
        if all(k in result for k in ("faultStartedUTC", "causeFoundUTC", "analysisRun", "stableResponseUTC")):
            event("completed", True)
            break
    except Exception as exc:
        print("Consulta temporalmente no disponible: " + str(exc), flush=True)
    time.sleep(2)
else:
    event("completed", False)
