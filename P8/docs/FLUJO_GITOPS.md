# Flujo GitOps de P8

## Qué hace cada parte

### Repositorio principal

Aquí están el código reutilizado, el workflow, Terraform, las políticas y la
documentación de P8. GitHub Actions valida el cambio, construye las imágenes,
genera el SBOM y las firma.

### Repositorio GitOps

Está separado del código principal:

<https://github.com/Jeremy142OO4/Practicas-SA-B-202300644-gitops>

Aquí solo se guardan los manifiestos que deben llegar al clúster. El cambio de
imagen se solicita mediante Pull Request.

### ArgoCD

ArgoCD revisa la rama `main` del repositorio GitOps. Cuando encuentra un cambio,
lo sincroniza en `sa-p8`. El pipeline no aplica manifiestos directamente.

### Argo Rollouts

Argo Rollouts mueve el tráfico en 10%, 30% y 60%. En cada etapa ejecuta una
validación. Si falla, detiene la promoción y conserva la versión estable.

## Reglas importantes

- No se guarda `kubeconfig` en GitHub.
- El workflow no ejecuta `kubectl apply` contra GKE.
- No se usa la etiqueta `latest`.
- Las imágenes se identifican con una versión.
- Los secretos no se escriben en el repositorio GitOps.
- Kyverno bloquea imágenes sin versión, contenedores con root y firmas no válidas.
