# Pruebas de P8

P8 reutiliza las pruebas existentes. No se copiaron dentro de esta carpeta.

- Unitarias e integración: `P7/tests/`.
- Verificación de categorías: `P7/scripts/verificar_distribucion_pruebas.py`.

El workflow de P8 ejecuta las pruebas y Helm antes de construir las imágenes.
Las pruebas de humo y carga no forman parte de esta versión porque el auxiliar
indicó que ya no se realizan.

El Rollout conserva los pesos 10%, 30% y 60%, con una validación de integración
al 30%. Si esa validación falla, Argo Rollouts detiene la promoción.
