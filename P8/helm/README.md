# Helm de P8

P8 reutiliza el chart de `P6/charts/sa-platform/`. No se copiaron las plantillas
de los microservicios.

Los archivos de esta carpeta separan los valores por ambiente:

- `values-dev.yaml`: valores para desarrollo.
- `values-prod.yaml`: valores para producción.

Las imágenes usan versiones como `v1.0.0`; no se usa `latest`.

Para validar el chart:

```bash
cd P6/charts/sa-platform
helm dependency build
helm lint . -f ../../../P8/helm/values-dev.yaml
helm lint . -f ../../../P8/helm/values-prod.yaml
```

El pipeline ejecuta esta validación antes de construir y publicar las imágenes.
