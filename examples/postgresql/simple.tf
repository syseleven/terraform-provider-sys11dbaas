resource "sys11dbaas_database" "postgresql" {
  name = "example-postgresql"
  application_config = {
    instances = 1
    type      = "postgresql"
    version   = 16.2
  }
  service_config = {
    disksize = 25
    flavor   = "m2.small"
    region   = "dus2"
    type     = "database"
  }
}
