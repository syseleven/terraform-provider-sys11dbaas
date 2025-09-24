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
	client          *sys11dbaassdk.Client
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

	r.client = providerData.client
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

	psqlRequest := &sys11dbaassdk.GetPostgreSQLRequestV1{
		UUID:         state.Uuid.ValueString(),
		Organization: r.organization.ValueString(),
		Project:      r.project.ValueString(),
	}

	var psqlDB *sys11dbaassdk.GetPostgreSQLResponseV1
	var err error
	errCount := 0
	for {
		psqlDB, err = r.client.GetPostgreSQLDBV1(psqlRequest)
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

	var maintenanceWindow *sys11dbaassdk.MaintenanceWindowV1
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

	var backupSchedule *sys11dbaassdk.PSQLScheduledBackupsV1
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

	createRequest := &sys11dbaassdk.CreatePostgreSQLRequestV1{
		Organization: r.organization.ValueString(),
		Project:      r.project.ValueString(),
		Name:         plan.Name.ValueString(),
		Description:  plan.Description.ValueString(),
		ServiceConfig: &sys11dbaassdk.PSQLServiceConfigRequestV1{
			Disksize:          sys11dbaassdk.Int64ToIntPtr(serviceConfig.Disksize.ValueInt64()),
			Type:              serviceConfig.ServiceConfigType.ValueString(),
			Flavor:            serviceConfig.Flavor.ValueString(),
			Region:            serviceConfig.Region.ValueString(),
			MaintenanceWindow: maintenanceWindow,
			RemoteIPs:         ipList,
		},
		ApplicationConfig: &sys11dbaassdk.PSQLApplicationConfigRequestV1{
			Type:             applicationConfig.ApplicationConfigType.ValueString(),
			Password:         applicationConfig.Password.ValueString(),
			Instances:        sys11dbaassdk.IntPtr(int(applicationConfig.Instances.ValueInt64())),
			Version:          applicationConfig.Version.ValueString(),
			ScheduledBackups: backupSchedule,
		},
	}

	d, _ := json.Marshal(createRequest)
	tflog.Debug(ctx, string(d), nil)

	// Create new db
	createResponse, err := r.client.CreatePostgreSQLDBV1(createRequest)
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
		getRequest := &sys11dbaassdk.GetPostgreSQLRequestV1{
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
			getResponse, err := r.client.GetPostgreSQLDBV1(getRequest)
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

	psqlRequest := &sys11dbaassdk.DeletePostgreSQLRequestV1{
		UUID:         state.Uuid.ValueString(),
		Organization: r.organization.ValueString(),
		Project:      r.project.ValueString(),
	}

	_, err := r.client.DeletePostgreSQLDBV1(psqlRequest)
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

	var maintenanceWindow *sys11dbaassdk.MaintenanceWindowV1
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

	var backupSchedule *sys11dbaassdk.PSQLScheduledBackupsV1
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

	updateRequest := &sys11dbaassdk.UpdatePostgreSQLRequestV1{
		UUID:         state.Uuid.ValueString(),
		Organization: r.organization.ValueString(),
		Project:      r.project.ValueString(),
		Name:         plan.Name.ValueString(),
		Description:  plan.Description.ValueString(),
		ServiceConfig: &sys11dbaassdk.PSQLServiceConfigUpdateRequestV1{
			Disksize:          sys11dbaassdk.Int64ToIntPtr(serviceConfig.Disksize.ValueInt64()),
			Type:              serviceConfig.ServiceConfigType.ValueString(),
			Flavor:            serviceConfig.Flavor.ValueString(),
			MaintenanceWindow: maintenanceWindow,
			RemoteIPs:         ipList,
		},
		ApplicationConfig: &sys11dbaassdk.PSQLApplicationConfigUpdateRequestV1{
			Password:         applicationConfig.Password.ValueString(),
			Instances:        sys11dbaassdk.IntPtr(int(applicationConfig.Instances.ValueInt64())),
			Version:          applicationConfig.Version.ValueString(),
			ScheduledBackups: backupSchedule,
		},
	}

	d, _ := json.Marshal(updateRequest)
	tflog.Debug(ctx, string(d), nil)

	// Update psql
	_, err := r.client.UpdatePostgreSQLDBV1(updateRequest)
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
		getRequest := &sys11dbaassdk.GetPostgreSQLRequestV1{
			Organization: r.organization.ValueString(),
			Project:      r.project.ValueString(),
			UUID:         state.Uuid.ValueString(),
		}
		getResponse, err := r.client.GetPostgreSQLDBV1(getRequest)
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

func psqlCreateResponseToModel(ctx context.Context, db *sys11dbaassdk.CreatePostgreSQLResponseV1, plan DatabaseModel, targetState *DatabaseModel) diag.Diagnostics {

	var diags diag.Diagnostics

	ctx = tflog.SetField(ctx, "conversion_create_source_response", db)
	tflog.Debug(ctx, "Converting create api response")

	ipList, d := types.ListValueFrom(ctx, types.StringType, (*db).ServiceConfig.RemoteIPs)

	var recoveryObjValue basetypes.ObjectValue
	if db.ApplicationConfig.Recovery != nil {
		recoveryObjValue, _ = RecoveryValue{
			Exclusive:  types.BoolValue(db.ApplicationConfig.Recovery.Exclusive),
			Source:     types.StringValue(db.ApplicationConfig.Recovery.Source),
			TargetLsn:  types.StringValue(db.ApplicationConfig.Recovery.TargetLSN),
			TargetName: types.StringValue(db.ApplicationConfig.Recovery.TargetName),
			TargetXid:  types.StringValue(db.ApplicationConfig.Recovery.TargetXID),
			TargetTime: types.StringValue(db.ApplicationConfig.Recovery.TargetTime.Format(time.RFC3339)),
		}.ToObjectValue(ctx)
	}

	scheduleObjVal, _ := ScheduleValue{
		Hour:   types.Int64Value(int64(*db.ApplicationConfig.ScheduledBackups.Schedule.Hour)),
		Minute: types.Int64Value(int64(*db.ApplicationConfig.ScheduledBackups.Schedule.Minute)),
	}.ToObjectValue(ctx)

	scheduledBackupsObjVal, _ := ScheduledBackupsValue{
		Schedule:  scheduleObjVal,
		Retention: types.Int64Value(int64(*db.ApplicationConfig.ScheduledBackups.Retention)),
	}.ToObjectValue(ctx)

	maintenanceWindowObjVal, _ := MaintenanceWindowValue{
		DayOfWeek:   types.Int64Value(int64(*db.ServiceConfig.MaintenanceWindow.DayOfWeek)),
		StartHour:   types.Int64Value(int64(*db.ServiceConfig.MaintenanceWindow.StartHour)),
		StartMinute: types.Int64Value(int64(*db.ServiceConfig.MaintenanceWindow.StartMinute)),
	}.ToObjectValue(ctx)

	diags.Append(d...)

	var targetServiceConfig ServiceConfigValue
	targetServiceConfig.Disksize = types.Int64Value(int64(*db.ServiceConfig.Disksize))
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
	targetApplicationConfig.Instances = types.Int64Value(int64(*db.ApplicationConfig.Instances))
	targetApplicationConfig.Version = types.StringValue(db.ApplicationConfig.Version)
	targetApplicationConfig.Hostname = types.StringValue(db.ApplicationConfig.Hostname)
	targetApplicationConfig.IpAddress = types.StringValue(db.ApplicationConfig.IPAddress)
	targetApplicationConfig.ScheduledBackups = scheduledBackupsObjVal
	targetApplicationConfig.Recovery = recoveryObjValue

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

func psqlGetResponseToModel(ctx context.Context, db *sys11dbaassdk.GetPostgreSQLResponseV1, targetState *DatabaseModel, previousState DatabaseModel) diag.Diagnostics {

	var diags diag.Diagnostics

	ctx = tflog.SetField(ctx, "conversion_read_source_response", db)
	tflog.Debug(ctx, "Converting read api response")

	ipList, d := types.ListValueFrom(ctx, types.StringType, (*db).ServiceConfig.RemoteIPs)

	var recoveryObjValue basetypes.ObjectValue
	if db.ApplicationConfig.Recovery != nil {
		recoveryObjValue, _ = RecoveryValue{
			Exclusive:  types.BoolValue(db.ApplicationConfig.Recovery.Exclusive),
			Source:     types.StringValue(db.ApplicationConfig.Recovery.Source),
			TargetLsn:  types.StringValue(db.ApplicationConfig.Recovery.TargetLSN),
			TargetName: types.StringValue(db.ApplicationConfig.Recovery.TargetName),
			TargetXid:  types.StringValue(db.ApplicationConfig.Recovery.TargetXID),
			TargetTime: types.StringValue(db.ApplicationConfig.Recovery.TargetTime.Format(time.RFC3339)),
		}.ToObjectValue(ctx)
	}

	scheduleObjVal, _ := ScheduleValue{
		Hour:   types.Int64Value(int64(*db.ApplicationConfig.ScheduledBackups.Schedule.Hour)),
		Minute: types.Int64Value(int64(*db.ApplicationConfig.ScheduledBackups.Schedule.Minute)),
	}.ToObjectValue(ctx)

	scheduledBackupsObjVal, _ := ScheduledBackupsValue{
		Schedule:  scheduleObjVal,
		Retention: types.Int64Value(int64(*db.ApplicationConfig.ScheduledBackups.Retention)),
	}.ToObjectValue(ctx)

	maintenanceWindowObjVal, _ := MaintenanceWindowValue{
		DayOfWeek:   types.Int64Value(int64(*db.ServiceConfig.MaintenanceWindow.DayOfWeek)),
		StartHour:   types.Int64Value(int64(*db.ServiceConfig.MaintenanceWindow.StartHour)),
		StartMinute: types.Int64Value(int64(*db.ServiceConfig.MaintenanceWindow.StartMinute)),
	}.ToObjectValue(ctx)

	diags.Append(d...)

	var targetServiceConfig ServiceConfigValue
	targetServiceConfig.Disksize = types.Int64Value(int64(*db.ServiceConfig.Disksize))
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
	targetApplicationConfig.Instances = types.Int64Value(int64(*db.ApplicationConfig.Instances))
	targetApplicationConfig.Version = types.StringValue(db.ApplicationConfig.Version)
	targetApplicationConfig.Hostname = types.StringValue(db.ApplicationConfig.Hostname)
	targetApplicationConfig.IpAddress = types.StringValue(db.ApplicationConfig.IPAddress)
	targetApplicationConfig.ScheduledBackups = scheduledBackupsObjVal
	targetApplicationConfig.Recovery = recoveryObjValue

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
