# Flujo y componentes

## Una solicitud

El gateway recibe la solicitud. Para una operacion de carrito, orders-ms
consulta catalog-ms. El encabezado `traceparent` mantiene la misma traza
entre los tres servicios. No todas las solicitudes recorren los tres:
una consulta directa al catalogo pasa solamente por gateway y catalogo.

Cada servicio expone metricas y escribe logs JSON con `trace_id`.
Prometheus recoge las metricas con ServiceMonitor. OpenTelemetry Collector
recibe las trazas y las envia a Tempo. Alloy consulta los logs de Kubernetes
y los envia a Loki. Grafana consulta los tres almacenes.

## Cuatro señales de oro

| Señal | Pregunta | Metrica |
|---|---|---|
| Latencia | ¿Las solicitudes terminan en menos de 500 ms? | `http_server_request_duration_seconds` |
| Trafico | ¿Cuantas solicitudes llegan? | `http_server_requests_total` |
| Errores | ¿Cuantas respuestas HTTP son 5xx? | El mismo contador filtrado por `status` |
| Saturacion | ¿Que pods consumen CPU y memoria? | Metricas de contenedores de Kubernetes |

Los paneles explican la pregunta que responden y muestran unidades.
Las rutas /health y /metrics se excluyen de los indicadores de negocio.
Un error GraphQL puede devolver HTTP 200: no se interpreta como un 5xx.
Las metricas `payments_approved_total` y `payments_rejected_total` cuentan
resultados de pago procesados; una excepcion de la API no es un pago rechazado.
No se inventan pagos para llenar el tablero: esos contadores pueden valer cero.

## SLO y presupuesto

La ventana del laboratorio es **una hora**. Disponibilidad: al menos 99% de
respuestas sin 5xx. Latencia: al menos 95% de solicitudes dentro de 500 ms.
La latencia p95 se presenta en otro panel; el indicador del SLO es la
proporcion de solicitudes que cumple el umbral, no el percentil en si.

Se elige una hora para evaluar una sesion de demostracion con los datos que
realmente se conservan. No se afirma cumplimiento de 30 dias. El 99% permite
hasta 1% de respuestas fallidas; 500 ms mantiene una espera corta para las
operaciones interactivas del carrito, permitiendo 5% de excepciones.
Son objetivos del laboratorio, no una promesa validada durante un mes.

Presupuesto consumido = proporcion de solicitudes malas / presupuesto permitido.
Por ejemplo, 2% de errores consume 200% del presupuesto de disponibilidad (1%).
Las reglas usan `increase(...[1h])` y calculan la ventana disponible desde
que Prometheus comenzo a recolectar. Antes de completar una hora, el resultado
es provisional. La retencion es dos dias y el almacenamiento actual es
temporal: apagar o reemplazar el pod puede perder ese historial.

## Alertas y entrega progresiva

Las alertas comparan la velocidad de consumo con el presupuesto: mas de 2
veces el ritmo permitido en ventanas de 1 y 5 minutos, durante 30 segundos,
y con al menos dos solicitudes observadas. Si ese ritmo se mantiene, el
presupuesto de una hora se agota en menos de media hora. Por eso requieren
revisar la candidata. Alertmanager recibe las alertas; el runbook indica la accion.
La latencia se evalua tambien por pod de la candidata con el mismo presupuesto
de 5%. Asi el trafico sano de stable no oculta una candidata lenta.

Argo Rollouts ejecuta el AnalysisTemplate y consulta Prometheus. No es
Prometheus quien ejecuta el despliegue. El analisis mide solamente el pod hash
de la candidata, usa solicitudes de negocio de la ultima ventana de un minuto
y exige datos validos. No convierte la ausencia de datos en un resultado sano.

La candidata debe mantenerse por debajo de 5% de 5xx y 500 ms de p95.
La primera medicion se retrasa 45 segundos para recoger muestras.
Una medicion fallida aborta la promocion (`failureLimit: 0`).
Los servicios stable/canary cambian sus selectores para conservar la version
estable. Se demuestra la falla con solicitudes HTTP de baja frecuencia,
sin un Job que decida el resultado y sin prueba de carga.

GitHub Actions construye y firma imagenes. Un PR actualiza GitOps.
ArgoCD aplica los cambios y Kyverno conserva las verificaciones de admision.
La misma aplicacion `p9-platform` administra `p10/platform`.

