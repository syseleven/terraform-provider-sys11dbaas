package provider

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	v2 "github.com/syseleven/sys11dbaas-sdk/database/v2"
)

// featuresDataSourceModel maps the data source schema data.
type featuresDataSourceModel struct {
	Type     string         `tfsdk:"type"`
	Features []FeatureModel `tfsdk:"features"`
}

// FeatureModel maps Feature schema data.
type FeatureModel struct {
	Name    types.String `tfsdk:"name"`
	Default types.String `tfsdk:"default"`
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &FeaturesDataSource{}
	_ datasource.DataSourceWithConfigure = &FeaturesDataSource{}
)

// NewFeaturesDataSource is a helper function to simplify the provider implementation.
func NewFeaturesDataSource() datasource.DataSource {
	return &FeaturesDataSource{}
}

// FeaturesDataSource is the data source implementation.
type FeaturesDataSource struct {
	client       *v2.TypedClient
	project      types.String
	organization types.String
}

// Configure adds the provider configured client to the data source.
func (d *FeaturesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *FeaturesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_features"
}

// Schema defines the schema for the data source.
func (d *FeaturesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches the list of available features.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "Type of database to retrieve features for.",
				Required:    true,
			},
			"features": schema.ListNestedAttribute{
				Description: "List of features.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Name of the feature.",
							Computed:    true,
						},
						"default": schema.StringAttribute{
							Description: "Default state of the feature.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *FeaturesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state featuresDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if state.Type == "postgresql" {
		features, err := d.client.ListFeatures(ctx, d.organization.ValueString(), d.project.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Read features",
				err.Error(),
			)
			return
		}

		for _, feature := range slices.SortedFunc(slices.Values(features), func(a, b v2.Feature) int {
			return strings.Compare(a.Id, b.Id)
		}) {
			featureState := FeatureModel{
				Name:    types.StringValue(feature.Id),
				Default: types.StringValue(string(feature.Default)),
			}

			state.Features = append(state.Features, featureState)
		}
	} else {
		resp.Diagnostics.Append(diag.NewAttributeErrorDiagnostic(path.Root("type"), "invalid data source type", "unsupported type selected for retrieving features"))
		return
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
