output "db-hostname" {
  value = resource.sys11dbaas_database_v2.db.application_config.private_networking.hostname
}

output "db-uuid" {
  value = resource.sys11dbaas_database_v2.db.uuid
}

output "db-status" {
  value = resource.sys11dbaas_database_v2.db.status
}

output "subnet-id" {
  value = resource.sys11dbaas_database_v2.db.application_config.private_networking.shared_subnet_id
}

output "network-id" {
  value = resource.sys11dbaas_database_v2.db.application_config.private_networking.shared_network_id
}
