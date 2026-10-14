# Informe de operacion

> Este informe se completa con los valores observados durante la demostracion.

## 1. SLOs

- Disponibilidad del API: 99% en 30 dias.
- Latencia p95 del API: 95% de solicitudes por debajo de 500 ms en 30 dias.

## 2. Falla inducida

Se publica una version candidata que aumenta la latencia de una ruta del
`api-gateway`. La falla es diferente a la usada en P8: el criterio de
reversion es una metrica real de Prometheus.

## 3. Tiempo de deteccion

- Inicio: **PENDIENTE**
- Alerta/AnalysisRun: **PENDIENTE**
- Tiempo total: **PENDIENTE**

## 4. Ruta de sintoma a causa

Alerta de error budget -> consulta PromQL -> trace en Tempo -> `trace_id` en
Loki -> pod, endpoint y error responsable.

## 5. Tiempo de recuperacion

- Rollout abortado: **PENDIENTE**
- Version estable disponible: **PENDIENTE**
- Tiempo total: **PENDIENTE**

## 6. Alerta faltante o excesiva

La alerta principal es la de consumo rapido del error budget. Se revisara si
el umbral produce ruido durante la prueba y se documentara cualquier ajuste.
