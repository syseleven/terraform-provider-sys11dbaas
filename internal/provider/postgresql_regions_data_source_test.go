package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestPostgresqlRegionsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "sys11dbaas_postgresql_regions" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of regions returned
					resource.TestCheckResourceAttr("data.sys11dbaas_postgresql_regions.test", "regions.#", "2"),

					// Verify the first region to ensure all attributes are set
					resource.TestCheckResourceAttr("data.sys11dbaas_postgresql_regions.test", "regions.0.id", "dus2"),
					resource.TestCheckResourceAttr("data.sys11dbaas_postgresql_regions.test", "regions.0.description", ""),
					resource.TestCheckResourceAttr("data.sys11dbaas_postgresql_regions.test", "regions.0.default", "false"),
				),
			},
		},
	})
}
