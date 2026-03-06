package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	v2 "github.com/syseleven/sys11dbaas-sdk/database/v2"
)

// postgresqlFlavorsDataSourceModel maps the data source schema data.
type postgresqlFlavorsDataSourceModel struct {
	PostgresqlFlavors []PostgresqlFlavorModel `tfsdk:"flavors"`
}

// PostgresqlFlavorModel maps postgresqlFlavor schema data.
type PostgresqlFlavorModel struct {
	ID          types.String `tfsdk:"id"`
	Description types.String `tfsdk:"description"`
	Default     types.Bool   `tfsdk:"default"`
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &postgresqlFlavorsDataSource{}
	_ datasource.DataSourceWithConfigure = &postgresqlFlavorsDataSource{}
)

// NewPostgresqlFlavorsDataSource is a helper function to simplify the provider implementation.
func NewPostgresqlFlavorsDataSource() datasource.DataSource {
	return &postgresqlFlavorsDataSource{}
}

// coffeesDataSource is the data source implementation.
type postgresqlFlavorsDataSource struct {
	client       *v2.TypedClient
	project      types.String
	organization types.String
}

// Configure adds the provider configured client to the data source.
func (d *postgresqlFlavorsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerData, ok := req.ProviderData.(*sys11DBaaSProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *sys11DBaaSProviderData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = providerData.client.V2()
	d.organization = providerData.organization
	d.project = providerData.project
}

// Metadata returns the data source type name.
func (d *postgresqlFlavorsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_postgresql_flavors"
}

// Schema defines the schema for the data source.
func (d *postgresqlFlavorsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the list of available flavors.",
		Attributes: map[string]schema.Attribute{
			"flavors": schema.ListNestedAttribute{
				Description: "List of flavors.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Textual identifier of the flavor.",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "Description of the flavor.",
							Computed:    true,
						},
						"default": schema.BoolAttribute{
							Description: "In case this flavor is the default, this field is true.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *postgresqlFlavorsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state postgresqlFlavorsDataSourceModel

	flavors, err := d.client.ListPostgreSQLFlavors(ctx, d.organization.ValueString(), d.project.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read PostgreSQL flavors",
			err.Error(),
		)
		return
	}

	for _, flavor := range flavors {
		flavorState := PostgresqlFlavorModel{
			ID:          types.StringValue(flavor.Id),
			Description: types.StringValue(flavor.Description),
			Default:     types.BoolValue(flavor.Default),
		}

		state.PostgresqlFlavors = append(state.PostgresqlFlavors, flavorState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
