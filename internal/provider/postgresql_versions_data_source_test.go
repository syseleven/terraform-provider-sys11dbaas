package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestPostgresqlVersionsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "sys11dbaas_postgresql_versions" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of versions returned
					resource.TestCheckResourceAttr("data.sys11dbaas_postgresql_versions.test", "versions.#", "31"),

					// Verify the first version to ensure all attributes are set
					resource.TestCheckResourceAttr("data.sys11dbaas_postgresql_versions.test", "versions.0.id", "15.6"),
					resource.TestCheckResourceAttr("data.sys11dbaas_postgresql_versions.test", "versions.0.description", ""),
					resource.TestCheckResourceAttr("data.sys11dbaas_postgresql_versions.test", "versions.0.default", "false"),
				),
			},
		},
	})
}
