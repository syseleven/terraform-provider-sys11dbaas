package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	sys11dbaassdk "github.com/syseleven/sys11dbaas-sdk"
)

type DatabaseModelV2 struct {
	ApplicationConfig types.Object `tfsdk:"application_config"`
	CreatedAt         types.String `tfsdk:"created_at"`
	CreatedBy         types.String `tfsdk:"created_by"`
	Description       types.String `tfsdk:"description"`
	LastModifiedAt    types.String `tfsdk:"last_modified_at"`
	LastModifiedBy    types.String `tfsdk:"last_modified_by"`
	Name              types.String `tfsdk:"name"`
	ServiceConfig     types.Object `tfsdk:"service_config"`
	Status            types.String `tfsdk:"status"`
	Phase             types.String `tfsdk:"phase"`
	ResourceStatus    types.String `tfsdk:"resource_status"`
	Uuid              types.String `tfsdk:"uuid"`
}

// resource

type DatabaseResourceV2 struct {
	client          *sys11dbaassdk.Client
	project         types.String
	organization    types.String
	waitForCreation types.Bool
}

func NewDatabaseResourceV2() resource.Resource {
	return &DatabaseResourceV2{}
}

// Metadata returns the resource type name.
func (r *DatabaseResourceV2) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database_v2"
}

// Configure adds the provider configured client to the resource.
func (r *DatabaseResourceV2) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = providerData.client
	r.organization = providerData.organization
	r.project = providerData.project
	r.waitForCreation = providerData.waitForCreation
}

// Read resource information.
func (r *DatabaseResourceV2) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state DatabaseModelV2
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	psqlRequest := &sys11dbaassdk.GetPostgreSQLRequestV2{
		UUID:         state.Uuid.ValueString(),
		Organization: r.organization.ValueString(),
		Project:      r.project.ValueString(),
	}

	var psqlDB *sys11dbaassdk.GetPostgreSQLResponseV2
	var err error
	errCount := 0
	for {
		psqlDB, err = r.client.GetPostgreSQLDBV2(psqlRequest)
		if err != nil {
			errCount++
			if errCount >= 3 {
				resp.Diagnostics.AddError(
					"Unable to Read database",
					err.Error(),
				)
				return
			}
			tflog.Warn(ctx, "Error reading updated database, retrying", map[string]interface{}{"error": err.Error()})
			continue
		} else {
			errCount = 0
			break
		}
	}

	diags = psqlGetResponseToModelV2(ctx, psqlDB, &state, state)
	ctx = tflog.SetField(ctx, "read_target_state", state)
	tflog.Debug(ctx, "Reading database", nil)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create resource.
