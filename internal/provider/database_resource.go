package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	database "github.com/syseleven/sys11dbaas-sdk/database/v1"
)

const CREATE_RETRY_LIMIT = 30 * time.Minute

type ApplicationConfigModelV1 struct {
	Hostname              basetypes.StringValue `tfsdk:"hostname"`
	Instances             basetypes.Int64Value  `tfsdk:"instances"`
	IpAddress             basetypes.StringValue `tfsdk:"ip_address"`
	Password              basetypes.StringValue `tfsdk:"password"`
	Recovery              basetypes.ObjectValue `tfsdk:"recovery"`
	ScheduledBackups      basetypes.ObjectValue `tfsdk:"scheduled_backups"`
	ApplicationConfigType basetypes.StringValue `tfsdk:"type"`
	Version               basetypes.StringValue `tfsdk:"version"`
}

func (m ApplicationConfigModelV1) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"hostname":   types.StringType,
		"instances":  types.Int64Type,
		"ip_address": types.StringType,
		"password":   types.StringType,
		"recovery": types.ObjectType{
			AttrTypes: RecoveryModel{}.AttributeTypes(),
		},
		"scheduled_backups": types.ObjectType{
			AttrTypes: ScheduledBackupsModel{}.AttributeTypes(),
		},
		"type":    types.StringType,
		"version": types.StringType,
	}
}

type ServiceConfigModelV1 struct {
	Disksize          basetypes.Int64Value  `tfsdk:"disksize"`
	Flavor            basetypes.StringValue `tfsdk:"flavor"`
	MaintenanceWindow basetypes.ObjectValue `tfsdk:"maintenance_window"`
	Region            basetypes.StringValue `tfsdk:"region"`
	RemoteIps         basetypes.ListValue   `tfsdk:"remote_ips"`
	ServiceConfigType basetypes.StringValue `tfsdk:"type"`
}

func (m ServiceConfigModelV1) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"disksize": types.Int64Type,
		"flavor":   types.StringType,
		"maintenance_window": types.ObjectType{
			AttrTypes: MaintenanceWindowModel{}.AttributeTypes(),
		},
		"region": types.StringType,
		"remote_ips": types.ListType{
			ElemType: types.StringType,
		},
		"type": types.StringType,
	}
}

type DatabaseModel struct {
	ApplicationConfig types.Object      `tfsdk:"application_config"`
	CreatedAt         timetypes.RFC3339 `tfsdk:"created_at"`
	CreatedBy         types.String      `tfsdk:"created_by"`
	Description       types.String      `tfsdk:"description"`
	LastModifiedAt    timetypes.RFC3339 `tfsdk:"last_modified_at"`
	LastModifiedBy    types.String      `tfsdk:"last_modified_by"`
	Name              types.String      `tfsdk:"name"`
	ServiceConfig     types.Object      `tfsdk:"service_config"`
	Status            types.String      `tfsdk:"status"`
	Phase             types.String      `tfsdk:"phase"`
	ResourceStatus    types.String      `tfsdk:"resource_status"`
	Uuid              types.String      `tfsdk:"uuid"`
}

// resource

