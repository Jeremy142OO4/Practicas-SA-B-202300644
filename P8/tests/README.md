# Pruebas de P8

P8 reutiliza las pruebas existentes. No se copiaron dentro de esta carpeta.

- Unitarias e integración: `P7/tests/`.
- Verificación de categorías: `P7/scripts/verificar_distribucion_pruebas.py`.
- Carga: `P5/load-tests/load.js`.

El workflow de P8 ejecuta las pruebas y Helm antes de construir las imágenes.
La prueba de carga permanece disponible, pero no se ejecutó porque el auxiliar
indicó que ya no era necesaria para esta práctica.

En el Rollout se ejecutan tres validaciones: humo al 10%, integración al 30% y
humo final al 60%. Si una validación falla, Argo Rollouts detiene la promoción.
