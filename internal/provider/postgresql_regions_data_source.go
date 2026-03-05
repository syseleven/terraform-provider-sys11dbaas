package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	v2 "github.com/syseleven/sys11dbaas-sdk/database/v2"
)

// postgresqlRegionsDataSourceModel maps the data source schema data.
type postgresqlRegionsDataSourceModel struct {
	PostgresqlRegions []PostgresqlRegionModel `tfsdk:"regions"`
}

// PostgresqlRegionModel maps postgresqlRegion schema data.
type PostgresqlRegionModel struct {
	ID          types.String `tfsdk:"id"`
	Description types.String `tfsdk:"description"`
	Default     types.Bool   `tfsdk:"default"`
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &postgresqlRegionsDataSource{}
	_ datasource.DataSourceWithConfigure = &postgresqlRegionsDataSource{}
)

// NewPostgresqlRegionsDataSource is a helper function to simplify the provider implementation.
func NewPostgresqlRegionsDataSource() datasource.DataSource {
	return &postgresqlRegionsDataSource{}
}

// coffeesDataSource is the data source implementation.
type postgresqlRegionsDataSource struct {
	client       *v2.TypedClient
	project      types.String
	organization types.String
}

// Configure adds the provider configured client to the data source.
func (d *postgresqlRegionsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *postgresqlRegionsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_postgresql_regions"
}

// Schema defines the schema for the data source.
func (d *postgresqlRegionsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the list of available regions.",
		Attributes: map[string]schema.Attribute{
			"regions": schema.ListNestedAttribute{
				Description: "List of regions.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Textual identifier of the region.",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "Description of the region.",
							Computed:    true,
						},
						"default": schema.BoolAttribute{
							Description: "In case this region is the default, this field is true.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *postgresqlRegionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state postgresqlRegionsDataSourceModel

	regions, err := d.client.ListPostgreSQLRegions(ctx, d.organization.ValueString(), d.project.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read PostgreSQL regions",
			err.Error(),
		)
		return
	}

	for _, region := range regions {
		regionstate := PostgresqlRegionModel{
			ID:          types.StringValue(region.Id),
			Description: types.StringValue(region.Description),
			Default:     types.BoolValue(region.Default),
		}

		state.PostgresqlRegions = append(state.PostgresqlRegions, regionstate)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
