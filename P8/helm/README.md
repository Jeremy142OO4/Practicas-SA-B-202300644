# Empaquetado Helm de P8

P8 reutiliza el chart compuesto de `P6/charts/sa-platform/` y sus subcharts,
sin copiar el código de las plantillas. Los archivos de esta carpeta son los
valores por ambiente que usa el pipeline para validar la parametrización.

- `values-dev.yaml`: una réplica inicial por servicio y namespace `sa-p8`.
- `values-prod.yaml`: réplicas de alta disponibilidad y política de pull
  `Always`.

Ambos ambientes usan etiquetas semánticas como `v1.0.0`; nunca se usa
`latest`. Los secretos son provistos por el gestor externo referenciado en el
repositorio GitOps.
