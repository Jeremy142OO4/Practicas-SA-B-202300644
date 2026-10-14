# Runbook de observabilidad

## 1. Ver la alerta

```bash
kubectl get prometheusrules -n observability
kubectl get prometheus -n observability
kubectl get alerts -A
```

## 2. Confirmar el sintoma

```bash
kubectl get rollout api-gateway -n sa-p9
kubectl get pods -n sa-p9 -o wide
kubectl describe rollout api-gateway -n sa-p9
```

En Grafana se revisan primero el trafico, los errores, la latencia p95 y la
saturacion. Se filtra el mismo intervalo en Loki usando el `trace_id`.

## 3. Llegar a la causa

1. La alerta indica el SLO afectado.
2. Prometheus muestra la consulta que cruzo el umbral.
3. Grafana abre el trace correspondiente en Tempo.
4. El `trace_id` del trace se busca en los logs JSON de Loki.
5. Se confirma el pod, endpoint y error concreto.

## 4. Recuperar

```bash
kubectl get analysisruns -n sa-p9
kubectl get rollout api-gateway -n sa-p9
kubectl argo rollouts get rollout api-gateway -n sa-p9
```

No se hace un despliegue directo. La recuperacion se realiza revirtiendo la
version en GitOps y dejando que ArgoCD sincronice el cambio.
