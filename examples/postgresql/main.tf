terraform {
  required_version = ">= 1.8"
  required_providers {
    sys11dbaas = {
      source  = "registry.terraform.io/syseleven/sys11dbaas"
      version = "~> 1"
    }
  }
}

variable "api_key" {
  type = string
}

variable "api_url" {
  type = string
}

variable "project" {
  type = string
}

variable "org" {
  type = string
}

provider "sys11dbaas" {
  url               = var.api_url
  api_key           = var.api_key
  project           = var.project
  organization      = var.org
  wait_for_creation = "true"
}

resource "sys11dbaas_database" "my_first_tf_db" {
  name        = "my-first-terraform-db"
  description = "this is my first terraform db"
  application_config = {
    instances = 3
    password  = "veryS3cretPassword"
    type      = "postgresql"
    version   = 16.2
    scheduled_backups = {
      schedule = {
        hour = 4
      }
    }
  }
  service_config = {
    disksize = 25
    flavor   = "m2c.small"
    region   = "dus2"
    type     = "database"
    remote_ips = [
      "0.0.0.0/0"
    ]
  }
}

output "my_first_tf_db" {
  value = [resource.sys11dbaas_database.my_first_tf_db.uuid, resource.sys11dbaas_database.my_first_tf_db.status, resource.sys11dbaas_database.my_first_tf_db.phase, resource.sys11dbaas_database.my_first_tf_db.resource_status]
}
