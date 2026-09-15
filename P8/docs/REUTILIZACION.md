# Reutilización entre prácticas

P8 no copia las aplicaciones ni los charts completos de las prácticas anteriores. Cada componente se consume desde su ubicación original dentro del repositorio.

| Componente | Ubicación original | Uso en P8 |
|---|---|---|
| Microservicios y Dockerfiles productivos | `P4/` | Build y publicación de las imágenes de la candidata |
| Chart Helm base | `P6/charts/sa-platform/` | Referencia para la parametrización y los recursos existentes |
| Pruebas unitarias e integración | `P7/tests/` | Puerta de calidad del pipeline |
| Verificación de categorías | `P7/scripts/verificar_distribucion_pruebas.py` | Comprobación de la suite heredada |
| Evidencias anteriores | `P5/Documentacion/` y `P7/evidencias/` | Contexto histórico, no sustituyen las evidencias de P8 |

## Regla de mantenimiento

Si un cambio pertenece al código de negocio, se realiza en `P4/`. Si pertenece a la entrega progresiva, seguridad de suministro, GitOps o infraestructura declarativa de P8, se agrega en `P8/` o en el repositorio GitOps independiente.
