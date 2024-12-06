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
    version   = 16.2
    password  = "veryS3cretPassword"
  }
  service_config = {
    disksize = 25
    flavor   = "m2c.small"
    region   = "dus2"
  }
}

output "db" {
  value = [resource.sys11dbaas_database.postgresql.uuid, resource.sys11dbaas_database.postgresql.status]
}

