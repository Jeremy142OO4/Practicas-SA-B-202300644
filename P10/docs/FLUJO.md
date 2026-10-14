# Flujo y componentes

## Flujo de una solicitud

1. `api-gateway` recibe la solicitud y crea o continua un trace.
2. El encabezado W3C `traceparent` viaja con la llamada al microservicio.
3. `orders-ms` registra el resultado de la operacion y llama a `catalog-ms`.
4. Cada servicio publica metricas HTTP en `/metrics` y escribe logs JSON.
5. OpenTelemetry Collector recibe las trazas y las envia a Tempo.
6. Prometheus consulta las metricas, calcula los SLO y ejecuta el analisis del canary.
7. Grafana une Prometheus, Loki y Tempo para investigar el mismo `trace_id`.

## Cuatro señales doradas

- **Latencia:** histograma `http_server_request_duration_seconds`.
- **Trafico:** contador `http_server_requests_total`.
- **Errores:** solicitudes con estado HTTP 5xx.
- **Saturacion:** uso de CPU y memoria del pod y de los nodos.

## Relacion con Argo Rollouts

El Rollout mantiene dos servicios, estable y canary. En cada paso de peso se
ejecuta un `AnalysisTemplate` que consulta directamente a Prometheus. Si el
porcentaje de errores o el p95 del canary supera el umbral, el AnalysisRun
falla y Argo Rollouts aborta la promocion.
