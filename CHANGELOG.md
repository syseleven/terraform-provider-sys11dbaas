# CHANGELOG

## 1.0.0 (Unreleased)

## 0.4.0

### NOTES

The `sys11dbaas_database_v2` resource is deprecated and will be removed in the next major version. Migrate to `sys11dbaas_database` using the `moved` block:

```hcl
moved {
  from = sys11dbaas_database_v2.my_database
  to   = sys11dbaas_database.my_database
}
```

### FEATURES

* new data source `sys11dbaas_postgresql_flavors`, `sys11dbaas_postgresql_regions`, `sys11dbaas_postgresql_versions` ([#12](https://github.com/syseleven/terraform-provider-sys11dbaas/pull/12))
* new data source `sys11dbaas_features` for discovering available feature toggles ([#41](https://github.com/syseleven/terraform-provider-sys11dbaas/pull/41))
* `sys11dbaas_database` now supports `application_config.features` for enabling feature toggles ([#41](https://github.com/syseleven/terraform-provider-sys11dbaas/pull/41))
* `sys11dbaas_database` can now be imported by UUID ([#17](https://github.com/syseleven/terraform-provider-sys11dbaas/pull/17))
* MoveState support for automatic migration from `sys11dbaas_database_v2` via Terraform's `moved` block ([#17](https://github.com/syseleven/terraform-provider-sys11dbaas/pull/17))
* new example for private networking ([d958e4d](https://github.com/syseleven/terraform-provider-sys11dbaas/commit/d958e4df83edaaa68efccadf28f4e7fb9d8bf785))
* new example for feature toggles ([#41](https://github.com/syseleven/terraform-provider-sys11dbaas/pull/41))

### IMPROVEMENTS

* `sys11dbaas_database_v2` is now deprecated; migrate to `sys11dbaas_database` ([#17](https://github.com/syseleven/terraform-provider-sys11dbaas/pull/17))
* `service_config.remote_ips` is now deprecated; use `application_config.public_networking` instead ([#17](https://github.com/syseleven/terraform-provider-sys11dbaas/pull/17))
* upgraded to latest sys11dbaas SDK with functional-options client ([#8](https://github.com/syseleven/terraform-provider-sys11dbaas/pull/8))
* resource lifecycle now waits for both `Status == "Ready"` and `ResourceStatus == "Synced"` before completing create/read/update ([#8](https://github.com/syseleven/terraform-provider-sys11dbaas/pull/8)) ([#17](https://github.com/syseleven/terraform-provider-sys11dbaas/pull/17))
* `created_at` and `last_modified_at` now use RFC 3339 time types ([#8](https://github.com/syseleven/terraform-provider-sys11dbaas/pull/8))
* config validators for networking activation hints ([#17](https://github.com/syseleven/terraform-provider-sys11dbaas/pull/17))
* `maintenance_window` timezone documented as UTC ([13639b0](https://github.com/syseleven/terraform-provider-sys11dbaas/commit/13639b081a4e79d0b84ede24a9ac77aa6bc14233))
* added import and migration documentation ([#17](https://github.com/syseleven/terraform-provider-sys11dbaas/pull/17))
* comprehensive acceptance test suite for resources and data sources ([#10](https://github.com/syseleven/terraform-provider-sys11dbaas/pull/10))
* migrated CI from GitLab CI to GitHub Actions ([#4](https://github.com/syseleven/terraform-provider-sys11dbaas/pull/4))
* Go 1.24+ with native toolchain support ([#11](https://github.com/syseleven/terraform-provider-sys11dbaas/pull/11))
* added automated dependency management via Renovate ([#11](https://github.com/syseleven/terraform-provider-sys11dbaas/pull/11))

### BUG FIXES

* fixed nil pointer dereference for `MaintenanceWindow`, `ScheduledBackups`, `PrivateNetworking` and `PublicNetworking`

## Migration from sys11dbaas_database_v2

There are two supported ways to migrate from `sys11dbaas_database_v2` to
`sys11dbaas_database`.

### Method 1: Terraform `moved` block (Recommended)

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

### Method 2: Import-based migration

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

