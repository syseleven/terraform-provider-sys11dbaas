resource "sys11dbaas_database" "postgresql" {
  name = "example-postgresql"
  application_config = {
    instances = 1
    type      = "postgresql"
    version   = 16.2
    password  = "veryS3cretPassword"
  }
  service_config = {
    disksize = 25
    flavor   = "m2c.small"
    region   = "dus2"
    type     = "database"
  }
}
