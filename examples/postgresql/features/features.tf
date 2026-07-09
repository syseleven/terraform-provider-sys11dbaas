data "sys11dbaas_features" "available" {
  type = "postgresql"
}

resource "sys11dbaas_database" "postgresql" {
  name = "example-postgresql"
  application_config = {
    instances = 1
    type      = "postgresql"
    version   = 17.5
    password  = "veryS3cretPassword"

    # The available feature set may change over time. Replace the example
    # key below with an actual feature name from the sys11dbaas_features
    # data source. Values must be "on" or "off".
    # Example:  data.sys11dbaas_features.available.features[0].name
    features = {
      example_feature = "on"
    }
  }
  service_config = {
    disksize = 25
    flavor   = "SCS-2V-4-50n"
    region   = "dus2"
  }
}

output "effective_features" {
  value = sys11dbaas_database.postgresql.application_config.effective_features
}
