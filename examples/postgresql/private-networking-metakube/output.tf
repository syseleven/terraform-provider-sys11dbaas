output "db-hostname" {
  description = "Hostname of the database"
  value       = resource.sys11dbaas_database_v2.db.application_config.private_networking.hostname
}
