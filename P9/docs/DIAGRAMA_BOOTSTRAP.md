# Diagrama de bootstrap y recuperacion de la Practica 9

El siguiente diagrama muestra el flujo de reconstruccion de la plataforma y
las pruebas de continuidad. Se presenta como un flujo de pasos y decisiones,
por lo que no depende de colores para entenderse.

![Diagrama de bootstrap de P9](./Diagrama%20Bootstrap.png)

## Como leer el diagrama

1. El flujo inicia confirmando que existe un backup valido y registrando la
   hora de inicio.
2. En la seccion **Bootstrap automatico** se ejecuta `bootstrap.sh`. Terraform
   prepara el estado remoto, crea GKE, configura Velero, Secret Manager e
   instala ArgoCD.
3. ArgoCD aplica `p9-root` y sincroniza las aplicaciones hijas: Velero,
   External Secrets, Argo Rollouts, las politicas y la plataforma.
4. En **Plataforma recuperada automaticamente** se muestra la dependencia
   entre PostgreSQL con su PVC, RabbitMQ, `orders-ms` y el API Gateway.
5. En **Acciones manuales** se decide si se necesita una reconstruccion desde
   cero o si el sistema ya esta funcionando. En una reconstruccion se aisla o
   destruye el entorno, se inicia el cronometro y se vuelve a ejecutar el
   bootstrap.
6. En **Pruebas de continuidad** se puede seleccionar la restauracion de datos
   o la prueba de perdida de nodo.
7. La restauracion comprueba que el PVC quede `Bound` y que el dato real
   `dr_marker` pueda consultarse en PostgreSQL.
8. La prueba de nodo comprueba el drenaje, la reprogramacion de replicas y el
   funcionamiento del servicio mediante el PDB.
9. Al final se registran las evidencias, el RTO y el RPO.

El repositorio de codigo reutiliza la infraestructura de `P8/terraform` como
modulo. El repositorio GitOps reutiliza los charts y manifiestos de P8 mediante
referencias remotas; P9 agrega el flujo de continuidad operativa, respaldos y
restauracion.
