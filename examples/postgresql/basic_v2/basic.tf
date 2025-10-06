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

resource "sys11dbaas_database_v2" "postgresql" {
  name = "example-postgresql-v2"
  application_config = {
    instances = 1
    type      = "postgresql"
    version   = 17.5
    password  = "veryS3cretPassword2"
    private_networking = {
      enabled            = true
      shared_subnet_cidr = "10.245.0.0/24"
      allowed_cidrs      = ["10.10.50.0/24", "10.245.0.0/24"]
    }
  }
  service_config = {
    disksize   = 25
    flavor     = "SCS-2V-4-50n"
    region     = "dus2"
    remote_ips = ["176.74.56.225/26"]
  }
}

output "db-hostname" {
  value = resource.sys11dbaas_database_v2.db.application_config.private_networking.hostname
}

output "db-uuid" {
  value = resource.sys11dbaas_database_v2.postgresql.uuid
}

output "db-status" {
  value = resource.sys11dbaas_database_v2.postgresql.status
}

output "db-subnet-id" {
  value = resource.sys11dbaas_database_v2.postgresql.application_config.private_networking.shared_subnet_id
}

output "db-network-id" {
  value = resource.sys11dbaas_database_v2.postgresql.application_config.private_networking.shared_network_id
}
