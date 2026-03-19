package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestDatabaseResourceV2(t *testing.T) {
	resourceName := acctest.RandomWithPrefix("create_read_v2")
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
resource "sys11dbaas_database_v2" "test" {
  name = "%s"
  application_config = {
    instances = 1
    type      = "postgresql"
    version = 17.4
    password = "test_test_test_test"
    public_networking = {
      enabled            = true
    }
  }
  	
  service_config = {
    disksize   = 25
    flavor     = "SCS-2V-4-50n"
    region     = "dus2"
  }
}
`, resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sys11dbaas_database_v2.test", "name", resourceName),
					resource.TestCheckResourceAttr("sys11dbaas_database_v2.test", "application_config.instances", "1"),
					resource.TestCheckResourceAttr("sys11dbaas_database_v2.test", "application_config.type", "postgresql"),
					resource.TestCheckResourceAttr("sys11dbaas_database_v2.test", "application_config.version", "17.4"),

					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("sys11dbaas_database_v2.test", "uuid"),
					resource.TestCheckResourceAttrSet("sys11dbaas_database_v2.test", "created_at"),
					resource.TestCheckResourceAttrSet("sys11dbaas_database_v2.test", "created_by"),
					resource.TestCheckResourceAttrSet("sys11dbaas_database_v2.test", "last_modified_at"),
					resource.TestCheckResourceAttrSet("sys11dbaas_database_v2.test", "last_modified_by"),
				),
			},
			// ImportState testing
			{
				ResourceName: "sys11dbaas_database_v2.test",
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					resource, ok := s.RootModule().Resources["sys11dbaas_database_v2.test"]
					if !ok {
						return "", fmt.Errorf("failed to find test resource in state")
					}
					return resource.Primary.Attributes["uuid"], nil
				},
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "uuid",
				ImportStateVerifyIgnore:              []string{"application_config.password"},
			},
			// Update and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
resource "sys11dbaas_database_v2" "test" {
  name = "%s"
  application_config = {
    instances = 1
    type      = "postgresql"
    version = 17.4
    password = "test_test_test_test"
    public_networking = {
      enabled            = true
      allowed_cidrs      = [
        "0.0.0.0/0"
      ]
    }
  }
  	
  service_config = {
    disksize   = 25
    flavor     = "SCS-2V-4-50n"
    region     = "dus2"
  }
}
`, resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sys11dbaas_database_v2.test", "application_config.public_networking.allowed_cidrs.0", "0.0.0.0/0"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