type DatabaseResource struct {
	client          *database.TypedClient
	project         types.String
	organization    types.String
	waitForCreation types.Bool
	waitForUpdate   types.Bool
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
	r.waitForUpdate = providerData.waitForUpdate
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

	psqlDB, err := r.client.GetPostgreSQL(ctx, r.organization.ValueString(), r.project.ValueString(), state.Uuid.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read database",
			err.Error(),
		)
		return
	}

	diags = psqlGetResponseToModel(ctx, psqlDB, &state)
	ctx = tflog.SetField(ctx, "read_target_state", state)
	tflog.Debug(ctx, "Reading database", nil)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
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

	var applicationConfig ApplicationConfigModelV1
	diags = plan.ApplicationConfig.As(ctx, &applicationConfig, basetypes.ObjectAsOptions{})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var scheduledBackups *database.PostgreSQLBackupSchedule
	if !applicationConfig.ScheduledBackups.IsUnknown() {
		var scheduledBackupsModel ScheduledBackupsModel
		resp.Diagnostics.Append(applicationConfig.ScheduledBackups.As(ctx, &scheduledBackupsModel, basetypes.ObjectAsOptions{})...)
		if resp.Diagnostics.HasError() {
			return
		}

		var schedule ScheduleModel
		resp.Diagnostics.Append(scheduledBackupsModel.Schedule.As(ctx, &schedule, basetypes.ObjectAsOptions{})...)
		if resp.Diagnostics.HasError() {
			return
		}

		scheduledBackups = &database.PostgreSQLBackupSchedule{
			Retention: scheduledBackupsModel.Retention.ValueInt64Pointer(),
			Schedule: &database.PostgreSQLBackupScheduleConfig{
				Hour:   schedule.Hour.ValueInt64Pointer(),
				Minute: schedule.Minute.ValueInt64Pointer(),
			},
		}
	}

	var recovery *database.PostgreSQLRecovery
	if !applicationConfig.Recovery.IsUnknown() && !applicationConfig.Recovery.IsNull() {
		var recoveryModel RecoveryModel
		resp.Diagnostics.Append(applicationConfig.Recovery.As(ctx, &recoveryModel, basetypes.ObjectAsOptions{})...)
		if resp.Diagnostics.HasError() {
			return
		}

		targetTime, err := time.Parse(time.RFC3339, recoveryModel.TargetTime.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("failed to parse recovery time", "Parsing recovery target time into RFC339 time failed")
			return
		}

		recovery = &database.PostgreSQLRecovery{
			Exclusive:  recoveryModel.Exclusive.ValueBoolPointer(),
			Source:     recoveryModel.Source.ValueStringPointer(),
			TargetLsn:  recoveryModel.TargetLsn.ValueStringPointer(),
			TargetName: recoveryModel.TargetName.ValueStringPointer(),
			TargetTime: &targetTime,
			TargetXid:  recoveryModel.TargetXid.ValueStringPointer(),
		}
	}

	var serviceConfig ServiceConfigModelV1
	resp.Diagnostics.Append(plan.ServiceConfig.As(ctx, &serviceConfig, basetypes.ObjectAsOptions{})...)
	if resp.Diagnostics.HasError() {
		return
	}

	var remoteIps []string
	if !serviceConfig.RemoteIps.IsNull() && !serviceConfig.RemoteIps.IsUnknown() {
		resp.Diagnostics.Append(serviceConfig.RemoteIps.ElementsAs(ctx, &remoteIps, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	var maintenanceWindow *database.PostgreSQLMaintenance
	if !serviceConfig.MaintenanceWindow.IsUnknown() {
		var maintenanceWindowModel MaintenanceWindowModel
		resp.Diagnostics.Append(serviceConfig.MaintenanceWindow.As(ctx, &maintenanceWindowModel, basetypes.ObjectAsOptions{})...)
		if resp.Diagnostics.HasError() {
			return
		}

		maintenanceWindow = &database.PostgreSQLMaintenance{
			DayOfWeek:   maintenanceWindowModel.DayOfWeek.ValueInt64Pointer(),
			StartHour:   maintenanceWindowModel.StartHour.ValueInt64Pointer(),
			StartMinute: maintenanceWindowModel.StartMinute.ValueInt64Pointer(),
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
			RemoteIps:         &remoteIps,
		},
		ApplicationConfig: database.PostgreSQLApplicationConfigV1{
			Type:             applicationConfig.ApplicationConfigType.ValueString(),
			Password:         applicationConfig.Password.ValueString(),
			Instances:        applicationConfig.Instances.ValueInt64Pointer(),
			Version:          applicationConfig.Version.ValueString(),
			ScheduledBackups: scheduledBackups,
			Recovery:         recovery,
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
	diags = psqlGetResponseToModel(ctx, createResponse, &plan)

	if r.waitForCreation.ValueBool() {
		for createResponse.Status != database.StateReady {
			select {
			case <-ctx.Done():
				return
			default:
				time.Sleep(30 * time.Second)
			}
			createResponse, err = r.client.GetPostgreSQL(ctx, r.organization.ValueString(), r.project.ValueString(), createResponse.Uuid)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error read database during wait",
					"Could not read database during wait, unexpected error: "+err.Error(),
				)
				return
			}
		}
	}

	diags = psqlGetResponseToModel(ctx, createResponse, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
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

	var applicationConfig ApplicationConfigModelV1
	diags = plan.ApplicationConfig.As(ctx, &applicationConfig, basetypes.ObjectAsOptions{})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var scheduledBackups *database.PostgreSQLBackupSchedule
	if !applicationConfig.ScheduledBackups.IsUnknown() {
		var scheduledBackupsModel ScheduledBackupsModel
		resp.Diagnostics.Append(applicationConfig.ScheduledBackups.As(ctx, &scheduledBackupsModel, basetypes.ObjectAsOptions{})...)
		if resp.Diagnostics.HasError() {
			return
		}

		var schedule ScheduleModel
		resp.Diagnostics.Append(scheduledBackupsModel.Schedule.As(ctx, &schedule, basetypes.ObjectAsOptions{})...)
		if resp.Diagnostics.HasError() {
			return
		}

		scheduledBackups = &database.PostgreSQLBackupSchedule{
			Retention: scheduledBackupsModel.Retention.ValueInt64Pointer(),
			Schedule: &database.PostgreSQLBackupScheduleConfig{
				Hour:   schedule.Hour.ValueInt64Pointer(),
				Minute: schedule.Minute.ValueInt64Pointer(),
			},
		}
	}

	var recovery *database.PostgreSQLRecovery
	if !applicationConfig.Recovery.IsUnknown() && !applicationConfig.Recovery.IsNull() {
		var recoveryModel RecoveryModel
		resp.Diagnostics.Append(applicationConfig.Recovery.As(ctx, &recoveryModel, basetypes.ObjectAsOptions{})...)
		if resp.Diagnostics.HasError() {
			return
		}

		targetTime, err := time.Parse(time.RFC3339, recoveryModel.TargetTime.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("failed to parse recovery time", "Parsing recovery target time into RFC339 time failed")
			return
		}

		recovery = &database.PostgreSQLRecovery{
			Exclusive:  recoveryModel.Exclusive.ValueBoolPointer(),
			Source:     recoveryModel.Source.ValueStringPointer(),
			TargetLsn:  recoveryModel.TargetLsn.ValueStringPointer(),
			TargetName: recoveryModel.TargetName.ValueStringPointer(),
			TargetTime: &targetTime,
			TargetXid:  recoveryModel.TargetXid.ValueStringPointer(),
		}
	}

	var serviceConfig ServiceConfigModelV1
	resp.Diagnostics.Append(plan.ServiceConfig.As(ctx, &serviceConfig, basetypes.ObjectAsOptions{})...)
	if resp.Diagnostics.HasError() {
		return
	}

	var remoteIps []string
	for _, e := range serviceConfig.RemoteIps.Elements() {
		remoteIps = append(remoteIps, strings.Trim(e.String(), "\""))
	}

	var maintenanceWindow *database.PostgreSQLMaintenance
	if !serviceConfig.MaintenanceWindow.IsUnknown() {
		var maintenanceWindowModel MaintenanceWindowModel
		resp.Diagnostics.Append(serviceConfig.MaintenanceWindow.As(ctx, &maintenanceWindowModel, basetypes.ObjectAsOptions{})...)
		if resp.Diagnostics.HasError() {
			return
		}

		maintenanceWindow = &database.PostgreSQLMaintenance{
			DayOfWeek:   maintenanceWindowModel.DayOfWeek.ValueInt64Pointer(),
			StartHour:   maintenanceWindowModel.StartHour.ValueInt64Pointer(),
			StartMinute: maintenanceWindowModel.StartMinute.ValueInt64Pointer(),
		}
	}

	updateRequest := database.PostgreSQLUpdateRequestV1{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueStringPointer(),
		ServiceConfig: database.PostgreSQLServiceConfigV1{
			Disksize:          serviceConfig.Disksize.ValueInt64Pointer(),
			Type:              serviceConfig.ServiceConfigType.ValueString(),
			Flavor:            serviceConfig.Flavor.ValueString(),
			Region:            serviceConfig.Region.ValueString(),
			MaintenanceWindow: maintenanceWindow,
			RemoteIps:         &remoteIps,
		},
		ApplicationConfig: database.PostgreSQLApplicationConfigV1{
			Type:             applicationConfig.ApplicationConfigType.ValueString(),
			Password:         applicationConfig.Password.ValueString(),
			Instances:        applicationConfig.Instances.ValueInt64Pointer(),
			Version:          applicationConfig.Version.ValueString(),
			ScheduledBackups: scheduledBackups,
			Recovery:         recovery,
		},
	}

	d, _ := json.Marshal(updateRequest)
	tflog.Debug(ctx, string(d), nil)

	// Update psql
	_, err := r.client.UpdatePostgreSQL(ctx, r.organization.ValueString(), r.project.ValueString(), plan.Uuid.ValueString(), updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating database",
			"Could not update database, unexpected error: "+err.Error(),
		)
		resp.Diagnostics.Append(diags...)
		return
	}

	response := database.PostgreSQLGetResponseV1{Status: database.StateNotReady}
	if r.waitForUpdate.ValueBool() {
		for response.Status != database.StateReady {
			select {
			case <-ctx.Done():
				return
			default:
				time.Sleep(30 * time.Second)
			}
			response, err = r.client.GetPostgreSQL(ctx, r.organization.ValueString(), r.project.ValueString(), plan.Uuid.ValueString())
			if err != nil {
				resp.Diagnostics.AddError(
					"Error updating database",
					"Could not update database, unexpected error: "+err.Error(),
				)
				return
			}
		}
	} else {
		response, err = r.client.GetPostgreSQL(ctx, r.organization.ValueString(), r.project.ValueString(), plan.Uuid.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating database",
				"Could not update database, unexpected error: "+err.Error(),
			)
			return
		}
	}

	diags = psqlGetResponseToModel(ctx, response, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Schema defines the schema for the resource.
func (r *DatabaseResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"application_config": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"hostname": schema.StringAttribute{
						Computed:    true,
						Description: "DNS name of the database in the format uuid.postgresql.syseleven.services.",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"instances": schema.Int64Attribute{
						Required:    true,
						Description: "Node count of the database cluster.",
						Validators: []validator.Int64{
							int64validator.AtMost(5),
						},
					},
					"ip_address": schema.StringAttribute{
						Computed:    true,
						Description: "Public IP address of the database. It will be 'Pending' if no address has been assigned yet.",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"password": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Sensitive:   true,
						Description: "Password for the admin user.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(16),
						},
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"recovery": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"exclusive": schema.BoolAttribute{
								Optional:    true,
								Computed:    true,
								Description: "Set to true, when the given target should be excluded.",
								Default:     booldefault.StaticBool(false),
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
								},
							},
							"source": schema.StringAttribute{
								Optional:    true,
								Computed:    true,
								Description: "UUID of the source database.",
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
							"target_lsn": schema.StringAttribute{
								Optional:    true,
								Computed:    true,
								Description: "LSN of the write-ahead log location up to which recovery will proceed. target_* parameters are mutually exclusive.",
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
							"target_name": schema.StringAttribute{
								Optional:    true,
								Computed:    true,
								Description: "Named restore point (created with pg_create_restore_point()) to which recovery will proceed. target_* parameters are mutually exclusive.",
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
							"target_time": schema.StringAttribute{
								Optional:    true,
								Computed:    true,
								Description: "Time stamp up to which recovery will proceed, expressed in RFC 3339 format. target_* parameters are mutually exclusive.",
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
							"target_xid": schema.StringAttribute{
								Optional:    true,
								Computed:    true,
								Description: "Transaction ID up to which recovery will proceed. target_* parameters are mutually exclusive.",
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
						},
						Optional: true,
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
							objectplanmodifier.RequiresReplaceIfConfigured(),
						},
					},
					"scheduled_backups": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"retention": schema.Int64Attribute{
								Optional:    true,
								Computed:    true,
								Description: "Duration in days for which backups should be stored.",
								Validators: []validator.Int64{
									int64validator.Between(7, 90),
								},
								Default: int64default.StaticInt64(7),
								PlanModifiers: []planmodifier.Int64{
									int64planmodifier.UseStateForUnknown(),
								},
							},
							"schedule": schema.SingleNestedAttribute{
								Attributes: map[string]schema.Attribute{
									"hour": schema.Int64Attribute{
										Optional:    true,
										Computed:    true,
										Description: "Hour when the full backup should start. If this value is omitted, a random hour between 1am and 5am will be generated.",
										Validators: []validator.Int64{
											int64validator.Between(0, 23),
										},
										PlanModifiers: []planmodifier.Int64{
											int64planmodifier.UseStateForUnknown(),
										},
									},
									"minute": schema.Int64Attribute{
										Optional:    true,
										Computed:    true,
										Description: "Minute when the full backup should start. If this value is omitted, a random minute will be generated.",
										Validators: []validator.Int64{
											int64validator.Between(0, 59),
										},
										PlanModifiers: []planmodifier.Int64{
											int64planmodifier.UseStateForUnknown(),
										},
									},
								},
								Optional:    true,
								Computed:    true,
								Description: "Schedules for the backup policy.",
							},
						},
						Optional:    true,
						Computed:    true,
						Description: "Scheduled backups policy for the database.",
					},
					"type": schema.StringAttribute{
						Required:    true,
						Description: "Type of the database. Currently only supports 'postgresql'.",
					},
					"version": schema.StringAttribute{
						Required:    true,
						Description: "Minor version of PostgreSQL.",
					},
				},
				Required: true,
			},
			"created_at": schema.StringAttribute{
				CustomType:  timetypes.RFC3339Type{},
				Computed:    true,
				Description: "Date when the database was created.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_by": schema.StringAttribute{
				Computed:    true,
				Description: "Initial creator of the database.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Fulltext description of the database.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(0, 128),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Default: stringdefault.StaticString(""),
			},
			"last_modified_at": schema.StringAttribute{
				CustomType:  timetypes.RFC3339Type{},
				Computed:    true,
				Description: "Date when the database was last modified.",
			},
			"last_modified_by": schema.StringAttribute{
				Computed:    true,
				Description: "User who last changed the database.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the database.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 64),
					stringvalidator.RegexMatches(regexp.MustCompile("^(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?$"), ""),
				},
			},
			"service_config": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"disksize": schema.Int64Attribute{
						Required:    true,
						Description: "Disksize in GB.",
						Validators: []validator.Int64{
							int64validator.Between(5, 500),
						},
					},
					"flavor": schema.StringAttribute{
						Required:    true,
						Description: "VM flavor to use.",
					},
					"maintenance_window": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"day_of_week": schema.Int64Attribute{
								Optional:    true,
								Computed:    true,
								Description: "Day of week as a cron time (0=Sun, 1=Mon, ..., 6=Sat). If omitted, a random day will be used.",
								PlanModifiers: []planmodifier.Int64{
									int64planmodifier.UseStateForUnknown(),
								},
							},
							"start_hour": schema.Int64Attribute{
								Optional:    true,
								Computed:    true,
								Description: "Hour when the maintenance window starts. If omitted, a random hour between 20 and 4 will be used.",
								PlanModifiers: []planmodifier.Int64{
									int64planmodifier.UseStateForUnknown(),
								},
							},
							"start_minute": schema.Int64Attribute{
								Optional:    true,
								Computed:    true,
								Description: "Minute when the maintenance window starts. If omitted, a random minute will be used.",
								PlanModifiers: []planmodifier.Int64{
									int64planmodifier.UseStateForUnknown(),
								},
							},
						},
						Optional:    true,
						Computed:    true,
						Description: "Maintenance window. This will be a time window for updates and maintenance. If omitted, a random window will be generated.",
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},
					},
					"region": schema.StringAttribute{
						Required:    true,
						Description: "Region for the database.",
					},
					"remote_ips": schema.ListAttribute{
						ElementType: types.StringType,
						Optional:    true,
						Computed:    true,
						Description: "List of IP addresses, that should be allowed to connect to the database.",
						PlanModifiers: []planmodifier.List{
							listplanmodifier.UseStateForUnknown(),
						},
					},
					"type": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Description: "Type of the service you want to create (default `database`)",
						Default:     stringdefault.StaticString("database"),
					},
				},
				Required: true,
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "Overall status of the database.",
			},
			"phase": schema.StringAttribute{
				Computed:    true,
				Description: "Detailed status of the database.",
			},
			"resource_status": schema.StringAttribute{
				Computed:    true,
				Description: "Sync status of the database.",
			},
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "UUID of the database.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func psqlGetResponseToModel(ctx context.Context, db database.PostgreSQLGetResponseV1, model *DatabaseModel) diag.Diagnostics {
	var diags diag.Diagnostics
	var conversionDiags diag.Diagnostics

	ctx = tflog.SetField(ctx, "conversion_read_source_response", db)
	tflog.Debug(ctx, "Converting read api response")

	var d diag.Diagnostics
	remoteIps, d := types.ListValueFrom(ctx, types.StringType, db.ServiceConfig.RemoteIps)
	diags.Append(d...)

	serviceConfig := ServiceConfigModelV1{
		Disksize:          types.Int64PointerValue(db.ServiceConfig.Disksize),
		ServiceConfigType: types.StringValue(db.ServiceConfig.Type),
		Flavor:            types.StringValue(db.ServiceConfig.Flavor),
		Region:            types.StringValue(db.ServiceConfig.Region),
		RemoteIps:         remoteIps,
	}

	if db.ServiceConfig.MaintenanceWindow != nil {
		maintenanceWindow := MaintenanceWindowModel{
			DayOfWeek:   types.Int64PointerValue(db.ServiceConfig.MaintenanceWindow.DayOfWeek),
			StartHour:   types.Int64PointerValue(db.ServiceConfig.MaintenanceWindow.StartHour),
			StartMinute: types.Int64PointerValue(db.ServiceConfig.MaintenanceWindow.StartMinute),
		}
		objectValue, conversionDiags := types.ObjectValueFrom(ctx, maintenanceWindow.AttributeTypes(), maintenanceWindow)
		diags.Append(conversionDiags...)
		serviceConfig.MaintenanceWindow = objectValue
	}

	model.ServiceConfig, conversionDiags = types.ObjectValueFrom(ctx, serviceConfig.AttributeTypes(), serviceConfig)
	diags.Append(conversionDiags...)

	applicationConfig := ApplicationConfigModelV1{
		ApplicationConfigType: types.StringValue(db.ApplicationConfig.Type),
		Instances:             types.Int64PointerValue(db.ApplicationConfig.Instances),
		Version:               types.StringValue(db.ApplicationConfig.Version),
	}

	// Extract password consistently - use plan password for create response
	if !model.ApplicationConfig.IsNull() && !model.ApplicationConfig.IsUnknown() {
		var planApplicationConfig ApplicationConfigModelV1
		planDiags := model.ApplicationConfig.As(ctx, &planApplicationConfig, basetypes.ObjectAsOptions{})
		diags.Append(planDiags...)
		if diags.HasError() {
			return diags
		}

		if !planApplicationConfig.Password.IsNull() && !planApplicationConfig.Password.IsUnknown() {
			passwordValue, passwordDiags := planApplicationConfig.Password.ToStringValue(ctx)
			diags.Append(passwordDiags...)
			if diags.HasError() {
				return diags
			}
			applicationConfig.Password = passwordValue
		} else {
			applicationConfig.Password = types.StringNull()
		}
	}

	if db.ApplicationConfig.ScheduledBackups != nil && db.ApplicationConfig.ScheduledBackups.Schedule != nil {
		scheduledBackups := ScheduledBackupsModel{
			Retention: types.Int64PointerValue(db.ApplicationConfig.ScheduledBackups.Retention),
		}

		schedule := ScheduleModel{
			Hour:   types.Int64PointerValue(db.ApplicationConfig.ScheduledBackups.Schedule.Hour),
			Minute: types.Int64PointerValue(db.ApplicationConfig.ScheduledBackups.Schedule.Minute),
		}
		objectValue, conversionDiags := types.ObjectValueFrom(ctx, schedule.AttributeTypes(), schedule)
		diags.Append(conversionDiags...)
		scheduledBackups.Schedule = objectValue

		objectValue, conversionDiags = types.ObjectValueFrom(ctx, scheduledBackups.AttributeTypes(), scheduledBackups)
		diags.Append(conversionDiags...)
		applicationConfig.ScheduledBackups = objectValue
	} else {
		applicationConfig.ScheduledBackups = types.ObjectNull(ScheduleModel{}.AttributeTypes())
	}

	if db.ApplicationConfig.Recovery != nil {
		recovery := RecoveryModel{
			Exclusive:  types.BoolPointerValue(db.ApplicationConfig.Recovery.Exclusive),
			Source:     types.StringPointerValue(db.ApplicationConfig.Recovery.Source),
			TargetLsn:  types.StringPointerValue(db.ApplicationConfig.Recovery.TargetLsn),
			TargetName: types.StringPointerValue(db.ApplicationConfig.Recovery.TargetName),
			TargetXid:  types.StringPointerValue(db.ApplicationConfig.Recovery.TargetXid),
			TargetTime: types.StringValue(db.ApplicationConfig.Recovery.TargetTime.Format(time.RFC3339)),
		}
		objectValue, conversionDiags := types.ObjectValueFrom(ctx, recovery.AttributeTypes(), recovery)
		diags.Append(conversionDiags...)
		applicationConfig.Recovery = objectValue
	} else {
		applicationConfig.Recovery = types.ObjectNull(RecoveryModel{}.AttributeTypes())
	}

	model.Uuid = types.StringValue(db.Uuid)
	model.Name = types.StringValue(db.Name)
	model.Description = types.StringValue(db.Description)
	model.Status = types.StringValue(db.Status)
	model.Phase = types.StringValue(db.Phase)
	model.ResourceStatus = types.StringValue(db.ResourceStatus)
	model.CreatedBy = types.StringValue(db.CreatedBy)
	model.CreatedAt = timetypes.NewRFC3339TimePointerValue(db.CreatedAt)
	model.LastModifiedBy = types.StringValue(db.LastModifiedBy)
	model.LastModifiedAt = timetypes.NewRFC3339TimePointerValue(db.LastModifiedAt)

	var applicationConfigDiags diag.Diagnostics
	model.ApplicationConfig, applicationConfigDiags = types.ObjectValueFrom(ctx, applicationConfig.AttributeTypes(), applicationConfig)
	diags.Append(applicationConfigDiags...)

	var serviceConfigDiags diag.Diagnostics
	model.ServiceConfig, serviceConfigDiags = types.ObjectValueFrom(ctx, serviceConfig.AttributeTypes(), serviceConfig)
	diags.Append(serviceConfigDiags...)

	ctx = tflog.SetField(ctx, "conversion_read_model", model)
	tflog.Debug(ctx, "Converted api read response to model")

	return diags
}
