# Diagrama de la Práctica 9

```mermaid
flowchart TD
    U[Una orden: bootstrap.sh] --> B[Bucket GCS de Terraform]
    B --> T1[Terraform fase 1]
    T1 --> GKE[Cluster GKE P9 y nodos]
    GKE --> T2[Terraform fase 2]
    T2 --> R[Bucket Velero, cuentas GCP y Secret Manager]
    GKE --> A[Terraform instala ArgoCD]
    A --> ROOT[Application raíz p9-root]
    ROOT --> C1[Argo Rollouts]
    ROOT --> C2[External Secrets]
    ROOT --> C3[Kyverno]
    ROOT --> C4[Velero y respaldos programados]
    ROOT --> C5[Aplicación P9]
    C2 --> SM[Secret Manager externo]
    C4 --> VB[Bucket externo de respaldos]
    C5 --> DB[PostgreSQL StatefulSet y PVC]
    C5 --> ORD[orders-ms: 2 réplicas]
    ORD --> DB
    ORD --> RMQ[RabbitMQ]
    VB --> REST[Restauración después de reconstruir]
    REST --> C5
```

El repositorio de código reutiliza la infraestructura de `P8/terraform` como
módulo. El repositorio GitOps reutiliza los charts y manifiestos de P8 mediante
referencias remotas; P9 solo agrega los recursos de continuidad operativa.
