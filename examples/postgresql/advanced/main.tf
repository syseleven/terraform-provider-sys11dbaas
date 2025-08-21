provider "sys11dbaas" {
  api_key      = var.api_key
  project      = var.project
  organization = var.org
}

resource "sys11dbaas_database_v2" "db" {
  name        = var.db_name
  description = var.db_description
  application_config = {
    instances = var.db_instances
    password  = var.db_password
    type      = var.db_type
    version   = var.db_version

    private_networking   = {
      enabled            = var.private_networking
      shared_subnet_cidr = var.shared_subnet_cidr
      allowed_cidrs      = var.allowed_cidrs
    }

    public_networking = {
      enabled         = var.public_networking
      allowed_cidrs   = var.db_remote_ips
    }

    scheduled_backups = {
      schedule = {
        hour = var.db_backup_hour
      }
    }
  }

  service_config = {
    disksize     = var.db_disk_size
    flavor       = var.db_flavor
    region       = var.region
    remote_ips   = var.db_remote_ips
  }
}
