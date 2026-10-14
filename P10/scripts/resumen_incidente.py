"""Muestra el resultado medido guardado por observar_incidente.py."""
from datetime import datetime
import json
from pathlib import Path

data = json.loads((Path(__file__).resolve().parents[1] / "evidencias" / "09-incidente.json").read_text())


def seconds(a, b):
    return round((datetime.fromisoformat(data[b].replace("Z", "+00:00")) -
                  datetime.fromisoformat(data[a].replace("Z", "+00:00"))).total_seconds(), 3)


print("INCIDENTE P10 - marcas de tiempo reales (UTC)")
for field in ("faultStartedUTC", "alertObservedUTC", "causeFoundUTC", "analysisFinishedUTC", "abortObservedUTC", "stableResponseUTC"):
    print(field + ": " + str(data.get(field, "Aun no observado")))
print("\nAnalysisRun: " + str(data.get("analysisRun", "Aun no observado")))
for metric in data.get("analysisStatus", {}).get("metricResults", []):
    print(metric["name"] + ": " + metric["phase"] + " " + metric["measurements"][-1]["value"])
print("\nTrace ID: " + str(data.get("traceID", "Aun no observado")))
print("Servicios: " + ", ".join(data.get("traceServices", [])))
print("Causa: " + json.dumps(data.get("faultLog", {})))
if "diagnosisSeconds" in data:
    print("Recorrido alerta -> causa: " + str(data["diagnosisSeconds"]) + " segundos")
if "stableResponseUTC" in data:
    print("Recuperacion desde inicio: " + str(seconds("faultStartedUTC", "stableResponseUTC")) + " segundos")
    print("Respuesta stable: " + str(data["stableResponseMilliseconds"]) + " ms")
print("\nCaptura completa: " + str(data.get("completed", False)))
