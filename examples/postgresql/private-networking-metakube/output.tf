output "db-hostname" {
  description = "Hostname of the database"
  value       = resource.sys11dbaas_database.db.application_config.private_networking.hostname
}
