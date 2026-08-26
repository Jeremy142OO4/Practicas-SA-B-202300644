# Evidencias: persistencia, upgrade/rollback y cero downtime

## 1. Persistencia de PostgreSQL

Se elimino unicamente el Pod `sa-platform-orders-database-0`, conservando su PVC.

- UID anterior del Pod: `bda7b73b-2ac7-4767-8765-44022eca4944`.
- UID nuevo del Pod: `c16acff4-106e-4717-9862-9d170135f7cf`.
- Datos antes: 9 ejecuciones y 4 resumenes.
- Datos despues: 9 ejecuciones y 4 resumenes.
- El registro testigo con `id_ejecucion=9` continuo existiendo.
- PVC despues de recrear el Pod: `Bound`, volumen `pvc-11a89708-3962-458b-9c36-acf235f813d2`, capacidad 8 GiB.

Esto demuestra que los datos residen en el volumen persistente y no en el sistema de archivos temporal del Pod.

## 2. Upgrade y rollback con Helm

El chart dispone de dos archivos pequenos de valores para cambiar solamente la etiqueta de la imagen del API Gateway:

- `values-version-1.yaml`: imagen `1.0.0-dev`.
- `values-version-2.yaml`: imagen `1.1.0-dev`.

El chart padre version 0.2.0 representa la segunda version. La prueba real produjo este historial relevante:

- Revision 9: upgrade a `sa-platform-0.2.0`, aplicacion 1.1.0.
- Revision 10: rollback a la revision 8, recuperando `sa-platform-0.1.0`, aplicacion 1.0.0.
- Despues del rollback se verifico la imagen `practicas-sa/api-gateway:1.0.0-dev`.
- Revision 12: estado final desplegado con `sa-platform-0.2.0`, aplicacion 1.1.0 e imagen `practicas-sa/api-gateway:1.1.0-dev`.

Comando de upgrade:

```bash
helm upgrade sa-platform P5/charts/sa-platform \
  -n default --reuse-values \
  -f P5/charts/sa-platform/values-version-2.yaml \
  --atomic --wait --timeout 5m
```

Comando de rollback usado en la prueba:

```bash
helm rollback sa-platform 8 -n default --wait --timeout 5m
```

## 3. Actualizacion sin interrupcion

La medicion se realizo desde un cliente temporal dentro del cluster, consumiendo el Service ClusterIP del API Gateway mientras Helm realizaba el upgrade de la imagen 1.0.0 a la 1.1.0.

- Solicitudes enviadas: 800.
- Respuestas HTTP 200: 800.
- Errores HTTP o de conexion: 0.
- Minimo de replicas listas observado: 2.
- Estado final: 2 de 2 replicas listas.
- Imagen final de ambos Pods: `practicas-sa/api-gateway:1.1.0-dev`.

Para evitar que una conexion sea enviada a un Pod durante su terminacion, el API Gateway incluye un hook `preStop` de 5 segundos y `terminationGracePeriodSeconds: 30`. Esto complementa `RollingUpdate` con `maxUnavailable: 0`, las readiness probes y el PodDisruptionBudget ya configurados.

`kubectl port-forward` no se utilizo como medidor de disponibilidad porque queda enlazado a un Pod concreto y su conexion se corta cuando ese Pod es reemplazado; el trafico real entre componentes utiliza el Service ClusterIP.
