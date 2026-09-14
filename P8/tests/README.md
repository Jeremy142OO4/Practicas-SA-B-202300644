# Pruebas de P8

P8 reutiliza sin copiar las pruebas de las practicas anteriores:

- Unitarias e integracion: `P7/tests/`.
- Verificacion de categorias: `P7/scripts/verificar_distribucion_pruebas.py`.
- Carga: `P5/load-tests/load.js`.

Los umbrales de promocion canary se declaran en el repositorio GitOps:

- Humo: tres ejecuciones, sin respuestas distintas de HTTP 2xx.
- Integracion: dos ejecuciones, salud y disponibilidad deben responder.
- Carga: error menor a 5% y percentil 95 menor a 500 ms.

La ejecucion contra un cluster real y sus reportes se agregaran en `P8/evidencias/`
despues de configurar Argo Rollouts.
