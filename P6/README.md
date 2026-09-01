# UNIVERSIDAD DE SAN CARLOS DE GUATEMALA

## FACULTAD DE INGENIERIA

### SOFTWARE AVANZADO

**Nombre:** Jeremy Estuardo Orellana Aldana  
**Carne:** 202300644  
**Seccion:** B

# Practica 6 - Despliegue de la plataforma en Amazon EKS

## 1. Alcance

Esta practica despliega en Amazon Elastic Kubernetes Service la plataforma construida en la Practica 5. No modifica la funcionalidad de la tienda: adapta el registro de imagenes, el almacenamiento, las credenciales y la exposicion publica a un entorno de nube real.

El despliegue contiene API Gateway, Auth, Catalog, Orders, Payments, cuatro bases PostgreSQL, RabbitMQ, dos CronJobs y el consumidor de resumenes. El flujo de pagos permanece asincrono mediante la cola durable `payment.requests`.

## 2. Requisitos previos

- Una cuenta de AWS activa con permisos para EKS, EC2, IAM y CloudFormation.
- AWS CloudShell con AWS CLI, Docker, `kubectl`, Helm y `eksctl`.
- Una cuenta de Docker Hub con el repositorio público `jeremy142004/sa-p6` y un token nuevo con permiso de escritura.
- Permisos para crear clústeres, nodos administrados, roles IAM, VPC, volúmenes EBS y balanceadores.

Antes de crear recursos se debe configurar un presupuesto en **AWS Billing and Cost Management > Budgets**. Una alerta informa sobre consumo, pero no detiene automaticamente los recursos. EKS, EC2, EBS y el balanceador pueden generar cargos desde el momento de su creación.

Verificar las herramientas en AWS CloudShell:

```bash
aws --version
kubectl version --client
helm version
eksctl version
aws sts get-caller-identity
```

Si `eksctl` no está disponible, instalarlo desde su distribución oficial:

```bash
ARCH=amd64
PLATFORM="$(uname -s)_$ARCH"
curl -sLO "https://github.com/eksctl-io/eksctl/releases/latest/download/eksctl_$PLATFORM.tar.gz"
curl -sL "https://github.com/eksctl-io/eksctl/releases/latest/download/eksctl_checksums.txt" \
  | grep "$PLATFORM" | sha256sum --check
tar -xzf "eksctl_$PLATFORM.tar.gz" -C /tmp
rm "eksctl_$PLATFORM.tar.gz"
sudo install -m 0755 /tmp/eksctl /usr/local/bin/eksctl
rm /tmp/eksctl
```

## 3. Estructura de P6

```text
P6/
|-- charts/sa-platform/          Chart adaptado a Amazon EKS
|-- scripts/p6.sh               Automatizacion para AWS CloudShell
|-- evidencias/                 Capturas requeridas para la entrega
|-- .gitignore                  Proteccion de credenciales locales
`-- README.md                   Procedimiento y respuestas teoricas
```

No se incluyen diagramas porque el enunciado no los solicita.

## 4. Configuracion elegida

| Elemento | Configuracion |
| --- | --- |
| Proveedor | Amazon Web Services |
| Kubernetes | Amazon EKS con grupo de nodos administrado |
| Region predeterminada | `us-east-1` |
| Nodos | 2 nodos `t3.medium` |
| Registro | Docker Hub público: `jeremy142004/sa-p6` |
| Entrada publica | Service `LoadBalancer` |
| Persistencia | `sa-p6-gp3`, volúmenes Amazon EBS cifrados |
| Namespace | `sa-p6` |

El tamaño de maquina se puede cambiar al ejecutar el script si la cuota o la carga real lo requieren.

## 5. Creacion del clúster y del registro

Clonar el repositorio y ejecutar el procedimiento desde AWS CloudShell:

```bash
cd ~/Practicas-SA-B-202300644/P6
chmod +x scripts/p6.sh
./scripts/p6.sh crear-infraestructura
```

El script crea un clúster EKS con exactamente dos nodos administrados, configura el contexto remoto de `kubectl`, asocia el proveedor OIDC, crea el rol mínimo requerido por EBS CSI e instala el complemento oficial de almacenamiento.

Verificacion:

```bash
kubectl config current-context
kubectl get nodes -o wide
kubectl get storageclass
```

Se debe capturar la consola de EKS mostrando el clúster y el grupo con dos nodos, junto con la salida de estos comandos.

## 6. Publicacion de imagenes

```bash
./scripts/p6.sh publicar-imagenes
```

El script solicita el token de Docker Hub de forma oculta, lo transmite por entrada estándar y cierra la sesión al terminar. El token nunca se escribe en el repositorio. Antes de ejecutar este paso se debe crear en Docker Hub un único repositorio público llamado `sa-p6` dentro de la cuenta `jeremy142004`.

Se construyen y publican estas seis imagenes:

- `api-gateway`
- `auth-ms`
- `catalog-ms`
- `orders-ms`
- `payments-ms`
- `scheduled-jobs`

Las seis imágenes se publican en un solo repositorio usando etiquetas diferentes:

```text
jeremy142004/sa-p6:api-gateway-1.0.0
jeremy142004/sa-p6:auth-ms-1.0.0
jeremy142004/sa-p6:catalog-ms-1.0.0
jeremy142004/sa-p6:orders-ms-1.0.0
jeremy142004/sa-p6:payments-ms-1.0.0
jeremy142004/sa-p6:scheduled-jobs-1.0.0
```

Verificacion:

La verificación se realiza en la pestaña **Tags** de `https://hub.docker.com/r/jeremy142004/sa-p6/tags`.