func (r *DatabaseResourceV2) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan DatabaseModelV2
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceConfig, diags := NewServiceConfigValueV2(plan.ServiceConfig.AttributeTypes(ctx), plan.ServiceConfig.Attributes())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var maintenanceWindow *sys11dbaassdk.MaintenanceWindowV2
	if !serviceConfig.MaintenanceWindow.IsUnknown() {
		maintenanceWindowObj, diags := NewMaintenanceWindowValueV2(serviceConfig.MaintenanceWindow.AttributeTypes(ctx), serviceConfig.MaintenanceWindow.Attributes())
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		maintenanceWindow, diags = maintenanceWindowObj.ToDBaaSSdkObject(ctx)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	applicationConfig, diags := NewApplicationConfigValueV2(plan.ApplicationConfig.AttributeTypes(ctx), plan.ApplicationConfig.Attributes())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var backupSchedule *sys11dbaassdk.PSQLScheduledBackupsV2
	if !applicationConfig.ScheduledBackups.IsUnknown() {
		scheduledBackupsObj, diags := NewScheduledBackupsValueV2(applicationConfig.ScheduledBackups.AttributeTypes(ctx), applicationConfig.ScheduledBackups.Attributes())
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		backupSchedule, diags = scheduledBackupsObj.ToDBaaSSdkObject(ctx)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	var privateNetworkConfig *sys11dbaassdk.PSQLPrivateNetworkConfigRequestV2
	if !applicationConfig.PrivateNetworkConfig.IsUnknown() {
		privateNetworkConfigObj, diags := NewPrivateNetworkConfigValueV2(applicationConfig.PrivateNetworkConfig.AttributeTypes(ctx), applicationConfig.PrivateNetworkConfig.Attributes())
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		privateNetworkConfig, diags = privateNetworkConfigObj.ToDBaaSSdkRequest(ctx)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	var publicNetworkConfig *sys11dbaassdk.PSQLPublicNetworkConfigRequestV2
	if !applicationConfig.PublicNetworkConfig.IsUnknown() {
		publicNetworkConfigObj, diags := NewPublicNetworkConfigValueV2(applicationConfig.PublicNetworkConfig.AttributeTypes(ctx), applicationConfig.PublicNetworkConfig.Attributes())
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		publicNetworkConfig, diags = publicNetworkConfigObj.ToDBaaSSdkRequest(ctx)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	createRequest := &sys11dbaassdk.CreatePostgreSQLRequestV2{
		Organization: r.organization.ValueString(),
		Project:      r.project.ValueString(),
		Name:         plan.Name.ValueString(),
		Description:  plan.Description.ValueString(),
		ServiceConfig: &sys11dbaassdk.PSQLServiceConfigRequestV2{
			Disksize:          sys11dbaassdk.Int64ToIntPtr(serviceConfig.Disksize.ValueInt64()),
			Type:              serviceConfig.ServiceConfigType.ValueString(),
			Flavor:            serviceConfig.Flavor.ValueString(),
			Region:            serviceConfig.Region.ValueString(),
			MaintenanceWindow: maintenanceWindow,
		},
		ApplicationConfig: &sys11dbaassdk.PSQLApplicationConfigRequestV2{
			Type:                 applicationConfig.ApplicationConfigType.ValueString(),
			Password:             applicationConfig.Password.ValueString(),
			Instances:            sys11dbaassdk.IntPtr(int(applicationConfig.Instances.ValueInt64())),
			Version:              applicationConfig.Version.ValueString(),
			ScheduledBackups:     backupSchedule,
			PrivateNetworkConfig: privateNetworkConfig,
			PublicNetworkConfig:  publicNetworkConfig,
		},
	}

	d, _ := json.Marshal(createRequest)
	tflog.Debug(ctx, string(d), nil)

	// Create new db
	createResponse, err := r.client.CreatePostgreSQLDBV2(createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating database",
			"Could not create database, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	targetState := DatabaseModelV2{}
	diags = psqlCreateResponseToModelV2(ctx, createResponse, plan, &targetState)

	ctx = tflog.SetField(ctx, "create_target_state", &targetState)
	tflog.Debug(ctx, "[CREATE] Created", nil)
	resp.Diagnostics.Append(diags...)

	tflog.Debug(ctx, "[CREATE] Wait for creation: "+r.waitForCreation.String(), nil)
	if r.waitForCreation.ValueBool() {
		getRequest := &sys11dbaassdk.GetPostgreSQLRequestV2{
			Organization: r.organization.ValueString(),
			Project:      r.project.ValueString(),
			UUID:         createResponse.UUID,
		}
		sleepFor := time.Duration(30 * time.Second)
		errCount := 0
		for retryCount := 0; targetState.Status.ValueString() != sys11dbaassdk.STATE_READY; {
			if retryCount == int((CREATE_RETRY_LIMIT / sleepFor).Abs()) {
				diags = resp.State.Set(ctx, &targetState)
				resp.Diagnostics.Append(diags...)
				resp.Diagnostics.AddError("RetryLimit reached during wait_for_creation", "The retry limit of "+CREATE_RETRY_LIMIT.String()+" was reached while waiting for creation of database")
				return
			}
			time.Sleep(sleepFor)
			getResponse, err := r.client.GetPostgreSQLDBV2(getRequest)
			if err != nil {
				errCount++
				if errCount >= 3 {
					resp.Diagnostics.AddError(
						"Error read database during wait",
						"Could not read database during wait, unexpected error: "+err.Error(),
					)
					return
				}
				tflog.Warn(ctx, "Error reading updated database, retrying", map[string]interface{}{"error": err.Error()})
				continue
			}
			errCount = 0
			diags = psqlGetResponseToModelV2(ctx, getResponse, &targetState, plan)

			ctx = tflog.SetField(ctx, "create_target_state", &targetState)
			tflog.Debug(ctx, "[CREATE] Current creation state", nil)
			resp.Diagnostics.Append(diags...)
			retryCount++
		}
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, &targetState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete resource.
func (r *DatabaseResourceV2) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Get current state
	var state DatabaseModelV2
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	psqlRequest := &sys11dbaassdk.DeletePostgreSQLRequestV2{
		UUID:         state.Uuid.ValueString(),
		Organization: r.organization.ValueString(),
		Project:      r.project.ValueString(),
	}

	_, err := r.client.DeletePostgreSQLDBV2(psqlRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Delete Database",
			err.Error(),
		)
		return
	}

	// Set refreshed state
	resp.State.RemoveResource(ctx)
}

// Update resource.
func (r *DatabaseResourceV2) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan
	var plan DatabaseModelV2
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Get current state
	var state DatabaseModelV2
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceConfig, diags := NewServiceConfigValueV2(plan.ServiceConfig.AttributeTypes(ctx), plan.ServiceConfig.Attributes())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var maintenanceWindow *sys11dbaassdk.MaintenanceWindowV2
	if !serviceConfig.MaintenanceWindow.IsUnknown() {
		maintenanceWindowObj, diags := NewMaintenanceWindowValueV2(serviceConfig.MaintenanceWindow.AttributeTypes(ctx), serviceConfig.MaintenanceWindow.Attributes())
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		maintenanceWindow, diags = maintenanceWindowObj.ToDBaaSSdkObject(ctx)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	applicationConfig, diags := NewApplicationConfigValueV2(plan.ApplicationConfig.AttributeTypes(ctx), plan.ApplicationConfig.Attributes())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var backupSchedule *sys11dbaassdk.PSQLScheduledBackupsV2
	if !applicationConfig.ScheduledBackups.IsUnknown() {
		scheduledBackupsObj, diags := NewScheduledBackupsValueV2(applicationConfig.ScheduledBackups.AttributeTypes(ctx), applicationConfig.ScheduledBackups.Attributes())
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		backupSchedule, diags = scheduledBackupsObj.ToDBaaSSdkObject(ctx)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	var privateNetworkConfig *sys11dbaassdk.PSQLPrivateNetworkConfigRequestV2
	if !applicationConfig.PrivateNetworkConfig.IsUnknown() {
		privateNetworkConfigObj, diags := NewPrivateNetworkConfigValueV2(applicationConfig.PrivateNetworkConfig.AttributeTypes(ctx), applicationConfig.PrivateNetworkConfig.Attributes())
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		privateNetworkConfig, diags = privateNetworkConfigObj.ToDBaaSSdkRequest(ctx)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	var publicNetworkConfig *sys11dbaassdk.PSQLPublicNetworkConfigRequestV2
	if !applicationConfig.PublicNetworkConfig.IsUnknown() {
		publicNetworkConfigObj, diags := NewPublicNetworkConfigValueV2(applicationConfig.PublicNetworkConfig.AttributeTypes(ctx), applicationConfig.PublicNetworkConfig.Attributes())
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		publicNetworkConfig, diags = publicNetworkConfigObj.ToDBaaSSdkRequest(ctx)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	updateRequest := &sys11dbaassdk.UpdatePostgreSQLRequestV2{
		UUID:         state.Uuid.ValueString(),
		Organization: r.organization.ValueString(),
		Project:      r.project.ValueString(),
		Name:         plan.Name.ValueString(),
		Description:  plan.Description.ValueString(),
		ServiceConfig: &sys11dbaassdk.PSQLServiceConfigUpdateRequestV2{
			Disksize:          sys11dbaassdk.Int64ToIntPtr(serviceConfig.Disksize.ValueInt64()),
			Type:              serviceConfig.ServiceConfigType.ValueString(),
			Flavor:            serviceConfig.Flavor.ValueString(),
			MaintenanceWindow: maintenanceWindow,
		},
		ApplicationConfig: &sys11dbaassdk.PSQLApplicationConfigUpdateRequestV2{
			Password:             applicationConfig.Password.ValueString(),
			Instances:            sys11dbaassdk.IntPtr(int(applicationConfig.Instances.ValueInt64())),
			Version:              applicationConfig.Version.ValueString(),
			ScheduledBackups:     backupSchedule,
			PrivateNetworkConfig: privateNetworkConfig,
			PublicNetworkConfig:  publicNetworkConfig,
		},
	}

	d, _ := json.Marshal(updateRequest)
	tflog.Debug(ctx, string(d), nil)

	// Update psql
	_, err := r.client.UpdatePostgreSQLDBV2(updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating database",
			"Could not update database, unexpected error: "+err.Error(),
		)
		resp.Diagnostics.Append(diags...)
		return
	}

	errCount := 0
	for i := 0; i < 10; i++ {
		time.Sleep(2 * time.Second) // give DBaaS time to propagate changes
		tflog.Debug(ctx, "wait 2 seconds to give DBaaS time to propagate", nil)

		// Update psql
		getRequest := &sys11dbaassdk.GetPostgreSQLRequestV2{
			Organization: r.organization.ValueString(),
			Project:      r.project.ValueString(),
			UUID:         state.Uuid.ValueString(),
		}
		getResponse, err := r.client.GetPostgreSQLDBV2(getRequest)
		if err != nil {
			errCount++
			if errCount >= 3 {
				resp.Diagnostics.AddError(
					"Error reading updated database",
					"Could not read updated database, unexpected error: "+err.Error(),
				)
				return
			}
			tflog.Warn(ctx, "Error reading updated database, retrying", map[string]interface{}{"error": err.Error()})
			continue
		}
		errCount = 0

		// Map response body to schema and populate Computed attribute values
		targetState := DatabaseModelV2{}
		diags = psqlGetResponseToModelV2(ctx, getResponse, &targetState, plan)

		ctx = tflog.SetField(ctx, "update_target_state", &targetState)
		tflog.Debug(ctx, "Updated State", nil)
		resp.Diagnostics.Append(diags...)

		// Set state to fully populated data
		diags = resp.State.Set(ctx, &targetState)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			continue
		}
	}
}

// Schema defines the schema for the resource.
func (r *DatabaseResourceV2) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = DatabaseResourceV2Schema(ctx)
}

func psqlCreateResponseToModelV2(ctx context.Context, db *sys11dbaassdk.CreatePostgreSQLResponseV2, plan DatabaseModelV2, targetState *DatabaseModelV2) diag.Diagnostics {

	var diags diag.Diagnostics

	ctx = tflog.SetField(ctx, "conversion_create_source_response", db)
	tflog.Debug(ctx, "Converting create api response")

	var recoveryObjValue basetypes.ObjectValue
	if db.ApplicationConfig.Recovery != nil {
		recoveryObjValue, _ = RecoveryValueV2{
			Exclusive:  types.BoolValue(db.ApplicationConfig.Recovery.Exclusive),
			Source:     types.StringValue(db.ApplicationConfig.Recovery.Source),
			TargetLsn:  types.StringValue(db.ApplicationConfig.Recovery.TargetLSN),
			TargetName: types.StringValue(db.ApplicationConfig.Recovery.TargetName),
			TargetXid:  types.StringValue(db.ApplicationConfig.Recovery.TargetXID),
			TargetTime: types.StringValue(db.ApplicationConfig.Recovery.TargetTime.Format(time.RFC3339)),
		}.ToObjectValue(ctx)
	}

	scheduleObjVal, _ := ScheduleValueV2{
		Hour:   types.Int64Value(int64(*db.ApplicationConfig.ScheduledBackups.Schedule.Hour)),
		Minute: types.Int64Value(int64(*db.ApplicationConfig.ScheduledBackups.Schedule.Minute)),
	}.ToObjectValue(ctx)

	scheduledBackupsObjVal, _ := ScheduledBackupsValueV2{
		Schedule:  scheduleObjVal,
		Retention: types.Int64Value(int64(*db.ApplicationConfig.ScheduledBackups.Retention)),
	}.ToObjectValue(ctx)

	maintenanceWindowObjVal, _ := MaintenanceWindowValueV2{
		DayOfWeek:   types.Int64Value(int64(*db.ServiceConfig.MaintenanceWindow.DayOfWeek)),
		StartHour:   types.Int64Value(int64(*db.ServiceConfig.MaintenanceWindow.StartHour)),
		StartMinute: types.Int64Value(int64(*db.ServiceConfig.MaintenanceWindow.StartMinute)),
	}.ToObjectValue(ctx)

	privateAllowedCIRDs, d := types.ListValueFrom(ctx, types.StringType, (*db).ApplicationConfig.PrivateNetworkConfig.AllowedCIDRs)
	diags.Append(d...)

	privateNetworkConfigObjVal, _ := PrivateNetworkConfigValueV2{
		Enabled:          types.BoolValue(db.ApplicationConfig.PrivateNetworkConfig.Enabled),
		Hostname:         types.StringValue(db.ApplicationConfig.PrivateNetworkConfig.Hostname),
		IPAddress:        types.StringValue(db.ApplicationConfig.PrivateNetworkConfig.IPAddress),
		AllowedCIDRs:     privateAllowedCIRDs,
		SharedSubnetCIDR: types.StringValue(*db.ApplicationConfig.PrivateNetworkConfig.SharedSubnetCIDR),
		SharedSubnetID:   types.StringValue(db.ApplicationConfig.PrivateNetworkConfig.SharedSubnetID),
		SharedNetworkID:  types.StringValue(db.ApplicationConfig.PrivateNetworkConfig.SharedNetworkID),
	}.ToObjectValue(ctx)

	publicAllowedCIRDs, d := types.ListValueFrom(ctx, types.StringType, (*db).ApplicationConfig.PublicNetworkConfig.AllowedCIDRs)
	diags.Append(d...)

	publicNetworkConfigObjVal, _ := PublicNetworkConfigValueV2{
		Enabled:      types.BoolValue(db.ApplicationConfig.PublicNetworkConfig.Enabled),
		Hostname:     types.StringValue(db.ApplicationConfig.PublicNetworkConfig.Hostname),
		IPAddress:    types.StringValue(db.ApplicationConfig.PublicNetworkConfig.IPAddress),
		AllowedCIDRs: publicAllowedCIRDs,
	}.ToObjectValue(ctx)

	var targetServiceConfig ServiceConfigValueV2
	targetServiceConfig.Disksize = types.Int64Value(int64(*db.ServiceConfig.Disksize))
	targetServiceConfig.ServiceConfigType = types.StringValue(db.ServiceConfig.Type)
	targetServiceConfig.Flavor = types.StringValue(db.ServiceConfig.Flavor)
	targetServiceConfig.Region = types.StringValue(db.ServiceConfig.Region)
	targetServiceConfig.MaintenanceWindow = maintenanceWindowObjVal

	targetServiceConfigObj, diags := targetServiceConfig.ToObjectValue(ctx)

	var targetApplicationConfig ApplicationConfigValueV2
	targetApplicationConfig.ApplicationConfigType = types.StringValue(db.ApplicationConfig.Type)
	targetApplicationConfig.Password = types.StringValue(strings.Trim(plan.ApplicationConfig.Attributes()["password"].String(), "\"")) // take this from the plan, since it is not included in the response
	targetApplicationConfig.Instances = types.Int64Value(int64(*db.ApplicationConfig.Instances))
	targetApplicationConfig.Version = types.StringValue(db.ApplicationConfig.Version)
	targetApplicationConfig.ScheduledBackups = scheduledBackupsObjVal
	targetApplicationConfig.Recovery = recoveryObjValue
	targetApplicationConfig.PrivateNetworkConfig = privateNetworkConfigObjVal
	targetApplicationConfig.PublicNetworkConfig = publicNetworkConfigObjVal

	targetApplicationConfigObj, diags := targetApplicationConfig.ToObjectValue(ctx)

	targetState.Uuid = types.StringValue(db.UUID)
	targetState.Name = types.StringValue(db.Name)
	targetState.Description = types.StringValue(db.Description)
	targetState.Status = types.StringValue(db.Status)
	targetState.Phase = types.StringValue(db.Phase)
	targetState.ResourceStatus = types.StringValue(db.ResourceStatus)
	targetState.CreatedBy = types.StringValue(db.CreatedBy)
	targetState.CreatedAt = types.StringValue(db.CreatedAt)
	targetState.LastModifiedBy = types.StringValue(db.LastModifiedBy)
	targetState.LastModifiedAt = types.StringValue(db.LastModifiedAt)
	targetState.ApplicationConfig = targetApplicationConfigObj
	targetState.ServiceConfig = targetServiceConfigObj

	ctx = tflog.SetField(ctx, "conversion_create_target_state", &targetState)
	tflog.Debug(ctx, "Converted api create response to state")

	return diags

}

func psqlGetResponseToModelV2(ctx context.Context, db *sys11dbaassdk.GetPostgreSQLResponseV2, targetState *DatabaseModelV2, previousState DatabaseModelV2) diag.Diagnostics {

	var diags diag.Diagnostics

	ctx = tflog.SetField(ctx, "conversion_read_source_response", db)
	tflog.Debug(ctx, "Converting read api response")

	var recoveryObjValue basetypes.ObjectValue
	if db.ApplicationConfig.Recovery != nil {
		recoveryObjValue, _ = RecoveryValueV2{
			Exclusive:  types.BoolValue(db.ApplicationConfig.Recovery.Exclusive),
			Source:     types.StringValue(db.ApplicationConfig.Recovery.Source),
			TargetLsn:  types.StringValue(db.ApplicationConfig.Recovery.TargetLSN),
			TargetName: types.StringValue(db.ApplicationConfig.Recovery.TargetName),
			TargetXid:  types.StringValue(db.ApplicationConfig.Recovery.TargetXID),
			TargetTime: types.StringValue(db.ApplicationConfig.Recovery.TargetTime.Format(time.RFC3339)),
		}.ToObjectValue(ctx)
	}

	scheduleObjVal, _ := ScheduleValueV2{
		Hour:   types.Int64Value(int64(*db.ApplicationConfig.ScheduledBackups.Schedule.Hour)),
		Minute: types.Int64Value(int64(*db.ApplicationConfig.ScheduledBackups.Schedule.Minute)),
	}.ToObjectValue(ctx)

	scheduledBackupsObjVal, _ := ScheduledBackupsValueV2{
		Schedule:  scheduleObjVal,
		Retention: types.Int64Value(int64(*db.ApplicationConfig.ScheduledBackups.Retention)),
	}.ToObjectValue(ctx)

	maintenanceWindowObjVal, _ := MaintenanceWindowValueV2{
		DayOfWeek:   types.Int64Value(int64(*db.ServiceConfig.MaintenanceWindow.DayOfWeek)),
		StartHour:   types.Int64Value(int64(*db.ServiceConfig.MaintenanceWindow.StartHour)),
		StartMinute: types.Int64Value(int64(*db.ServiceConfig.MaintenanceWindow.StartMinute)),
	}.ToObjectValue(ctx)

	privateAllowedCIRDs, d := types.ListValueFrom(ctx, types.StringType, (*db).ApplicationConfig.PrivateNetworkConfig.AllowedCIDRs)
	diags.Append(d...)

	privateNetworkConfigObjVal, _ := PrivateNetworkConfigValueV2{
		Enabled:          types.BoolValue(db.ApplicationConfig.PrivateNetworkConfig.Enabled),
		Hostname:         types.StringValue(db.ApplicationConfig.PrivateNetworkConfig.Hostname),
		IPAddress:        types.StringValue(db.ApplicationConfig.PrivateNetworkConfig.IPAddress),
		AllowedCIDRs:     privateAllowedCIRDs,
		SharedSubnetCIDR: types.StringValue(*db.ApplicationConfig.PrivateNetworkConfig.SharedSubnetCIDR),
		SharedSubnetID:   types.StringValue(db.ApplicationConfig.PrivateNetworkConfig.SharedSubnetID),
		SharedNetworkID:  types.StringValue(db.ApplicationConfig.PrivateNetworkConfig.SharedNetworkID),
	}.ToObjectValue(ctx)

	publicAllowedCIRDs, d := types.ListValueFrom(ctx, types.StringType, (*db).ApplicationConfig.PublicNetworkConfig.AllowedCIDRs)
	diags.Append(d...)

	publicNetworkConfigObjVal, _ := PublicNetworkConfigValueV2{
		Enabled:      types.BoolValue(db.ApplicationConfig.PublicNetworkConfig.Enabled),
		Hostname:     types.StringValue(db.ApplicationConfig.PublicNetworkConfig.Hostname),
		IPAddress:    types.StringValue(db.ApplicationConfig.PublicNetworkConfig.IPAddress),
		AllowedCIDRs: publicAllowedCIRDs,
	}.ToObjectValue(ctx)

	var targetServiceConfig ServiceConfigValueV2
	targetServiceConfig.Disksize = types.Int64Value(int64(*db.ServiceConfig.Disksize))
	targetServiceConfig.ServiceConfigType = types.StringValue(db.ServiceConfig.Type)
	targetServiceConfig.Flavor = types.StringValue(db.ServiceConfig.Flavor)
	targetServiceConfig.Region = types.StringValue(db.ServiceConfig.Region)
	targetServiceConfig.MaintenanceWindow = maintenanceWindowObjVal

	targetServiceConfigObj, diags := targetServiceConfig.ToObjectValue(ctx)

	previousPassword := strings.Trim(previousState.ApplicationConfig.Attributes()["password"].String(), "\"")

	var targetApplicationConfig ApplicationConfigValueV2
	targetApplicationConfig.ApplicationConfigType = types.StringValue(db.ApplicationConfig.Type)
	targetApplicationConfig.Password = types.StringValue(previousPassword)
	targetApplicationConfig.Instances = types.Int64Value(int64(*db.ApplicationConfig.Instances))
	targetApplicationConfig.Version = types.StringValue(db.ApplicationConfig.Version)
	targetApplicationConfig.ScheduledBackups = scheduledBackupsObjVal
	targetApplicationConfig.Recovery = recoveryObjValue
	targetApplicationConfig.PrivateNetworkConfig = privateNetworkConfigObjVal
	targetApplicationConfig.PublicNetworkConfig = publicNetworkConfigObjVal

	targetApplicationConfigObj, diags := targetApplicationConfig.ToObjectValue(ctx)

	targetState.Uuid = types.StringValue(db.UUID)
	targetState.Name = types.StringValue(db.Name)
	targetState.Description = types.StringValue(db.Description)
	targetState.Status = types.StringValue(db.Status)
	targetState.Phase = types.StringValue(db.Phase)
	targetState.ResourceStatus = types.StringValue(db.ResourceStatus)
	targetState.CreatedBy = types.StringValue(db.CreatedBy)
	targetState.CreatedAt = types.StringValue(db.CreatedAt)
	targetState.LastModifiedBy = types.StringValue(db.LastModifiedBy)
	targetState.LastModifiedAt = types.StringValue(db.LastModifiedAt)
	targetState.ApplicationConfig = targetApplicationConfigObj
	targetState.ServiceConfig = targetServiceConfigObj

	ctx = tflog.SetField(ctx, "conversion_read_target_state", targetState)
	tflog.Debug(ctx, "Converted api read response to state")

	return diags

}
