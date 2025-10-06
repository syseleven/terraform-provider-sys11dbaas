output "vm-ip" {
  description = "IP address of the VM"
  value       = openstack_networking_floatingip_v2.fip.address
}

output "db-hostname" {
  description = "Hostname of the database"
  value       = resource.sys11dbaas_database_v2.db.application_config.private_networking.hostname
}
