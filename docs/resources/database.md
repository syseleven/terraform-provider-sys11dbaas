---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sys11dbaas_database Resource - terraform-provider-sys11dbaas"
subcategory: ""
description: |-
  
---

# sys11dbaas_database (Resource)



## Example Usage

```terraform
terraform {
  required_providers {
    sys11dbaas = {
      source = "syseleven/sys11dbaas"
    }
  }
}

provider "sys11dbaas" {
  api_key      = "s11_prak_..."
  project      = "0123456789"
  organization = "0123-456-78-9"
}

resource "sys11dbaas_database" "postgresql" {
  name = "example-postgresql"
  application_config = {
    instances = 1
    type      = "postgresql"
    version   = 17.5
    password  = "veryS3cretPassword"
  }
  service_config = {
    disksize = 25
    flavor   = "SCS-2V-4-50n"
    region   = "dus2"
  }
}

output "db" {
  value = [resource.sys11dbaas_database.postgresql.uuid, resource.sys11dbaas_database.postgresql.status]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `application_config` (Attributes) (see [below for nested schema](#nestedatt--application_config))
- `name` (String) Name of the database.
- `service_config` (Attributes) (see [below for nested schema](#nestedatt--service_config))

### Optional

- `description` (String) Fulltext description of the database.

### Read-Only

- `created_at` (String) Date when the database was created.
- `created_by` (String) Initial creator of the database.
- `last_modified_at` (String) Date when the database was last modified.
- `last_modified_by` (String) User who last changed the database.
- `phase` (String) Detailed status of the database.
- `resource_status` (String) Sync status of the database.
- `status` (String) Overall status of the database.
- `uuid` (String) UUID of the database.

<a id="nestedatt--application_config"></a>
### Nested Schema for `application_config`

Required:

- `instances` (Number) Node count of the database cluster.
- `type` (String) Type of the database. Currently only supports 'postgresql'.
- `version` (String) Minor version of PostgreSQL.

Optional:

- `password` (String, Sensitive) Password for the admin user.
- `recovery` (Attributes) (see [below for nested schema](#nestedatt--application_config--recovery))
- `scheduled_backups` (Attributes) Scheduled backups policy for the database. (see [below for nested schema](#nestedatt--application_config--scheduled_backups))

Read-Only:

- `hostname` (String) DNS name of the database in the format uuid.postgresql.syseleven.services.
- `ip_address` (String) Public IP address of the database. It will be 'Pending' if no address has been assigned yet.

<a id="nestedatt--application_config--recovery"></a>
### Nested Schema for `application_config.recovery`

Optional:

- `exclusive` (Boolean) Set to true, when the given target should be excluded.
- `source` (String) UUID of the source database.
- `target_lsn` (String) LSN of the write-ahead log location up to which recovery will proceed. target_* parameters are mutually exclusive.
- `target_name` (String) Named restore point (created with pg_create_restore_point()) to which recovery will proceed. target_* parameters are mutually exclusive.
- `target_time` (String) Time stamp up to which recovery will proceed, expressed in RFC 3339 format. target_* parameters are mutually exclusive.
- `target_xid` (String) Transaction ID up to which recovery will proceed. target_* parameters are mutually exclusive.


<a id="nestedatt--application_config--scheduled_backups"></a>
### Nested Schema for `application_config.scheduled_backups`

Optional:

- `retention` (Number) Duration in days for which backups should be stored.
- `schedule` (Attributes) Schedules for the backup policy. (see [below for nested schema](#nestedatt--application_config--scheduled_backups--schedule))

<a id="nestedatt--application_config--scheduled_backups--schedule"></a>
### Nested Schema for `application_config.scheduled_backups.schedule`

Optional:

- `hour` (Number) Hour when the full backup should start. If this value is omitted, a random hour between 1am and 5am will be generated.
- `minute` (Number) Minute when the full backup should start. If this value is omitted, a random minute will be generated.




<a id="nestedatt--service_config"></a>
### Nested Schema for `service_config`

Required:

- `disksize` (Number) Disksize in GB.
- `flavor` (String) VM flavor to use.
- `region` (String) Region for the database.

Optional:

- `maintenance_window` (Attributes) Maintenance window. This will be a time window for updates and maintenance. If omitted, a random window will be generated. (see [below for nested schema](#nestedatt--service_config--maintenance_window))
- `remote_ips` (List of String) List of IP addresses, that should be allowed to connect to the database.
- `type` (String) Type of the service you want to create (default `database`)

<a id="nestedatt--service_config--maintenance_window"></a>
### Nested Schema for `service_config.maintenance_window`

Optional:

- `day_of_week` (Number) Day of week as a cron time (0=Sun, 1=Mon, ..., 6=Sat). If omitted, a random day will be used.
- `start_hour` (Number) Hour when the maintenance window starts. If omitted, a random hour between 20 and 4 will be used.
- `start_minute` (Number) Minute when the maintenance window starts. If omitted, a random minute will be used.
