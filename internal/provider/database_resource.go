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
	database "github.com/syseleven/sys11dbaas-sdk/database/v1"
)

const CREATE_RETRY_LIMIT = 30 * time.Minute

type DatabaseModel struct {
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

type DatabaseResource struct {
	client          *database.TypedClient
	project         types.String
	organization    types.String
	waitForCreation types.Bool
}

func NewDatabaseResource() resource.Resource {
	return &DatabaseResource{}
}

// Metadata returns the resource type name.
func (r *DatabaseResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database"
}

// Configure adds the provider configured client to the resource.
func (r *DatabaseResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = providerData.client.V1()
	r.organization = providerData.organization
	r.project = providerData.project
	r.waitForCreation = providerData.waitForCreation
}

// Read resource information.
func (r *DatabaseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state DatabaseModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var psqlDB database.PostgreSQLGetResponseV1
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

	diags = psqlGetResponseToModel(ctx, psqlDB, &state, state)
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
func (r *DatabaseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan DatabaseModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceConfig, diags := NewServiceConfigValue(plan.ServiceConfig.AttributeTypes(ctx), plan.ServiceConfig.Attributes())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var ipList []string

	for _, e := range serviceConfig.RemoteIps.Elements() {
		ipList = append(ipList, strings.Trim(e.String(), "\""))
	}

	var maintenanceWindow *database.PostgreSQLMaintenance
	if !serviceConfig.MaintenanceWindow.IsUnknown() {
		maintenanceWindowObj, diags := NewMaintenanceWindowValue(serviceConfig.MaintenanceWindow.AttributeTypes(ctx), serviceConfig.MaintenanceWindow.Attributes())
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

	applicationConfig, diags := NewApplicationConfigValue(plan.ApplicationConfig.AttributeTypes(ctx), plan.ApplicationConfig.Attributes())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var backupSchedule *database.PostgreSQLBackupSchedule
	if !applicationConfig.ScheduledBackups.IsUnknown() {
		scheduledBackupsObj, diags := NewScheduledBackupsValue(applicationConfig.ScheduledBackups.AttributeTypes(ctx), applicationConfig.ScheduledBackups.Attributes())
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

	createRequest := database.PostgreSQLCreateRequestV1{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueStringPointer(),
		ServiceConfig: database.PostgreSQLServiceConfigV1{
			Disksize:          serviceConfig.Disksize.ValueInt64Pointer(),
			Type:              serviceConfig.ServiceConfigType.ValueString(),
			Flavor:            serviceConfig.Flavor.ValueString(),
			Region:            serviceConfig.Region.ValueString(),
			MaintenanceWindow: maintenanceWindow,
			RemoteIps:         &ipList,
		},
		ApplicationConfig: database.PostgreSQLApplicationConfigV1{
			Type:             applicationConfig.ApplicationConfigType.ValueString(),
			Password:         applicationConfig.Password.ValueString(),
			Instances:        applicationConfig.Instances.ValueInt64Pointer(),
			Version:          applicationConfig.Version.ValueString(),
			ScheduledBackups: backupSchedule,
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
	targetState := DatabaseModel{}
	diags = psqlCreateResponseToModel(ctx, createResponse, plan, &targetState)

	ctx = tflog.SetField(ctx, "create_target_state", &targetState)
	tflog.Debug(ctx, "[CREATE] Created", nil)
	resp.Diagnostics.Append(diags...)

	tflog.Debug(ctx, "[CREATE] Wait for creation: "+r.waitForCreation.String(), nil)
	if r.waitForCreation.ValueBool() {
		sleepFor := time.Duration(30 * time.Second)
		errCount := 0
		for retryCount := 0; targetState.Status.ValueString() != database.StateReady; {
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
			diags = psqlGetResponseToModel(ctx, getResponse, &targetState, plan)

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
func (r *DatabaseResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Get current state
	var state DatabaseModel
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
func (r *DatabaseResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan
	var plan DatabaseModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Get current state
	var state DatabaseModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceConfig, diags := NewServiceConfigValue(plan.ServiceConfig.AttributeTypes(ctx), plan.ServiceConfig.Attributes())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var ipList []string

	for _, e := range serviceConfig.RemoteIps.Elements() {
		ipList = append(ipList, strings.Trim(e.String(), "\""))
	}

	var maintenanceWindow *database.PostgreSQLMaintenance
	if !serviceConfig.MaintenanceWindow.IsUnknown() {
		maintenanceWindowObj, diags := NewMaintenanceWindowValue(serviceConfig.MaintenanceWindow.AttributeTypes(ctx), serviceConfig.MaintenanceWindow.Attributes())
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

	applicationConfig, diags := NewApplicationConfigValue(plan.ApplicationConfig.AttributeTypes(ctx), plan.ApplicationConfig.Attributes())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var backupSchedule *database.PostgreSQLBackupSchedule
	if !applicationConfig.ScheduledBackups.IsUnknown() {
		scheduledBackupsObj, diags := NewScheduledBackupsValue(applicationConfig.ScheduledBackups.AttributeTypes(ctx), applicationConfig.ScheduledBackups.Attributes())
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

	updateRequest := database.PostgreSQLUpdateRequestV1{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueStringPointer(),
		ServiceConfig: database.PostgreSQLServiceConfigV1{
			Disksize:          serviceConfig.Disksize.ValueInt64Pointer(),
			Type:              serviceConfig.ServiceConfigType.ValueString(),
			Flavor:            serviceConfig.Flavor.ValueString(),
			MaintenanceWindow: maintenanceWindow,
			RemoteIps:         &ipList,
		},
		ApplicationConfig: database.PostgreSQLApplicationConfigV1{
			Password:         applicationConfig.Password.ValueString(),
			Instances:        applicationConfig.Instances.ValueInt64Pointer(),
			Version:          applicationConfig.Version.ValueString(),
			ScheduledBackups: backupSchedule,
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

		// Update psql
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
		targetState := DatabaseModel{}
		diags = psqlGetResponseToModel(ctx, getResponse, &targetState, plan)

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
func (r *DatabaseResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = DatabaseResourceSchema(ctx)
}

func psqlCreateResponseToModel(ctx context.Context, db database.PostgreSQLGetResponseV1, plan DatabaseModel, targetState *DatabaseModel) diag.Diagnostics {

	var diags diag.Diagnostics

	ctx = tflog.SetField(ctx, "conversion_create_source_response", db)
	tflog.Debug(ctx, "Converting create api response")

	ipList, d := types.ListValueFrom(ctx, types.StringType, db.ServiceConfig.RemoteIps)

	var recoveryObjValue basetypes.ObjectValue
	if db.ApplicationConfig.Recovery != nil {
		recoveryObjValue, _ = RecoveryValue{
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
		scheduleObjVal, _ := ScheduleValue{
			Hour:   types.Int64PointerValue(db.ApplicationConfig.ScheduledBackups.Schedule.Hour),
			Minute: types.Int64PointerValue(db.ApplicationConfig.ScheduledBackups.Schedule.Minute),
		}.ToObjectValue(ctx)

		scheduledBackupsObjVal, _ = ScheduledBackupsValue{
			Schedule:  scheduleObjVal,
			Retention: types.Int64PointerValue(db.ApplicationConfig.ScheduledBackups.Retention),
		}.ToObjectValue(ctx)
	} else {
		scheduledBackupsObjVal = types.ObjectNull(ScheduledBackupsValue{}.AttributeTypes(ctx))
	}

	var maintenanceWindowObjVal basetypes.ObjectValue
	if db.ServiceConfig.MaintenanceWindow != nil {
		maintenanceWindowObjVal, _ = MaintenanceWindowValue{
			DayOfWeek:   types.Int64PointerValue(db.ServiceConfig.MaintenanceWindow.DayOfWeek),
			StartHour:   types.Int64PointerValue(db.ServiceConfig.MaintenanceWindow.StartHour),
			StartMinute: types.Int64PointerValue(db.ServiceConfig.MaintenanceWindow.StartMinute),
		}.ToObjectValue(ctx)
	} else {
		maintenanceWindowObjVal = types.ObjectNull(MaintenanceWindowValue{}.AttributeTypes(ctx))
	}

	diags.Append(d...)

	var targetServiceConfig ServiceConfigValue
	targetServiceConfig.Disksize = types.Int64PointerValue(db.ServiceConfig.Disksize)
	targetServiceConfig.ServiceConfigType = types.StringValue(db.ServiceConfig.Type)
	targetServiceConfig.Flavor = types.StringValue(db.ServiceConfig.Flavor)
	targetServiceConfig.Region = types.StringValue(db.ServiceConfig.Region)
	targetServiceConfig.MaintenanceWindow = maintenanceWindowObjVal
	targetServiceConfig.RemoteIps = ipList

	targetServiceConfigObj, diags := targetServiceConfig.ToObjectValue(ctx)

	// Extract password from plan for state consistency
	// The API response doesn't include the password for security reasons, but Terraform
	// needs to store it in state to detect changes. We must use the planned password
	// value to ensure the sensitive application_config object remains consistent
	// between plan and apply phases, preventing "inconsistent values for sensitive attribute" errors.
	planPassword := ""
	if passwordAttr, exists := plan.ApplicationConfig.Attributes()["password"]; exists {
		planPassword = strings.Trim(passwordAttr.String(), "\"")
	}

	var targetApplicationConfig ApplicationConfigValue
	targetApplicationConfig.ApplicationConfigType = types.StringValue(db.ApplicationConfig.Type)
	targetApplicationConfig.Password = types.StringValue(planPassword) // take this from the plan, since it is not included in the response
	targetApplicationConfig.Instances = types.Int64PointerValue(db.ApplicationConfig.Instances)
	targetApplicationConfig.Version = types.StringValue(db.ApplicationConfig.Version)
	targetApplicationConfig.Hostname = types.StringPointerValue(db.ApplicationConfig.Hostname)
	targetApplicationConfig.IpAddress = types.StringPointerValue(db.ApplicationConfig.IpAddress)
	targetApplicationConfig.ScheduledBackups = scheduledBackupsObjVal
	targetApplicationConfig.Recovery = recoveryObjValue

	targetApplicationConfigObj, diags := targetApplicationConfig.ToObjectValue(ctx)

	targetState.Uuid = types.StringValue(db.Uuid)
	targetState.Name = types.StringValue(db.Name)
	targetState.Description = types.StringValue(db.Description)
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

func psqlGetResponseToModel(ctx context.Context, db database.PostgreSQLGetResponseV1, targetState *DatabaseModel, previousState DatabaseModel) diag.Diagnostics {

	var diags diag.Diagnostics

	ctx = tflog.SetField(ctx, "conversion_read_source_response", db)
	tflog.Debug(ctx, "Converting read api response")

	ipList, d := types.ListValueFrom(ctx, types.StringType, db.ServiceConfig.RemoteIps)

	var recoveryObjValue basetypes.ObjectValue
	if db.ApplicationConfig.Recovery != nil {
		recoveryObjValue, _ = RecoveryValue{
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
		scheduleObjVal, _ := ScheduleValue{
			Hour:   types.Int64PointerValue(db.ApplicationConfig.ScheduledBackups.Schedule.Hour),
			Minute: types.Int64PointerValue(db.ApplicationConfig.ScheduledBackups.Schedule.Minute),
		}.ToObjectValue(ctx)

		scheduledBackupsObjVal, _ = ScheduledBackupsValue{
			Schedule:  scheduleObjVal,
			Retention: types.Int64PointerValue(db.ApplicationConfig.ScheduledBackups.Retention),
		}.ToObjectValue(ctx)
	} else {
		scheduledBackupsObjVal = types.ObjectNull(ScheduledBackupsValue{}.AttributeTypes(ctx))
	}

	var maintenanceWindowObjVal basetypes.ObjectValue
	if db.ServiceConfig.MaintenanceWindow != nil {
		maintenanceWindowObjVal, _ = MaintenanceWindowValue{
			DayOfWeek:   types.Int64PointerValue(db.ServiceConfig.MaintenanceWindow.DayOfWeek),
			StartHour:   types.Int64PointerValue(db.ServiceConfig.MaintenanceWindow.StartHour),
			StartMinute: types.Int64PointerValue(db.ServiceConfig.MaintenanceWindow.StartMinute),
		}.ToObjectValue(ctx)
	} else {
		maintenanceWindowObjVal = types.ObjectNull(MaintenanceWindowValue{}.AttributeTypes(ctx))
	}

	diags.Append(d...)

	var targetServiceConfig ServiceConfigValue
	targetServiceConfig.Disksize = types.Int64PointerValue(db.ServiceConfig.Disksize)
	targetServiceConfig.ServiceConfigType = types.StringValue(db.ServiceConfig.Type)
	targetServiceConfig.Flavor = types.StringValue(db.ServiceConfig.Flavor)
	targetServiceConfig.Region = types.StringValue(db.ServiceConfig.Region)
	targetServiceConfig.MaintenanceWindow = maintenanceWindowObjVal
	targetServiceConfig.RemoteIps = ipList

	targetServiceConfigObj, diags := targetServiceConfig.ToObjectValue(ctx)

	// Extract password from previous state for consistency
	// During read operations, the API doesn't return the password. We must preserve
	// the password from the previous Terraform state to maintain consistency of the
	// sensitive application_config object and prevent state drift or validation errors.
	previousPassword := ""
	if passwordAttr, exists := previousState.ApplicationConfig.Attributes()["password"]; exists {
		previousPassword = strings.Trim(passwordAttr.String(), "\"")
	}

	var targetApplicationConfig ApplicationConfigValue
	targetApplicationConfig.ApplicationConfigType = types.StringValue(db.ApplicationConfig.Type)
	targetApplicationConfig.Password = types.StringValue(previousPassword)
	targetApplicationConfig.Instances = types.Int64PointerValue(db.ApplicationConfig.Instances)
	targetApplicationConfig.Version = types.StringValue(db.ApplicationConfig.Version)
	targetApplicationConfig.Hostname = types.StringPointerValue(db.ApplicationConfig.Hostname)
	targetApplicationConfig.IpAddress = types.StringPointerValue(db.ApplicationConfig.IpAddress)
	targetApplicationConfig.ScheduledBackups = scheduledBackupsObjVal
	targetApplicationConfig.Recovery = recoveryObjValue

	targetApplicationConfigObj, diags := targetApplicationConfig.ToObjectValue(ctx)

	targetState.Uuid = types.StringValue(db.Uuid)
	targetState.Name = types.StringValue(db.Name)
	targetState.Description = types.StringValue(db.Description)
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
