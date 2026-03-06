package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestPostgresqlFlavorsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "sys11dbaas_postgresql_flavors" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of flavors returned
					resource.TestCheckResourceAttr("data.sys11dbaas_postgresql_flavors.test", "flavors.#", "7"),

					// Verify the first flavor to ensure all attributes are set
					resource.TestCheckResourceAttr("data.sys11dbaas_postgresql_flavors.test", "flavors.0.id", "SCS-2V-4-50n"),
					resource.TestCheckResourceAttr("data.sys11dbaas_postgresql_flavors.test", "flavors.0.description", "2/4/50"),
					resource.TestCheckResourceAttr("data.sys11dbaas_postgresql_flavors.test", "flavors.0.default", "true"),
				),
			},
		},
	})
}
