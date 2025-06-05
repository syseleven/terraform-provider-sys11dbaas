// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"
	"strconv"

	sys11dbaassdk "github.com/syseleven/sys11dbaas-sdk"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure Sys11DBaaSProvider satisfies various provider interfaces.
var _ provider.Provider = &Sys11DBaaSProvider{}

// Sys11DBaaSProvider defines the provider implementation.
type Sys11DBaaSProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Sys11DBaaSProvider maps provider schema data to a Go type.
type Sys11DBaaSProviderModel struct {
	URL             types.String `tfsdk:"url"`
	ApiKey          types.String `tfsdk:"api_key"`
	Project         types.String `tfsdk:"project"`
	Organization    types.String `tfsdk:"organization"`
	WaitForCreation types.Bool   `tfsdk:"wait_for_creation"`
}

type sys11DBaaSProviderData struct {
	client          *sys11dbaassdk.Client
	project         types.String `tfsdk:"project"`
	organization    types.String `tfsdk:"organization"`
	waitForCreation types.Bool   `tfsdk:"wait_for_creation"`
}

func (p *Sys11DBaaSProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "sys11dbaas"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *Sys11DBaaSProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"url": schema.StringAttribute{
				Required:    false,
				Optional:    true,
				Description: "URL of the DBaaS API. If omitted, the `SYS11DBAAS_URL` environment variable is used. Otherwise fallbacks to https://dbaas.apis.syseleven.de",
			},
			"api_key": schema.StringAttribute{
				Required:    false,
				Optional:    true,
				Description: "API key to use for authentication to the DBaaS API. If omitted, the `SYS11DBAAS_API_KEY` environment variable is used.",
			},
			"organization": schema.StringAttribute{
				Required:    false,
				Optional:    true,
				Description: "ID of your organization. If omitted, the `SYS11DBAAS_ORGANIZATION` environment variable is used.",
			},
			"project": schema.StringAttribute{
				Required:    false,
				Optional:    true,
				Description: "ID of your project. If omitted, the `SYS11DBAAS_PROJECT` environment variable is used.",
			},
			"wait_for_creation": schema.BoolAttribute{
				Required:    false,
				Optional:    true,
				Description: "Whether to wait for the service to be created. If omitted, the `SYS11DBAAS_WAIT_FOR_CREATION` environment variable is used. Defaults to true",
			},
		},
	}
}

func (p *Sys11DBaaSProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Sys11DBaaS client")

	// Retrieve provider data from configuration
	var config Sys11DBaaSProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.URL.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("url"),
			"Unknown Sys11DBaaS API Url",
			"The provider cannot create the Sys11DBaaS API client as there is an unknown configuration value for the Sys11DBaaS API url. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the SYS11DBAAS_URL environment variable.",
		)
	}

	if config.ApiKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Unknown Sys11DBaaS API ApiKey",
			"The provider cannot create the Sys11DBaaS API client as there is an unknown configuration value for the Sys11DBaaS API username. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the SYS11DBAAS_USERNAME environment variable.",
		)
	}

	if config.Organization.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("organization"),
			"Unknown Sys11DBaaS API org",
			"The provider cannot create the Sys11DBaaS API client as there is an unknown configuration value for the Sys11DBaaS API username. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the SYS11DBAAS_USERNAME environment variable.",
		)
	}

	if config.Project.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("project"),
			"Unknown Sys11DBaaS API project",
			"The provider cannot create the Sys11DBaaS API client as there is an unknown configuration value for the Sys11DBaaS API username. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the SYS11DBAAS_USERNAME environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.
	var url string
	var set bool
	if url, set = os.LookupEnv("SYS11DBAAS_URL"); !set {
		url = "https://dbaas.apis.syseleven.de"
	}

	apikey := os.Getenv("SYS11DBAAS_API_KEY")
	organization := os.Getenv("SYS11DBAAS_ORGANIZATION")
	project := os.Getenv("SYS11DBAAS_PROJECT")

	var waitForCreation bool
	if waitForCreationEnv, set := os.LookupEnv("SYS11DBAAS_WAIT_FOR_CREATION"); !set {
		waitForCreation = true
	} else {
		waitForCreation, _ = strconv.ParseBool(waitForCreationEnv)
	}

	if !config.URL.IsNull() {
		url = config.URL.ValueString()
	}

	if !config.ApiKey.IsNull() {
		apikey = config.ApiKey.ValueString()
	}

	if !config.Organization.IsNull() {
		organization = config.Organization.ValueString()
	}

	if !config.Project.IsNull() {
		project = config.Project.ValueString()
	}

	if !config.WaitForCreation.IsNull() {
		waitForCreation = config.WaitForCreation.ValueBool()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if url == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("url"),
			"Missing Sys11DBaaS API Url",
			"The provider cannot create the Sys11DBaaS API client as there is a missing or empty value for the Sys11DBaaS API url. "+
				"Set the url value in the configuration or use the SYS11DBAAS_URL environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if apikey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Missing Sys11DBaaS API ApiKey",
			"The provider cannot create the Sys11DBaaS API client as there is a missing or empty value for the Sys11DBaaS API ApiKey. "+
				"Set the api_key value in the configuration or use the SYS11DBAAS_API_KEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if organization == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("organization"),
			"Missing Sys11DBaaS organization",
			"The provider cannot create the Sys11DBaaS API client as there is a missing or empty value for the Sys11DBaaS organization. "+
				"Set the organization value in the configuration or use the SYS11DBAAS_ORGANIZATION environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if project == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("project"),
			"Missing Sys11DBaaS project",
			"The provider cannot create the Sys11DBaaS API client as there is a missing or empty value for the Sys11DBaaS project. "+
				"Set the project value in the configuration or use the SYS11DBAAS_PROJECT environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "sys11dbaas_url", url)
	ctx = tflog.SetField(ctx, "sys11dbaas_api_key", apikey)
	ctx = tflog.SetField(ctx, "sys11dbaas_organization", organization)
	ctx = tflog.SetField(ctx, "sys11dbaas_project", project)
	ctx = tflog.SetField(ctx, "sys11dbaas_wait_for_creation", waitForCreation)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "sys11dbaas_api_key")

	tflog.Debug(ctx, "Creating Sys11DBaaS client")

	agent := "sys11dbaas-terraform/" + p.version

	// Create a new Sys11DBaaS client using the configuration values
	client, err := sys11dbaassdk.NewClient(url, apikey, agent, 60, sys11dbaassdk.AuthModeApiKey)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Sys11DBaaS API Client",
			"An unexpected error occurred when creating the Sys11DBaaS API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Sys11DBaaS Client Error: "+err.Error(),
		)
		return
	}

	// Make the Sys11DBaaS client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = &sys11DBaaSProviderData{
		client:          client,
		project:         types.StringValue(project),
		organization:    types.StringValue(organization),
		waitForCreation: types.BoolValue(waitForCreation),
	}
	resp.ResourceData = &sys11DBaaSProviderData{
		client:          client,
		project:         types.StringValue(project),
		organization:    types.StringValue(organization),
		waitForCreation: types.BoolValue(waitForCreation),
	}

	tflog.Info(ctx, "Configured Sys11DBaaS client", map[string]any{"success": true})
}

func (p *Sys11DBaaSProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewDatabaseResource,
		NewDatabaseResourceV2,
	}
}

func (p *Sys11DBaaSProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &Sys11DBaaSProvider{
			version: version,
		}
	}
}