## 7. Credenciales y Secrets

Las credenciales reales no se almacenan en Git. Este comando genera valores criptograficamente aleatorios en `P6/valores-secretos.yaml`:

```bash
./scripts/p6.sh preparar-secretos
git check-ignore valores-secretos.yaml
```

El segundo comando debe confirmar que el archivo esta ignorado. Helm transforma esos valores en objetos `Secret`; los Pods los reciben mediante `envFrom`. Si el archivo ya existe, el script se detiene y no lo reemplaza.

## 8. Despliegue con Helm

```bash
./scripts/p6.sh desplegar
```

El comando descarga las dependencias del chart, incorpora los scripts SQL de inicializacion, configura las seis etiquetas públicas de Docker Hub y realiza un `helm upgrade --install` atomico. Si el despliegue no queda listo, Helm revierte la operacion.

Verificacion completa:

```bash
./scripts/p6.sh verificar
```

La verificacion exige que todos los Pods esten `Ready`, espera la direccion del balanceador y realiza una peticion mediante la IP publica:

```text
http://IP_PUBLICA:4000/health
```

Esa solicitud se origina fuera del clúster y sirve como evidencia de acceso desde internet.

## 9. Persistencia y adaptacion de red

Las cuatro bases PostgreSQL y RabbitMQ solicitan volúmenes de 8 GiB con la `StorageClass` `sa-p6-gp3`. El complemento Amazon EBS CSI aprovisiona volúmenes `gp3` cifrados. Los PVC deben aparecer como `Bound`:

```bash
kubectl -n sa-p6 get pvc
kubectl get pv
```

El API Gateway es el unico Service de tipo `LoadBalancer`. Los demas servicios permanecen como `ClusterIP`. La politica predeterminada niega trafico y luego permite de forma explicita el acceso publico al Gateway y los flujos internos necesarios.

## 10. Comunicacion asincrona

Orders publica las solicitudes de pago en `payment.requests`. Payments consume el mensaje, registra el resultado y notifica a Orders. La solicitud inicial devuelve `PAGO_EN_PROCESO` y posteriormente la orden cambia a `PAGADA` o `PAGO_RECHAZADO`.

Para la evidencia se debe realizar el flujo funcional desde la IP publica, mostrar los logs de Orders y Payments y comprobar que ambos interactuaron con RabbitMQ:

```bash
kubectl -n sa-p6 logs deployment/orders-ms --tail=100
kubectl -n sa-p6 logs deployment/payments-ms --tail=100
kubectl -n sa-p6 get pods
```

## 11. Evidencias requeridas

Guardar las capturas en `P6/evidencias` con estos nombres sugeridos:

1. `01-cluster-eks.png`: clúster administrado en la consola de AWS.
2. `02-dos-nodos.png`: los dos nodos y el contexto remoto.
3. `03-docker-hub.png`: el repositorio público con las seis etiquetas.
4. `04-pods-running.png`: Pods y StatefulSets disponibles.
5. `05-pvc-bound.png`: PVC enlazados y `sa-p6-gp3`.
6. `06-ip-publica.png`: Service `LoadBalancer` con IP externa.
7. `07-peticion-internet.png`: respuesta exitosa mediante la direccion publica.
8. `08-comunicacion-asincrona.png`: solicitud y procesamiento asincrono del pago.
9. `09-costos.png`: estimacion o reporte de facturacion.
10. `10-recursos-eliminados.png`: consola después de eliminar el clúster.

