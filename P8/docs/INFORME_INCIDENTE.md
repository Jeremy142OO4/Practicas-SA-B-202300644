# Informe de incidente - fallo controlado

## Qué falló

Se probó una versión que no pasó el análisis configurado para el canary. La
métrica quedó en estado `Failed` y superó el límite permitido.

Esta prueba se realizó antes de reconstruir el clúster anterior con Terraform.
La captura se conserva en
[`evidencias/07-fallo-rollout-historico.png`](../evidencias/07-fallo-rollout-historico.png).

## Cómo se detectó

Argo Rollouts informó:

```text
RolloutAborted: Rollout aborted update to revision 5
Metric assessed Failed
failed (1) > failureLimit (0)
```

El análisis fue el que detuvo la promoción.

## Cómo se contuvo

Argo Rollouts abortó la promoción antes de continuar con el siguiente paso. La
versión estable quedó atendiendo el servicio. No se hizo un cambio manual en
Kubernetes.

## Tiempo de recuperación

La captura disponible no conserva los minutos exactos del incidente. Sí confirma
que la promoción se detuvo automáticamente. Después se reconstruyó el clúster
con Terraform y se verificó nuevamente el estado saludable.

## Cómo prevenirlo

- Mantener la prueba de integración antes de promover.
- Usar imágenes versionadas y firmadas.
- Revisar las imágenes con Trivy.
- Mantener Kyverno en modo `Enforce`.
- Revisar ArgoCD después de cada cambio.
