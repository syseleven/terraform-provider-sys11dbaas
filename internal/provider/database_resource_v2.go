package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	apiv2 "github.com/syseleven/sys11dbaas-sdk/apiv2"
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
	client          *apiv2.TypedClient
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

	r.client = providerData.client.V2()
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

	var psqlDB apiv2.PostgreSQLGetResponse
	var err error
	errCount := 0
	for {
		psqlDB, err = r.client.GetPostgreSQL(ctx, r.organization.ValueString(), r.project.ValueString(), state.Uuid.ValueString())
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

	var maintenanceWindow *apiv2.PostgreSQLMaintenance
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

	var backupSchedule *apiv2.PostgreSQLBackupSchedule
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

	var privateNetworking *apiv2.PostgreSQLPrivateNetworking
	if !applicationConfig.PrivateNetworking.IsUnknown() {
		privateNetworkingObj, diags := NewPrivateNetworkingValueV2(applicationConfig.PrivateNetworking.AttributeTypes(ctx), applicationConfig.PrivateNetworking.Attributes())
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		privateNetworking, diags = privateNetworkingObj.ToDBaaSSdkRequest(ctx)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	var publicNetworking *apiv2.PostgreSQLPublicNetworking
	if !applicationConfig.PublicNetworking.IsUnknown() {
		publicNetworkingObj, diags := NewPublicNetworkingValueV2(applicationConfig.PublicNetworking.AttributeTypes(ctx), applicationConfig.PublicNetworking.Attributes())
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		publicNetworking, diags = publicNetworkingObj.ToDBaaSSdkRequest(ctx)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	createRequest := apiv2.PostgreSQLCreateRequest{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueStringPointer(),
		ServiceConfig: apiv2.PostgreSQLServiceConfig{
			Disksize:          serviceConfig.Disksize.ValueInt64Pointer(),
			Type:              serviceConfig.ServiceConfigType.ValueString(),
			Flavor:            serviceConfig.Flavor.ValueString(),
			Region:            serviceConfig.Region.ValueString(),
			MaintenanceWindow: maintenanceWindow,
		},
		ApplicationConfig: apiv2.PostgreSQLApplicationConfig{
			Type:              applicationConfig.ApplicationConfigType.ValueString(),
			Password:          applicationConfig.Password.ValueString(),
			Instances:         applicationConfig.Instances.ValueInt64Pointer(),
			Version:           applicationConfig.Version.ValueString(),
			ScheduledBackups:  backupSchedule,
			PrivateNetworking: privateNetworking,
			PublicNetworking:  publicNetworking,
		},
	}

	d, _ := json.Marshal(createRequest)
	tflog.Debug(ctx, string(d), nil)

	// Create new db
	createResponse, err := r.client.CreatePostgreSQL(ctx, r.organization.ValueString(), r.project.ValueString(), createRequest)
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
		sleepFor := time.Duration(30 * time.Second)
		errCount := 0
		for retryCount := 0; targetState.Status.ValueString() != apiv2.StateReady; {
			if retryCount == int((CREATE_RETRY_LIMIT / sleepFor).Abs()) {
				diags = resp.State.Set(ctx, &targetState)
				resp.Diagnostics.Append(diags...)
				resp.Diagnostics.AddError("RetryLimit reached during wait_for_creation", "The retry limit of "+CREATE_RETRY_LIMIT.String()+" was reached while waiting for creation of database")
				return
			}
			time.Sleep(sleepFor)
			getResponse, err := r.client.GetPostgreSQL(ctx, r.organization.ValueString(), r.project.ValueString(), createResponse.Uuid)
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

	_, err := r.client.DeletePostgreSQL(ctx, r.organization.ValueString(), r.project.ValueString(), state.Uuid.ValueString())
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

	var maintenanceWindow *apiv2.PostgreSQLMaintenance
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

	var backupSchedule *apiv2.PostgreSQLBackupSchedule
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

	var privateNetworking *apiv2.PostgreSQLPrivateNetworking
	if !applicationConfig.PrivateNetworking.IsUnknown() {
		privateNetworkingObj, diags := NewPrivateNetworkingValueV2(applicationConfig.PrivateNetworking.AttributeTypes(ctx), applicationConfig.PrivateNetworking.Attributes())
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		privateNetworking, diags = privateNetworkingObj.ToDBaaSSdkRequest(ctx)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	var publicNetworking *apiv2.PostgreSQLPublicNetworking
	if !applicationConfig.PublicNetworking.IsUnknown() {
		publicNetworkingObj, diags := NewPublicNetworkingValueV2(applicationConfig.PublicNetworking.AttributeTypes(ctx), applicationConfig.PublicNetworking.Attributes())
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		publicNetworking, diags = publicNetworkingObj.ToDBaaSSdkRequest(ctx)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	updateRequest := apiv2.PostgreSQLCreateRequest{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueStringPointer(),
		ServiceConfig: apiv2.PostgreSQLServiceConfig{
			Disksize:          serviceConfig.Disksize.ValueInt64Pointer(),
			Type:              serviceConfig.ServiceConfigType.ValueString(),
			Flavor:            serviceConfig.Flavor.ValueString(),
			MaintenanceWindow: maintenanceWindow,
		},
		ApplicationConfig: apiv2.PostgreSQLApplicationConfig{
			Password:          applicationConfig.Password.ValueString(),
			Instances:         applicationConfig.Instances.ValueInt64Pointer(),
			Version:           applicationConfig.Version.ValueString(),
			ScheduledBackups:  backupSchedule,
			PrivateNetworking: privateNetworking,
			PublicNetworking:  publicNetworking,
		},
	}

	d, _ := json.Marshal(updateRequest)
	tflog.Debug(ctx, string(d), nil)

	// Update psql
	_, err := r.client.UpdatePostgreSQL(ctx, r.organization.ValueString(), r.project.ValueString(), state.Uuid.ValueString(), updateRequest)
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

		getResponse, err := r.client.GetPostgreSQL(ctx, r.organization.ValueString(), r.project.ValueString(), state.Uuid.ValueString())
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

func (r *DatabaseResourceV2) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}

// Schema defines the schema for the resource.
func (r *DatabaseResourceV2) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = DatabaseResourceV2Schema(ctx)
}

func psqlCreateResponseToModelV2(ctx context.Context, db apiv2.PostgreSQLGetResponse, plan DatabaseModelV2, targetState *DatabaseModelV2) diag.Diagnostics {
	var diags diag.Diagnostics

	ctx = tflog.SetField(ctx, "conversion_create_source_response", db)
	tflog.Debug(ctx, "Converting create api response")

	var recoveryObjValue basetypes.ObjectValue
	if db.ApplicationConfig.Recovery != nil {
		recoveryObjValue, _ = RecoveryValueV2{
			Exclusive:  types.BoolPointerValue(db.ApplicationConfig.Recovery.Exclusive),
			Source:     types.StringPointerValue(db.ApplicationConfig.Recovery.Source),
			TargetLsn:  types.StringPointerValue(db.ApplicationConfig.Recovery.TargetLsn),
			TargetName: types.StringPointerValue(db.ApplicationConfig.Recovery.TargetName),
			TargetXid:  types.StringPointerValue(db.ApplicationConfig.Recovery.TargetXid),
			TargetTime: types.StringValue(db.ApplicationConfig.Recovery.TargetTime.Format(time.RFC3339)),
		}.ToObjectValue(ctx)
	}

	var scheduledBackupsObjVal basetypes.ObjectValue
	if db.ApplicationConfig.ScheduledBackups != nil && db.ApplicationConfig.ScheduledBackups.Schedule != nil {
		scheduleObjVal, _ := ScheduleValueV2{
			Hour:   types.Int64PointerValue(db.ApplicationConfig.ScheduledBackups.Schedule.Hour),
			Minute: types.Int64PointerValue(db.ApplicationConfig.ScheduledBackups.Schedule.Minute),
		}.ToObjectValue(ctx)

		scheduledBackupsObjVal, _ = ScheduledBackupsValueV2{
			Schedule:  scheduleObjVal,
			Retention: types.Int64PointerValue(db.ApplicationConfig.ScheduledBackups.Retention),
		}.ToObjectValue(ctx)
	} else {
		scheduledBackupsObjVal = types.ObjectNull(ScheduledBackupsValueV2{}.AttributeTypes(ctx))
	}

	var maintenanceWindowObjVal basetypes.ObjectValue
	if db.ServiceConfig.MaintenanceWindow != nil {
		maintenanceWindowObjVal, _ = MaintenanceWindowValueV2{
			DayOfWeek:   types.Int64PointerValue(db.ServiceConfig.MaintenanceWindow.DayOfWeek),
			StartHour:   types.Int64PointerValue(db.ServiceConfig.MaintenanceWindow.StartHour),
			StartMinute: types.Int64PointerValue(db.ServiceConfig.MaintenanceWindow.StartMinute),
		}.ToObjectValue(ctx)
	} else {
		maintenanceWindowObjVal = types.ObjectNull(MaintenanceWindowValueV2{}.AttributeTypes(ctx))
	}

	var privateNetworkingObjVal basetypes.ObjectValue
	if db.ApplicationConfig.PrivateNetworking != nil {
		var privateAllowedCidrs types.List
		if db.ApplicationConfig.PrivateNetworking.AllowedCidrs != nil {
			var d diag.Diagnostics
			privateAllowedCidrs, d = types.ListValueFrom(ctx, types.StringType, db.ApplicationConfig.PrivateNetworking.AllowedCidrs)
			diags.Append(d...)
		} else {
			privateAllowedCidrs = types.ListNull(types.StringType)
		}

		var sharedSubnetCIDR types.String
		if db.ApplicationConfig.PrivateNetworking.SharedSubnetCidr != nil {
			sharedSubnetCIDR = types.StringPointerValue(db.ApplicationConfig.PrivateNetworking.SharedSubnetCidr)
		} else {
			sharedSubnetCIDR = types.StringNull()
		}

		privateNetworkingObjVal, _ = PrivateNetworkingValueV2{
			Enabled:          types.BoolPointerValue(db.ApplicationConfig.PrivateNetworking.Enabled),
			Hostname:         types.StringPointerValue(db.ApplicationConfig.PrivateNetworking.Hostname),
			IPAddress:        types.StringPointerValue(db.ApplicationConfig.PrivateNetworking.IpAddress),
			AllowedCIDRs:     privateAllowedCidrs,
			SharedSubnetCIDR: sharedSubnetCIDR,
			SharedSubnetID:   types.StringPointerValue(db.ApplicationConfig.PrivateNetworking.SharedSubnetId),
			SharedNetworkID:  types.StringPointerValue(db.ApplicationConfig.PrivateNetworking.SharedNetworkId),
		}.ToObjectValue(ctx)
	} else {
		privateNetworkingObjVal = types.ObjectNull(PrivateNetworkingValueV2{}.AttributeTypes(ctx))
	}

	var publicNetworkingObjVal basetypes.ObjectValue
	if db.ApplicationConfig.PublicNetworking != nil {
		var publicAllowedCidrs types.List
		if db.ApplicationConfig.PublicNetworking.AllowedCidrs != nil {
			var d diag.Diagnostics
			publicAllowedCidrs, d = types.ListValueFrom(ctx, types.StringType, db.ApplicationConfig.PublicNetworking.AllowedCidrs)
			diags.Append(d...)
		} else {
			publicAllowedCidrs = types.ListNull(types.StringType)
		}

		publicNetworkingObjVal, _ = PublicNetworkingValueV2{
			Enabled:      types.BoolPointerValue(db.ApplicationConfig.PublicNetworking.Enabled),
			Hostname:     types.StringPointerValue(db.ApplicationConfig.PublicNetworking.Hostname),
			IPAddress:    types.StringPointerValue(db.ApplicationConfig.PublicNetworking.IpAddress),
			AllowedCidrs: publicAllowedCidrs,
		}.ToObjectValue(ctx)
	} else {
		publicNetworkingObjVal = types.ObjectNull(PublicNetworkingValueV2{}.AttributeTypes(ctx))
	}

	var targetServiceConfig ServiceConfigValueV2
	targetServiceConfig.Disksize = types.Int64PointerValue(db.ServiceConfig.Disksize)
	targetServiceConfig.ServiceConfigType = types.StringValue(db.ServiceConfig.Type)
	targetServiceConfig.Flavor = types.StringValue(db.ServiceConfig.Flavor)
	targetServiceConfig.Region = types.StringValue(db.ServiceConfig.Region)
	targetServiceConfig.MaintenanceWindow = maintenanceWindowObjVal

	targetServiceConfigObj, diags := targetServiceConfig.ToObjectValue(ctx)

	// Extract password consistently - use plan password for create response
	planPassword := ""
	if passwordAttr, exists := plan.ApplicationConfig.Attributes()["password"]; exists {
		planPassword = strings.Trim(passwordAttr.String(), "\"")
	}

	var targetApplicationConfig ApplicationConfigValueV2
	targetApplicationConfig.ApplicationConfigType = types.StringValue(db.ApplicationConfig.Type)
	targetApplicationConfig.Password = types.StringValue(planPassword) // take this from the plan, since it is not included in the response
	targetApplicationConfig.Instances = types.Int64PointerValue(db.ApplicationConfig.Instances)
	targetApplicationConfig.Version = types.StringValue(db.ApplicationConfig.Version)
	targetApplicationConfig.ScheduledBackups = scheduledBackupsObjVal
	targetApplicationConfig.Recovery = recoveryObjValue
	targetApplicationConfig.PrivateNetworking = privateNetworkingObjVal
	targetApplicationConfig.PublicNetworking = publicNetworkingObjVal

	targetApplicationConfigObj, diags := targetApplicationConfig.ToObjectValue(ctx)

	targetState.Uuid = types.StringValue(db.Uuid)
	targetState.Name = types.StringValue(db.Name)
	targetState.Description = types.StringPointerValue(db.Description)
	targetState.Status = types.StringValue(db.Status)
	targetState.Phase = types.StringValue(db.Phase)
	targetState.ResourceStatus = types.StringValue(db.ResourceStatus)
	targetState.CreatedBy = types.StringValue(db.CreatedBy)
	targetState.CreatedAt = types.StringValue(db.CreatedAt.Format(time.RFC3339))
	targetState.LastModifiedBy = types.StringValue(db.LastModifiedBy)
	targetState.LastModifiedAt = types.StringValue(db.LastModifiedAt.Format(time.RFC3339))
	targetState.ApplicationConfig = targetApplicationConfigObj
	targetState.ServiceConfig = targetServiceConfigObj

	ctx = tflog.SetField(ctx, "conversion_create_target_state", &targetState)
	tflog.Debug(ctx, "Converted api create response to state")

	return diags

}

func psqlGetResponseToModelV2(ctx context.Context, db apiv2.PostgreSQLGetResponse, targetState *DatabaseModelV2, previousState DatabaseModelV2) diag.Diagnostics {

	var diags diag.Diagnostics

	ctx = tflog.SetField(ctx, "conversion_read_source_response", db)
	tflog.Debug(ctx, "Converting read api response")

	var recoveryObjValue basetypes.ObjectValue
	if db.ApplicationConfig.Recovery != nil {
		recoveryObjValue, _ = RecoveryValueV2{
			Exclusive:  types.BoolPointerValue(db.ApplicationConfig.Recovery.Exclusive),
			Source:     types.StringPointerValue(db.ApplicationConfig.Recovery.Source),
			TargetLsn:  types.StringPointerValue(db.ApplicationConfig.Recovery.TargetLsn),
			TargetName: types.StringPointerValue(db.ApplicationConfig.Recovery.TargetName),
			TargetXid:  types.StringPointerValue(db.ApplicationConfig.Recovery.TargetXid),
			TargetTime: types.StringValue(db.ApplicationConfig.Recovery.TargetTime.Format(time.RFC3339)),
		}.ToObjectValue(ctx)
	}

	var scheduledBackupsObjVal basetypes.ObjectValue
	if db.ApplicationConfig.ScheduledBackups != nil && db.ApplicationConfig.ScheduledBackups.Schedule != nil {
		scheduleObjVal, _ := ScheduleValueV2{
			Hour:   types.Int64PointerValue(db.ApplicationConfig.ScheduledBackups.Schedule.Hour),
			Minute: types.Int64PointerValue(db.ApplicationConfig.ScheduledBackups.Schedule.Minute),
		}.ToObjectValue(ctx)

		scheduledBackupsObjVal, _ = ScheduledBackupsValueV2{
			Schedule:  scheduleObjVal,
			Retention: types.Int64PointerValue(db.ApplicationConfig.ScheduledBackups.Retention),
		}.ToObjectValue(ctx)
	} else {
		scheduledBackupsObjVal = types.ObjectNull(ScheduledBackupsValueV2{}.AttributeTypes(ctx))
	}

	var maintenanceWindowObjVal basetypes.ObjectValue
	if db.ServiceConfig.MaintenanceWindow != nil {
		maintenanceWindowObjVal, _ = MaintenanceWindowValueV2{
			DayOfWeek:   types.Int64PointerValue(db.ServiceConfig.MaintenanceWindow.DayOfWeek),
			StartHour:   types.Int64PointerValue(db.ServiceConfig.MaintenanceWindow.StartHour),
			StartMinute: types.Int64PointerValue(db.ServiceConfig.MaintenanceWindow.StartMinute),
		}.ToObjectValue(ctx)
	} else {
		maintenanceWindowObjVal = types.ObjectNull(MaintenanceWindowValueV2{}.AttributeTypes(ctx))
	}

	var privateNetworkingObjVal basetypes.ObjectValue
	if db.ApplicationConfig.PrivateNetworking != nil {
		tflog.Debug(ctx, "got private networking")
		var privateAllowedCidrs types.List
		if db.ApplicationConfig.PrivateNetworking.AllowedCidrs != nil {
			var d diag.Diagnostics
			privateAllowedCidrs, d = types.ListValueFrom(ctx, types.StringType, db.ApplicationConfig.PrivateNetworking.AllowedCidrs)
			diags.Append(d...)
		} else {
			privateAllowedCidrs = types.ListNull(types.StringType)
		}

		var sharedSubnetCIDRRead types.String
		if db.ApplicationConfig.PrivateNetworking.SharedSubnetCidr != nil {
			sharedSubnetCIDRRead = types.StringPointerValue(db.ApplicationConfig.PrivateNetworking.SharedSubnetCidr)
		} else {
			sharedSubnetCIDRRead = types.StringNull()
		}

		privateNetworkingObjVal, _ = PrivateNetworkingValueV2{
			Enabled:          types.BoolPointerValue(db.ApplicationConfig.PrivateNetworking.Enabled),
			Hostname:         types.StringPointerValue(db.ApplicationConfig.PrivateNetworking.Hostname),
			IPAddress:        types.StringPointerValue(db.ApplicationConfig.PrivateNetworking.IpAddress),
			AllowedCIDRs:     privateAllowedCidrs,
			SharedSubnetCIDR: sharedSubnetCIDRRead,
			SharedSubnetID:   types.StringPointerValue(db.ApplicationConfig.PrivateNetworking.SharedSubnetId),
			SharedNetworkID:  types.StringPointerValue(db.ApplicationConfig.PrivateNetworking.SharedNetworkId),
		}.ToObjectValue(ctx)
	} else {
		privateNetworkingObjVal = types.ObjectNull(PrivateNetworkingValueV2{}.AttributeTypes(ctx))
	}

	var publicNetworkingObjVal basetypes.ObjectValue
	if db.ApplicationConfig.PublicNetworking != nil {
		var publicAllowedCidrs types.List
		if db.ApplicationConfig.PublicNetworking.AllowedCidrs != nil {
			var d diag.Diagnostics
			publicAllowedCidrs, d = types.ListValueFrom(ctx, types.StringType, db.ApplicationConfig.PublicNetworking.AllowedCidrs)
			diags.Append(d...)
		} else {
			publicAllowedCidrs = types.ListNull(types.StringType)
		}

		publicNetworkingObjVal, _ = PublicNetworkingValueV2{
			Enabled:      types.BoolPointerValue(db.ApplicationConfig.PublicNetworking.Enabled),
			Hostname:     types.StringPointerValue(db.ApplicationConfig.PublicNetworking.Hostname),
			IPAddress:    types.StringPointerValue(db.ApplicationConfig.PublicNetworking.IpAddress),
			AllowedCidrs: publicAllowedCidrs,
		}.ToObjectValue(ctx)
	} else {
		publicNetworkingObjVal = types.ObjectNull(PublicNetworkingValueV2{}.AttributeTypes(ctx))
	}

	var targetServiceConfig ServiceConfigValueV2
	targetServiceConfig.Disksize = types.Int64PointerValue(db.ServiceConfig.Disksize)
	targetServiceConfig.ServiceConfigType = types.StringValue(db.ServiceConfig.Type)
	targetServiceConfig.Flavor = types.StringValue(db.ServiceConfig.Flavor)
	targetServiceConfig.Region = types.StringValue(db.ServiceConfig.Region)
	targetServiceConfig.MaintenanceWindow = maintenanceWindowObjVal

	targetServiceConfigObj, diags := targetServiceConfig.ToObjectValue(ctx)

	// Extract password consistently - use previous state password for read response
	previousPassword := ""
	if passwordAttr, exists := previousState.ApplicationConfig.Attributes()["password"]; exists {
		previousPassword = strings.Trim(passwordAttr.String(), "\"")
	}

	var targetApplicationConfig ApplicationConfigValueV2
	targetApplicationConfig.ApplicationConfigType = types.StringValue(db.ApplicationConfig.Type)
	targetApplicationConfig.Password = types.StringValue(previousPassword)
	targetApplicationConfig.Instances = types.Int64PointerValue(db.ApplicationConfig.Instances)
	targetApplicationConfig.Version = types.StringValue(db.ApplicationConfig.Version)
	targetApplicationConfig.ScheduledBackups = scheduledBackupsObjVal
	targetApplicationConfig.Recovery = recoveryObjValue
	targetApplicationConfig.PrivateNetworking = privateNetworkingObjVal
	targetApplicationConfig.PublicNetworking = publicNetworkingObjVal

	targetApplicationConfigObj, diags := targetApplicationConfig.ToObjectValue(ctx)

	targetState.Uuid = types.StringValue(db.Uuid)
	targetState.Name = types.StringValue(db.Name)
	targetState.Description = types.StringPointerValue(db.Description)
	targetState.Status = types.StringValue(db.Status)
	targetState.Phase = types.StringValue(db.Phase)
	targetState.ResourceStatus = types.StringValue(db.ResourceStatus)
	targetState.CreatedBy = types.StringValue(db.CreatedBy)
	targetState.CreatedAt = types.StringValue(db.CreatedAt.Format(time.RFC3339))
	targetState.LastModifiedBy = types.StringValue(db.LastModifiedBy)
	targetState.LastModifiedAt = types.StringValue(db.LastModifiedAt.Format(time.RFC3339))
	targetState.ApplicationConfig = targetApplicationConfigObj
	targetState.ServiceConfig = targetServiceConfigObj

	ctx = tflog.SetField(ctx, "conversion_read_target_state", targetState)
	tflog.Debug(ctx, "Converted api read response to state")

	return diags

}
