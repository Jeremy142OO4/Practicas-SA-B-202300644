# Practica 8 - GitOps, entrega progresiva y seguridad de la cadena de suministro

P8 evoluciona el flujo de P7 hacia GitOps. El repositorio de codigo valida y
publica artefactos; el repositorio GitOps es la fuente de verdad; ArgoCD es el
unico componente que aplica cambios al cluster; y Argo Rollouts controla la
promocion canary.

## Reutilizacion sin duplicar codigo

- Microservicios y Dockerfiles: `P4/`.
- Chart Helm y subcharts: `P6/charts/sa-platform/`.
- Prueba de carga: `P5/load-tests/load.js`.
- Pruebas unitarias e integracion: `P7/tests/`.
- Verificador de la suite: `P7/scripts/verificar_distribucion_pruebas.py`.

P8 agrega solamente configuracion, infraestructura, politicas, manifiestos,
documentacion y evidencias nuevas.

## Componentes implementados

- Terraform para namespace, ResourceQuota, LimitRange, ServiceAccounts y RBAC.
- Valores Helm separados para `dev` y `prod`, lint sobre el chart de P6.
- Workflow manual en `.github/workflows/p8-ci.yml` con pruebas, Helm lint,
  Terraform, Trivy, SBOM, build, firma Cosign y PR al repo GitOps.
- Rollout canary con pesos de 10%, 30% y 60% y tres AnalysisTemplates.
- Politicas Kyverno para bloquear `latest`, exigir recursos, impedir root y
  verificar firmas Cosign.
- ExternalSecret sin credenciales en texto plano.

## Tabla de enlaces de entrega

| Item | Enlace o dato |
|---|---|
| Repositorio GitOps | https://github.com/Jeremy142OO4/Practicas-SA-B-202300644-gitops |
| Aplicacion en ArgoCD | `sa-platform` en namespace `argocd`; destino `sa-p8` |
| Ejecucion exitosa del pipeline | Pendiente de ejecucion manual del workflow |
| Reversion automatica | Pendiente de evidencia del cluster en `evidencias/` |
| Despliegue rechazado por politica | `../Practicas-SA-B-202300644-gitops/policies/` y evidencia pendiente |
| Bloqueo por vulnerabilidad critica | Job `supply-chain` de `p8-ci.yml`; run pendiente |
| Imagen firmada | `ghcr.io/jeremy142oo4/sa-api-gateway:v1.0.0`, al publicar la version |
| Reporte de prueba de carga | `P5/load-tests/load.js`; reporte de ejecucion pendiente |
| Video demostrativo | Pendiente de grabacion; incluir minutaje por criterio |

## Restricciones de seguridad

El workflow de P8 no contiene kubeconfig, `kubectl apply`, `kubectl set image`
ni `helm upgrade`. La promocion de una imagen se solicita mediante Pull
Request al repositorio GitOps y la aplicacion al cluster queda a cargo de
ArgoCD.

El workflow de P7 se conserva, pero ambos workflows permanecen manuales
durante el desarrollo para evitar ejecuciones automaticas.

## Evidencias faltantes para el cierre

1. Ejecutar el workflow manual con una version semantica.
2. Configurar el cluster y capturar ArgoCD Synced/Healthy.
3. Capturar una promocion canary exitosa y una reversion automatica.
4. Capturar un rechazo de Kyverno y un bloqueo de Trivy.
5. Completar `docs/INFORME_INCIDENTE.md`, el diagrama y el video.
