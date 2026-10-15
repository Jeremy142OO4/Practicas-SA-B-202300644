#!/usr/bin/env bash
# Diagnostico independiente, sin puntajes. No reemplaza ni modifica al verificador.
set -euo pipefail
[[ $(kubectl config current-context) == gke_sa-202300644_us-central1-a_p9-202300644 ]] || exit 1
PORT=19091
kubectl port-forward -n observability svc/p10-monitoring-kube-promet-prometheus "$PORT:9090" >/dev/null 2>&1 &
PF=$!
trap 'kill "$PF" 2>/dev/null || true' EXIT
for attempt in $(seq 1 20); do
  kill -0 "$PF" 2>/dev/null || { echo 'No se pudo abrir el acceso a Prometheus'; exit 1; }
  if curl -fsS "http://127.0.0.1:$PORT/-/ready" >/dev/null 2>&1; then break; fi
  sleep 1
done
query() { curl -fsSG --max-time 10 "http://127.0.0.1:$PORT/api/v1/query" --data-urlencode "query=$1"; }

echo "Fecha UTC: $(date -u +%Y-%m-%dT%H:%M:%SZ)"
kubectl get applications -n argocd
echo 'Targets del gateway: una recoleccion por pod, todos saludables:'
curl -fsS "http://127.0.0.1:$PORT/api/v1/targets" | jq -e '
  [.data.activeTargets[] | select((.labels.pod // "") | startswith("api-gateway-"))]
  | sort_by(.labels.pod) | group_by(.labels.pod)
  | length > 0 and all(.[]; length == 1 and .[0].health == "up")'
echo 'Reglas de P10 evaluandose sin errores:'
curl -fsS "http://127.0.0.1:$PORT/api/v1/rules" | jq -e '
  [.data.groups[] | select(.name | startswith("p10.slo")) | .rules[]]
  | length == 12 and all(.[]; .health == "ok")'
kubectl get rollout api-gateway -n sa-p9 -o json | jq '{phase:.status.phase,replicas:.spec.replicas,ready:.status.readyReplicas,hash:.status.currentPodHash,image:.spec.template.spec.containers[0].image}'
echo 'Proveedores reales del AnalysisTemplate:'
AT=$(kubectl get analysistemplate api-gateway-metrics -n sa-p9 -o json)
jq -e 'all(.spec.metrics[]; .provider.prometheus != null and .provider.job == null)' <<< "$AT"
echo 'Promocion condicionada al analisis:'
kubectl get rollout api-gateway -n sa-p9 -o json | jq -e 'any(.spec.strategy.canary.steps[]?; .analysis != null)'
HASH=$(kubectl get rollout api-gateway -n sa-p9 -o jsonpath='{.status.currentPodHash}')
echo "Consultas completas con hash real: $HASH"
while IFS= read -r encoded; do
  metric=$(printf '%s' "$encoded" | base64 -d)
  name=$(jq -r .name <<< "$metric")
  promql=$(jq -r --arg h "$HASH" '.provider.prometheus.query | gsub("\\{\\{args\\.podHash\\}\\}"; $h)' <<< "$metric")
  echo "$name"
  query "$promql" | jq -e '.status=="success" and (.data.result | length > 0)'
  query "$promql" | jq '.data.result'
done < <(jq -r '.spec.metrics[] | @base64' <<< "$AT")

echo 'SLI y presupuesto con solicitudes recientes:'
for metric in slo:api_availability:ratio slo:api_latency:ratio slo:api_availability:error_budget:ratio slo:api_latency:error_budget:ratio payments_approved_total payments_rejected_total; do
  echo "$metric"
  query "$metric" | jq -e '.status=="success" and (.data.result | length > 0)'
  query "$metric" | jq '.data.result'
done
echo 'Reglas propias y runbooks:'
kubectl get prometheusrule p10-slos -n observability -o json | jq '{registros:[.spec.groups[].rules[]|select(.record)|.record],alertas:[.spec.groups[].rules[]|select(.alert)|{alert,expr,runbook:.annotations.runbook_url}]}'
echo 'Fuentes configuradas de Grafana (sin imprimir credenciales):'
kubectl get cm -n observability -o json | python3 -c '
import json,sys,yaml
types=set()
for cm in json.load(sys.stdin)["items"]:
 for raw in cm.get("data",{}).values():
  try: obj=yaml.safe_load(raw)
  except yaml.YAMLError: continue
  if isinstance(obj,dict):
   for ds in obj.get("datasources",[]):
    if isinstance(ds,dict): types.add(ds.get("type",""))
print("Fuentes:",", ".join(sorted(types)))
assert {"prometheus","loki","tempo"}.issubset(types)
'
echo 'Analisis recientes y evidencia de aborto:'
kubectl get analysisruns -n sa-p9 --sort-by=.metadata.creationTimestamp
echo 'Registro JSON y trace_id:'
POD=$(kubectl get pods -n sa-p9 -l app.kubernetes.io/name=api-gateway -o jsonpath='{.items[0].metadata.name}')
LINE=$(kubectl logs -n sa-p9 "$POD" --tail=30 | awk '/^\{/{line=$0} END{print line}')
jq -e '.trace_id != null and (.trace_id|length)==32' <<< "$LINE"
echo 'Fin del diagnostico independiente. Este resultado no modifica una nota.'
