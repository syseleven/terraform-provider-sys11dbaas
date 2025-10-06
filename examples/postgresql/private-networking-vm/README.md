# Connecting DBaaS via private networking to a VM

This example demostrates how to connect a DBaaS instance to a VM using private networking

## Prerequisites

- Terraform
- SSH Client
- SysEleven IAM token with the following permissions:
  - OpenStack: Editor
  - DBaaS: Create/Delete/Read/Update Databases

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

terraform init
terraform apply
```

## Verifying functionality

```bash
# Connect to the VM. IP address is given by Terraform output
ssh -i ssh_private.key -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null ubuntu@${VM_IP_ADDRESS}

# Verify database connection. Hostname is given by Terraform output aswell.
psql -h ${DBAAS_HOSTNAME} -U admin
```

## Remove example resources

```bash
terraform destroy
```
