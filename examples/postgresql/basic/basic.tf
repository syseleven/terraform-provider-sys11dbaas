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

resource "sys11dbaas_database_v2" "postgresql-v2" {
  name = "example-postgresql-v2"
  application_config = {
    instances = 1
    type      = "postgresql"
    version   = 17.5
    password  = "veryS3cretPassword2"
    private_network_config = {
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

output "db-v2-uuid" {
  value = resource.sys11dbaas_database_v2.postgresql-v2.uuid
}

output "db-v2-status" {
  value = resource.sys11dbaas_database_v2.postgresql-v2.status
}

output "db-v2-subnet" {
  value = resource.sys11dbaas_database_v2.postgresql-v2.application_config.private_network_config.shared_subnet_id
}

output "db-v2-network" {
  value = resource.sys11dbaas_database_v2.postgresql-v2.application_config.private_network_config.shared_network_id
}
