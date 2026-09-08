# Práctica 7 - Integración y despliegue continuo

## Objetivo

Automatizar la compilación, las pruebas, la publicación de imágenes y el despliegue en Kubernetes de los microservicios desarrollados en las prácticas anteriores.

## Pipeline

El workflow se encuentra en `.github/workflows/p7-ci-cd.yml` y contiene cuatro fases:

1. **Build:** compila los tres servicios Go y valida los dos servicios Python.
2. **Test:** ejecuta y reporta por separado pruebas unitarias, de integración y E2E.
3. **Docker:** construye las cinco imágenes y las publica en Docker Hub después de superar build y test.
4. **Deploy:** actualiza GKE con Helm únicamente al publicar un release o solicitarlo manualmente.

## Estrategia de pruebas 70/20/10

La indicación del auxiliar se implementa sobre una base inicial de diez casos:

| Categoría | Casos | Proporción | Alcance |
|---|---:|---:|---|
| Unitarias | 7 | 70% | Reglas de carrito, órdenes y pagos con dependencias aisladas |
| Integración | 2 | 20% | Interacción del servicio de órdenes con catálogo, repositorio y mensajería |
| E2E | 1 | 10% | Contrato completo de build de imágenes y despliegue |

`scripts/verificar_distribucion_pruebas.py` falla si una prueba carece de categoría o si la distribución deja de ser 7/2/1. Cuando aumente la suite se debe conservar la proporción 70/20/10.

## Ejecución local

```bash
cd P7
python -m pip install -r requirements-ci.txt
python scripts/verificar_distribucion_pruebas.py
pytest -m unit -v
pytest -m integration -v
pytest -m e2e -v
```

Las pruebas unitarias y de integración ligera no requieren Docker. La validación real de las imágenes y del flujo E2E se realizará en la laptop Fedora con Docker.

## Configuración de GitHub

Para publicar en Docker Hub se deben crear los secretos `DOCKERHUB_USERNAME` y `DOCKERHUB_TOKEN`. El token solo se configura como GitHub Secret y nunca se almacena en el repositorio.

Para desplegar en GKE se deben crear:

- Secret `GCP_SA_KEY`: credencial JSON de una cuenta de servicio con permisos mínimos sobre GKE.
- Secret `P7_VALUES_YAML`: contenido del archivo privado de valores de Helm con las conexiones de bases de datos, JWT, cifrado y contraseña de RabbitMQ.
- Variable `GCP_PROJECT_ID`: identificador del proyecto de GCP.
- Variable `GKE_CLUSTER`: nombre del clúster.
- Variable `GKE_ZONE`: zona del clúster, por ejemplo `us-central1-a`.

El entorno de GitHub `production` puede configurarse con aprobación manual para proteger el despliegue.
El clúster debe contener el Secret `dockerhub-cred` en el namespace `sa-p7` para descargar las imágenes privadas de Docker Hub.

## Evidencias pendientes

Guardar en `P7/evidencias` las capturas de build, test 70/20/10, imágenes publicadas y despliegue exitoso.
