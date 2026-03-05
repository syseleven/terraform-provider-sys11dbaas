package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	v2 "github.com/syseleven/sys11dbaas-sdk/database/v2"
)

// postgresqlVersionsDataSourceModel maps the data source schema data.
type postgresqlVersionsDataSourceModel struct {
	PostgresqlVersions []PostgresqlVersionModel `tfsdk:"versions"`
}

// PostgresqlVersionModel maps postgresqlVersion schema data.
type PostgresqlVersionModel struct {
	ID          types.String `tfsdk:"id"`
	Description types.String `tfsdk:"description"`
	Default     types.Bool   `tfsdk:"default"`
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &postgresqlVersionsDataSource{}
	_ datasource.DataSourceWithConfigure = &postgresqlVersionsDataSource{}
)

// NewPostgresqlVersionsDataSource is a helper function to simplify the provider implementation.
func NewPostgresqlVersionsDataSource() datasource.DataSource {
	return &postgresqlVersionsDataSource{}
}

// coffeesDataSource is the data source implementation.
type postgresqlVersionsDataSource struct {
	client       *v2.TypedClient
	project      types.String
	organization types.String
}

// Configure adds the provider configured client to the data source.
func (d *postgresqlVersionsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *postgresqlVersionsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_postgresql_versions"
}

// Schema defines the schema for the data source.
func (d *postgresqlVersionsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the list of available postgresql versions.",
		Attributes: map[string]schema.Attribute{
			"versions": schema.ListNestedAttribute{
				Description: "List of versions.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Textual identifier of the version.",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "Textual identifier of the version .",
							Computed:    true,
						},
						"default": schema.BoolAttribute{
							Description: "In case this version is the default, this field is true.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *postgresqlVersionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state postgresqlVersionsDataSourceModel

	versions, err := d.client.ListPostgreSQLVersions(ctx, d.organization.ValueString(), d.project.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read PostgreSQL versions",
			err.Error(),
		)
		return
	}

	for _, Version := range versions {
		versionstate := PostgresqlVersionModel{
			ID:          types.StringValue(Version.Id),
			Description: types.StringValue(Version.Description),
			Default:     types.BoolValue(Version.Default),
		}

		state.PostgresqlVersions = append(state.PostgresqlVersions, versionstate)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