Las capturas se agregaran cuando exista el despliegue real en AWS.

## 12. Costos aproximados y reduccion

El despliegue genera cargos por el plano de control de EKS, dos instancias EC2, discos de arranque, cinco volúmenes EBS `gp3`, IPv4 pública y el balanceador externo. EKS con soporte estándar cobra USD 0.10 por hora por clúster, además de los demás recursos. El monto exacto depende de la región, las horas activas, el tráfico y los créditos de la cuenta; se debe adjuntar el reporte real de AWS Cost Explorer.

Para reducir costos:

- Mantener el clúster solo durante la demostracion.
- Usar nodos pequeños compatibles con la carga.
- No reservar una IP estatica innecesaria.
- Conservar replicas y discos en el minimo requerido.
- Configurar AWS Budgets antes del despliegue.
- Eliminar clúster, discos y balanceador al finalizar.

## 13. Eliminacion de recursos

Después de guardar todas las evidencias:

```bash
./scripts/p6.sh eliminar
```

Docker Hub se conserva porque es independiente de AWS. Después se debe revisar manualmente que no queden recursos facturables:

```bash
aws eks list-clusters --region us-east-1
aws ec2 describe-instances --region us-east-1 --filters Name=instance-state-name,Values=pending,running
aws ec2 describe-volumes --region us-east-1 --filters Name=status,Values=available,in-use
aws elbv2 describe-load-balancers --region us-east-1
aws elb describe-load-balancers --region us-east-1
```

## 14. Respuestas teoricas

### 14.1 Que es un clúster administrado y que diferencias tiene frente a uno local

Un clúster administrado es un entorno de Kubernetes cuyo plano de control es operado por un proveedor. Amazon EKS ofrece la API de Kubernetes, administra la disponibilidad del plano de control y lo integra con VPC, IAM, EBS y Elastic Load Balancing. En un clúster local, la computadora del estudiante simula o aloja estos componentes y normalmente no ofrece balanceadores, discos ni direcciones publicas reales.

La administración no elimina la responsabilidad del estudiante: todavía debe definir los recursos de la aplicación, imagenes, permisos, credenciales, politicas de red, escalado, observabilidad y costos.

### 14.2 Que es un Service LoadBalancer y como lo implementa AWS

Es un Service que solicita al proveedor una entrada externa para un conjunto de Pods. Al crearlo en EKS, la integración de AWS aprovisiona un Elastic Load Balancer público. El balanceador entrega el tráfico al Service y Kubernetes lo distribuye a los Pods disponibles del API Gateway.

### 14.3 Que es un registro de contenedores y por que es necesario

Es un servicio que almacena imagenes y etiquetas de contenedores en una ubicación accesible. Docker Hub permite publicar las imágenes y, al utilizar un repositorio público, los nodos de EKS pueden descargarlas sin credenciales adicionales. Las imagenes locales no son accesibles para nodos remotos; por eso deben publicarse antes del despliegue.

### 14.4 Que administra el proveedor y que administra el estudiante

AWS administra la infraestructura del plano de control y la API de Kubernetes. El estudiante administra el grupo de nodos, versión, cantidad y tamaño de instancias, cargas, recursos, imágenes, Secrets, NetworkPolicies, almacenamiento, exposición, monitoreo y eliminación.

### 14.5 Que costos genera y como pueden reducirse

Los costos principales provienen de nodos, discos, balanceador, administración del clúster, registro y transferencia. Se reducen limitando tiempo de ejecución y tamaño, usando creditos educativos, evitando recursos ociosos, configurando presupuestos y eliminando todos los recursos al concluir.

## 15. Material oficial consultado

- [Creación y administración de clústeres con eksctl](https://docs.aws.amazon.com/eks/latest/eksctl/creating-and-managing-clusters.html)
- [Almacenamiento EBS CSI en Amazon EKS](https://docs.aws.amazon.com/eks/latest/userguide/ebs-csi.html)
- [Balanceo de carga en Amazon EKS](https://docs.aws.amazon.com/eks/latest/best-practices/load-balancing.html)
- [Precios de Amazon EKS](https://aws.amazon.com/eks/pricing/)
- [Publicacion de imagenes en Docker Hub](https://docs.docker.com/docker-hub/repos/manage/hub-images/push/)
- [Uso seguro de tokens de Docker Hub](https://docs.docker.com/security/access-tokens/)
