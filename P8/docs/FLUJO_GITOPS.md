# Flujo GitOps de P8

## Responsabilidades

### Repositorio de código

Contiene el código de P4 y la configuración de P8 necesaria para validar y publicar una imagen. GitHub Actions ejecuta build, pruebas, análisis de seguridad, generación de SBOM y firma.

### Repositorio GitOps

Contiene únicamente los manifiestos declarativos de Kubernetes, la aplicación de ArgoCD, los Rollouts, las políticas y la configuración segura de secretos. El cambio de imagen se realiza mediante Pull Request.

### Clúster

ArgoCD observa el repositorio GitOps y es el único componente autorizado para aplicar cambios de aplicación. Argo Rollouts controla la promoción progresiva y ejecuta los análisis antes de avanzar.

## Restricciones

- El pipeline de P8 no debe contener kubeconfig.
- El pipeline de P8 no debe ejecutar `kubectl apply`, `kubectl set image` ni `helm upgrade` contra el clúster.
- La imagen debe usar una versión semántica derivada de Git.
- No se debe utilizar la etiqueta `latest`.
- Los secretos no deben almacenarse en texto plano en el repositorio GitOps.

