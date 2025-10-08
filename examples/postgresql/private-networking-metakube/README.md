# Connecting DBaaS via private networking to a MetaKube cluster

This example demostrates how to connect a DBaaS instance to a MetaKube cluster using private networking

## Prerequisites

- Terraform
- kubectl
- SysEleven IAM token with the following permissions:
  - OpenStack: Editor
  - DBaaS: Create/Delete/Read/Update Databases
  - MetaKube: Create/Delete/Read/Update Clusters

## Deploying the example

```bash
export SYS11DBAAS_API_KEY=<SysEleven IAM token>
export SYS11DBAAS_ORGANIZATION=<SysEleven IAM organization UUID>
export SYS11DBAAS_PROJECT=<SysEleven IAM project ID>

export OS_AUTH_URL=https://keystone.cloud.syseleven.net:5000/v3/
export OS_AUTH_TYPE=v3applicationcredential
export OS_APPLICATION_CREDENTIAL_ID=s11auth
export OS_APPLICATION_CREDENTIAL_SECRET=$SYS11DBAAS_API_KEY
export OS_IDENTITY_API_VERSION=3

export METAKUBE_TOKEN=$SYS11DBAAS_API_KEY

terraform init
terraform apply
```

## Verifying functionality

```bash
KUBECONFIG=./kubeconfig  kubectl run --rm=true --restart=Never --stdin --tty dbaas-connector --image=alpine/psql --command -- sh
psql -U admin -h ${DBAAS_HOSTNAME}
```

## Remove example resources

```bash
terraform destroy
```
