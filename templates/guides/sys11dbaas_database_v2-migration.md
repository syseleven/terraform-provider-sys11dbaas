---
subcategory: ""
page_title: "Migration from sys11dbaas_database_v2 to sys11dbaas_database"
description: |-
  Migration from sys11dbaas_database_v2 to sys11dbaas_database
---

The terraform provider version 0.4.0 deprecated the usage of sys11dbaas_database_v2 by making sys11dbaas_database
forward compatible. This means customers are expected to consume sys11dbaas_database from now on, with internal api
versions hidden in implementation. The sys11dbaas_database_v2 resource will be removed with the release of the next
major version of the terraform provider itself.

There are two supported ways to migrate from `sys11dbaas_database_v2` to
`sys11dbaas_database`. Before starting the migration ensure you're using a version a provider version >= 0.4.0.

## Method 1: Terraform `moved` block (Recommended)

Add a `moved` block, rename the resource type, and run `terraform apply`:

```hcl
moved {
  from = sys11dbaas_database_v2.my_database
  to   = sys11dbaas_database.my_database
}

resource "sys11dbaas_database" "my_database" {
  name = "example-postgresql"

  application_config = {
    instances = 1
    type      = "postgresql"
    version   = "17.5"
  }

  service_config = {
    disksize = 25
    flavor   = "SCS-2V-4-50n"
    region   = "dus2"
  }
}
```

Terraform will automatically migrate the state to the new resource type using
the provider's built-in MoveState transformer.
For more information on the `moved` block, see the [official Terraform
documentation](https://developer.hashicorp.com/terraform/language/move).

## Method 2: Import-based migration

If the `moved` block approach doesn't work for your setup:

1. Add the new resource definition in your configuration:
   ```hcl
   resource "sys11dbaas_database" "my_database" {
     name = "example-postgresql"
     # ... copy over your configuration ...
   }
   ```

2. Import the existing database by UUID:
   ```
   $ terraform import sys11dbaas_database.my_database <database-uuid>
   ```

3. Remove the old resource from state:
   ```
   $ terraform state rm sys11dbaas_database_v2.my_database
   ```
   **Warning:** This step is invasive. Back up your state file before proceeding.

4. Verify:
   ```
   $ terraform plan
   ```
