"""Consulta las APIs reales a traves de los port-forward de la laptop."""
import json
from urllib.parse import urlencode
from urllib.request import urlopen


def consulta(base, path, **params):
    url = base + path + ("?" + urlencode(params) if params else "")
    with urlopen(url, timeout=20) as response:
        return json.load(response)


print("PROMETHEUS - solicitudes observadas por servicio")
query = 'sum by (service) (http_server_requests_total{service=~"api-gateway|catalog-ms|orders-ms"})'
for item in consulta("http://127.0.0.1:19090", "/api/v1/query", query=query)["data"]["result"]:
    print(f'  {item["metric"]["service"]}: {item["value"][1]}')

print("\nGRAFANA - estado de la base de datos")
print(" ", consulta("http://127.0.0.1:23000", "/api/health")["database"])

print("\nTEMPO - trazas recibidas")
traces = consulta("http://127.0.0.1:23200", "/api/search",
                  q='{ name = "/orders/graphql" }', limit=1)
trace_id = ""
for item in traces.get("traces", []):
    print(f'  {item["rootServiceName"]} {item["rootTraceName"]}')
    print(f'  trace_id: {item["traceID"]}')
    trace_id = item["traceID"]
    print("  servicios: " + ", ".join(sorted(item.get("serviceStats", {}))))

print("\nLOKI - logs JSON con trace_id")
logs = consulta("http://127.0.0.1:13100", "/loki/api/v1/query_range",
                query='{namespace="sa-p9"} |= "trace_id"' + (f' |= "{trace_id}"' if trace_id else ''), limit=3)
for stream in logs["data"]["result"]:
    for _, line in stream["values"]:
        payload = json.loads(line[line.index("{"):])
        print(f'  {payload["service"]} {payload["route"]} status={payload["status"]}')
        print(f'  trace_id: {payload["trace_id"]}')
