# Terraform de P8

Aqui se agrega unicamente la infraestructura propia de P8:

- Namespace `sa-p8`.
- `ResourceQuota`.
- `LimitRange`.
- Cuentas de servicio y RBAC con privilegios minimos.

## Uso local

```bash
terraform init
terraform fmt -check
terraform validate
terraform plan -var-file=terraform.tfvars
terraform apply -var-file=terraform.tfvars
```

`terraform.tfvars` es local y no debe confirmarse. El plan y el apply
requieren un kubeconfig proporcionado fuera del repositorio. La aplicacion y
actualizacion de manifiestos de los microservicios no se realiza desde
Terraform ni desde GitHub Actions; esa responsabilidad corresponde a ArgoCD.
