# Connecting DBaaS via private networking to a VM

This example demostrates how to connect a DBaaS instance to a VM using private networking

## Prerequisites

- Terraform
- SSH Client
- SSH public key

## Deploying the example

```bash
terraform init
terraform apply
```

## Verifying functionality

```bash
# Connect to the VM. IP address is given by Terraform output
ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null ubuntu@${IP_ADDRESS}

# Verify database connection. Hostname is given by Terraform output aswell.

```
