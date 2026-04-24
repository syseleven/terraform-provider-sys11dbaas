package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestDatabaseResource(t *testing.T) {
	resourceName := acctest.RandomWithPrefix("create_read")
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
    public_networking = {
      enabled            = true
      allowed_cidrs = [
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
					resource.TestCheckResourceAttr("sys11dbaas_database.test", "name", resourceName),
					resource.TestCheckResourceAttr("sys11dbaas_database.test", "application_config.instances", "1"),
					resource.TestCheckResourceAttr("sys11dbaas_database.test", "application_config.type", "postgresql"),
					resource.TestCheckResourceAttr("sys11dbaas_database.test", "application_config.version", "17.4"),
					resource.TestCheckResourceAttr("sys11dbaas_database.test", "application_config.public_networking.allowed_cidrs.0", "0.0.0.0/0"),
					resource.TestCheckResourceAttr("sys11dbaas_database.test", "status", "ClusterIsReady"),

					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("sys11dbaas_database.test", "uuid"),
					resource.TestCheckResourceAttrSet("sys11dbaas_database.test", "created_at"),
					resource.TestCheckResourceAttrSet("sys11dbaas_database.test", "created_by"),
					resource.TestCheckResourceAttrSet("sys11dbaas_database.test", "last_modified_at"),
					resource.TestCheckResourceAttrSet("sys11dbaas_database.test", "last_modified_by"),
				),
			},
			// ImportState testing
			{
				ResourceName: "sys11dbaas_database.test",
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					resource, ok := s.RootModule().Resources["sys11dbaas_database.test"]
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
resource "sys11dbaas_database" "test" {
  name = "%s"
  application_config = {
    instances = 1
    type      = "postgresql"
    version = 17.4
    password = "test_test_test_test"
    public_networking = {
      enabled            = true
      allowed_cidrs      = [
        "1.1.1.1/32"
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
					resource.TestCheckResourceAttr("sys11dbaas_database.test", "application_config.public_networking.enabled", "true"),
					resource.TestCheckResourceAttr("sys11dbaas_database.test", "application_config.public_networking.allowed_cidrs.0", "1.1.1.1/32"),
					resource.TestCheckResourceAttr("sys11dbaas_database.test", "status", "ClusterIsReady"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestDatabaseResourceWithNetworkMigration(t *testing.T) {
	resourceName := acctest.RandomWithPrefix("migrate_network")
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
    remote_ips = ["0.0.0.0/0"]
  }
}
`, resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sys11dbaas_database.test", "name", resourceName),
					resource.TestCheckResourceAttr("sys11dbaas_database.test", "application_config.instances", "1"),
					resource.TestCheckResourceAttr("sys11dbaas_database.test", "application_config.type", "postgresql"),
					resource.TestCheckResourceAttr("sys11dbaas_database.test", "application_config.version", "17.4"),
					resource.TestCheckNoResourceAttr("sys11dbaas_database.test", "application_config.public_networking.#"),
					resource.TestCheckResourceAttr("sys11dbaas_database.test", "service_config.remote_ips.0", "0.0.0.0/0"),
					resource.TestCheckResourceAttr("sys11dbaas_database.test", "status", "ClusterIsReady"),

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

    public_networking = {
        enabled = true
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
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sys11dbaas_database.test", "application_config.public_networking.enabled", "true"),
					resource.TestCheckResourceAttr("sys11dbaas_database.test", "application_config.public_networking.allowed_cidrs.0", "0.0.0.0/0"),
					resource.TestCheckNoResourceAttr("sys11dbaas_database.test", "service_config.remote_ips"),
					resource.TestCheckResourceAttr("sys11dbaas_database.test", "status", "ClusterIsReady"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestDatabaseResourceWithLegacyNetworking(t *testing.T) {
	resourceName := acctest.RandomWithPrefix("create_read")
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
    remote_ips = ["0.0.0.0/0"]
  }
}
`, resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sys11dbaas_database.test", "name", resourceName),
					resource.TestCheckResourceAttr("sys11dbaas_database.test", "application_config.instances", "1"),
					resource.TestCheckResourceAttr("sys11dbaas_database.test", "application_config.type", "postgresql"),
					resource.TestCheckResourceAttr("sys11dbaas_database.test", "application_config.version", "17.4"),
					resource.TestCheckResourceAttr("sys11dbaas_database.test", "service_config.remote_ips.0", "0.0.0.0/0"),
					resource.TestCheckResourceAttr("sys11dbaas_database.test", "status", "ClusterIsReady"),

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
    remote_ips = ["1.1.1.1/32"]
  }
}
`, resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sys11dbaas_database.test", "service_config.remote_ips.0", "1.1.1.1/32"),
					resource.TestCheckResourceAttr("sys11dbaas_database.test", "status", "ClusterIsReady"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestDatabaseCheckCompatibilityWithRelease(t *testing.T) {
	resourceName := acctest.RandomWithPrefix("release_compatibility")
	resource.ParallelTest(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"sys11dbaas": {
						Source:            "syseleven/sys11dbaas",
						VersionConstraint: "0.3.4",
					},
				},
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
    remote_ips = ["0.0.0.0/0"]
  }
}
`, resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sys11dbaas_database.test", "name", resourceName),
					resource.TestCheckResourceAttr("sys11dbaas_database.test", "application_config.instances", "1"),
					resource.TestCheckResourceAttr("sys11dbaas_database.test", "application_config.type", "postgresql"),
					resource.TestCheckResourceAttr("sys11dbaas_database.test", "application_config.version", "17.4"),
					resource.TestCheckNoResourceAttr("sys11dbaas_database.test", "service_config.public_networking"),
					resource.TestCheckResourceAttr("sys11dbaas_database.test", "service_config.remote_ips.0", "0.0.0.0/0"),
					resource.TestCheckResourceAttr("sys11dbaas_database.test", "status", "ClusterIsReady"),

					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("sys11dbaas_database.test", "uuid"),
					resource.TestCheckResourceAttrSet("sys11dbaas_database.test", "created_at"),
					resource.TestCheckResourceAttrSet("sys11dbaas_database.test", "created_by"),
					resource.TestCheckResourceAttrSet("sys11dbaas_database.test", "last_modified_at"),
					resource.TestCheckResourceAttrSet("sys11dbaas_database.test", "last_modified_by"),
				),
			},
			{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
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
    remote_ips = ["1.1.1.1/32"]
  }
}
`, resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr("sys11dbaas_database.test", "service_config.public_networking"),
					resource.TestCheckResourceAttr("sys11dbaas_database.test", "service_config.remote_ips.0", "1.1.1.1/32"),
					resource.TestCheckResourceAttr("sys11dbaas_database.test", "status", "ClusterIsReady"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestDatabaseResourceV2Move(t *testing.T) {
	resourceName := acctest.RandomWithPrefix("move_v2")
	v2Config := fmt.Sprintf(`
resource "sys11dbaas_database_v2" "test" {
  name = "%s"
  application_config = {
    instances = 1
    type      = "postgresql"
    version   = 17.4
    password = "test_test_test_test"

    public_networking = {
         enabled = true
         allowed_cidrs = ["0.0.0.0/0"]
    }
  }

  service_config = {
    disksize   = 25
    flavor     = "SCS-2V-4-50n"
    region     = "dus2"
  }
}
`, resourceName)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + v2Config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sys11dbaas_database_v2.test", "name", resourceName),
					resource.TestCheckResourceAttr("sys11dbaas_database_v2.test", "application_config.instances", "1"),
					resource.TestCheckResourceAttr("sys11dbaas_database_v2.test", "application_config.type", "postgresql"),
					resource.TestCheckResourceAttr("sys11dbaas_database_v2.test", "application_config.version", "17.4"),
					resource.TestCheckResourceAttr("sys11dbaas_database_v2.test", "status", "ClusterIsReady"),

					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("sys11dbaas_database_v2.test", "uuid"),
					resource.TestCheckResourceAttrSet("sys11dbaas_database_v2.test", "created_at"),
					resource.TestCheckResourceAttrSet("sys11dbaas_database_v2.test", "created_by"),
					resource.TestCheckResourceAttrSet("sys11dbaas_database_v2.test", "last_modified_at"),
					resource.TestCheckResourceAttrSet("sys11dbaas_database_v2.test", "last_modified_by"),
				),
			},
			// Move
			{
				Config: providerConfig + fmt.Sprintf(`
resource "sys11dbaas_database" "test" {
  name = "%s"
  application_config = {
    instances = 1
    type      = "postgresql"
    version   = 17.4
    password = "test_test_test_test"

    public_networking = {
         enabled = true
         allowed_cidrs = ["0.0.0.0/0"]
    }
  }

  service_config = {
    disksize   = 25
    flavor     = "SCS-2V-4-50n"
    region     = "dus2"
  }
}

moved {
    from = sys11dbaas_database_v2.test
    to = sys11dbaas_database.test
}
`, resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sys11dbaas_database.test", "name", resourceName),
					resource.TestCheckResourceAttr("sys11dbaas_database.test", "application_config.instances", "1"),
					resource.TestCheckResourceAttr("sys11dbaas_database.test", "application_config.type", "postgresql"),
					resource.TestCheckResourceAttr("sys11dbaas_database.test", "application_config.version", "17.4"),
					resource.TestCheckResourceAttr("sys11dbaas_database.test", "application_config.public_networking.allowed_cidrs.0", "0.0.0.0/0"),
					resource.TestCheckResourceAttr("sys11dbaas_database.test", "status", "ClusterIsReady"),

					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("sys11dbaas_database.test", "uuid"),
					resource.TestCheckResourceAttrSet("sys11dbaas_database.test", "created_at"),
					resource.TestCheckResourceAttrSet("sys11dbaas_database.test", "created_by"),
					resource.TestCheckResourceAttrSet("sys11dbaas_database.test", "last_modified_at"),
					resource.TestCheckResourceAttrSet("sys11dbaas_database.test", "last_modified_by"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
