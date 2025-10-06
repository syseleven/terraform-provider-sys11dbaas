variable "region" {
  type        = string
  description = "In which region the resources should be deployed to."
}

variable "ssh_publickey" {
  type        = string
  description = "ssh-rsa public key in authorized_keys format (ssh-rsa AAAAB3Nz [...] ABAAACAC62Lw== user@host)"
}

variable "database_password" {
  type        = string
  sensitive   = true
  description = "Password for the database 'admin' user."
}
