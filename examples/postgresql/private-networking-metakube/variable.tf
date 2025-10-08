variable "region" {
  type        = string
  description = "In which region the resources should be deployed to."
}

variable "database_password" {
  type        = string
  sensitive   = true
  description = "Password for the database 'admin' user."
}
