package provider

import (
	"fmt"
	"terraform-provider-sys11dbaas/internal/testhelpers"
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

func TestDatabaseResourceV2ImportV1(t *testing.T) {
	var uuid *string
	resourceName := acctest.RandomWithPrefix("import_v1_to_v2")
	v1Config := fmt.Sprintf(`
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
    remote_ips = ["0.0.0.0/0"]
  }
}
`, resourceName)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + v1Config,
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
					func(s *terraform.State) error {
						resources := s.RootModule().Resources
						id := resources["sys11dbaas_database.test"].Primary.Attributes["uuid"]
						uuid = &id
						delete(resources, "sys11dbaas_database.test")
						return nil
					},
				),
			},
			// ImportState testing
			{
				ResourceName: "sys11dbaas_database_v2.test",
				Config: providerConfig + fmt.Sprintf(`
resource "sys11dbaas_database_v2" "test" {
  name = "%s"
  application_config = {
    instances = 1
    type      = "postgresql"
    version   = 17.4
    password = "test_test_test_test"
    public_networking = {
    	allowed_cidrs = ["0.0.0.0/0"]
    }
  }

  service_config = {
    disksize   = 25
    flavor     = "SCS-2V-4-50n"
    region     = "dus2"
  }
}
				`, resourceName),
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					return *uuid, nil
				},
				ImportState: true,
				ImportStateCheck: testhelpers.ComposeAggregateImportStateCheckFunc(
					// this is explicitly not using ImportCheckResourceAttr to avoid capturing the uuid at invocation time
					func(is []*terraform.InstanceState) error {
						if is[0].Attributes["uuid"] != *uuid {
							return fmt.Errorf("Expected matching uuid after import")
						}
						return nil
					},
					testhelpers.ImportCheckResourceAttr("name", resourceName),
					testhelpers.ImportCheckResourceAttr("application_config.instances", "1"),
					testhelpers.ImportCheckResourceAttr("application_config.type", "postgresql"),
					testhelpers.ImportCheckResourceAttr("application_config.version", "17.4"),
					testhelpers.ImportCheckResourceAttr("application_config.public_networking.enabled", "true"),
					testhelpers.ImportCheckResourceAttr("application_config.public_networking.allowed_cidrs.0", "0.0.0.0/0"),

					// Verify dynamic values have any value set in the state.
					testhelpers.ImportCheckResourceAttrSet("uuid"),
					testhelpers.ImportCheckResourceAttrSet("created_at"),
					testhelpers.ImportCheckResourceAttrSet("created_by"),
					testhelpers.ImportCheckResourceAttrSet("last_modified_at"),
					testhelpers.ImportCheckResourceAttrSet("last_modified_by"),
				),
				ImportStateVerify: false,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
