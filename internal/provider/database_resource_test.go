package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDatabaseResource(t *testing.T) {
	resourceName := acctest.RandomWithPrefix("create_read_v1")
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
resource "sys11dbaas_database" "test" {
  name = "%s"
  application_config = {
    instances = 1
    type      = "postgresql"
    version   = 17.4
    password = "test_test_test_test"
  }

  service_config = {
    disksize   = 25
    flavor     = "SCS-2V-4-50n"
    region     = "dus2"
  }
}
`, resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sys11dbaas_database.test", "name", resourceName),
					resource.TestCheckResourceAttr("sys11dbaas_database.test", "application_config.instances", "1"),
					resource.TestCheckResourceAttr("sys11dbaas_database.test", "application_config.type", "postgresql"),
					resource.TestCheckResourceAttr("sys11dbaas_database.test", "application_config.version", "17.4"),

					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("sys11dbaas_database.test", "uuid"),
					resource.TestCheckResourceAttrSet("sys11dbaas_database.test", "created_at"),
					resource.TestCheckResourceAttrSet("sys11dbaas_database.test", "created_by"),
					resource.TestCheckResourceAttrSet("sys11dbaas_database.test", "last_modified_at"),
					resource.TestCheckResourceAttrSet("sys11dbaas_database.test", "last_modified_by"),
				),
			},
			// Update and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
resource "sys11dbaas_database" "test" {
  name = "%s"
  application_config = {
    instances = 1
    type      = "postgresql"
    version = 17.4
    password = "test_test_test_test"
  }
  	
  service_config = {
    disksize   = 25
    flavor     = "SCS-2V-4-50n"
    region     = "dus2"
    remote_ips = [
      "0.0.0.0/0"
    ]
  }
}
`, resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sys11dbaas_database.test", "service_config.remote_ips.0", "0.0.0.0/0"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
