output "db" {
  value = [resource.sys11dbaas_database.db.uuid, resource.sys11dbaas_database.db.status, resource.sys11dbaas_database.db.phase, resource.sys11dbaas_database.db.resource_status]
}
