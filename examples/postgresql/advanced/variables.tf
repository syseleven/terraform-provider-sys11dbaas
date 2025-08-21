# Datacenter / Project Configuration

variable "api_key" {
  type    = string
  default = ""
}

variable "api_url" {
  type    = string
  default = "https://dbaas.apis.syseleven.de"
}

variable "project" {
  type    = string
  default = ""
}

variable "region" {
  type    = string
  default = "dus2"
}

variable "org" {
  type    = string
  default = ""
}

# PostgreSQL DB Configuration

variable "db_name" {
  type    = string
  default = "postgres-db"
}

variable "db_description" {
  type    = string
  default = "My Postgres DB"
}

variable "db_password" {
  type    = string
  default = ""
}

variable "db_type" {
  type    = string
  default = "postgresql"
}

variable "db_version" {
  type    = number
  default = 17.5
}

variable "db_disk_size" {
  type    = number
  default = 25
}

variable "db_flavor" {
  type    = string
  default = "SCS-2V-4-50n"
}

variable "db_instances" {
  type    = number
  default = 2
}

variable "db_backup_hour" {
  type    = number
  default = 4
}

variable "db_remote_ips" {
  type    = list(string)
  default = ["0.0.0.0/0"]
}
