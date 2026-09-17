# Terraform de P8

Aqui se define la infraestructura de P8 y el clúster GKE:

- Clúster zonal GKE y node pool principal.
- Namespace `sa-p8`.
- `ResourceQuota`.
- `LimitRange`.
- Cuentas de servicio y RBAC con privilegios minimos.
- Secreto fuente administrado por Terraform con variables sensibles locales.
- Workload Identity y acceso al plano de control limitado a una red autorizada.

## Uso local

```bash
terraform init
terraform fmt -check
terraform validate
terraform plan -var-file=terraform.tfvars
terraform apply -var-file=terraform.tfvars
```

`terraform.tfvars` es local y no debe confirmarse. Debe contener la IP publica
autorizada para administrar GKE y los valores sensibles `jwt_secret` y
`encryption_key`. Terraform administra la
infraestructura base; la aplicacion y actualizacion de manifiestos de los
microservicios corresponde exclusivamente a ArgoCD.
