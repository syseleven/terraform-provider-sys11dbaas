package provider

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
)

func DatabaseResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"application_config": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"hostname": schema.StringAttribute{
						Computed:            true,
						Description:         "DNS name of the database in the format uuid.postgresql.syseleven.services.",
						MarkdownDescription: "DNS name of the database in the format uuid.postgresql.syseleven.services.",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"instances": schema.Int64Attribute{
						Required:            true,
						Description:         "Node count of the database cluster.",
						MarkdownDescription: "Node count of the database cluster.",
						Validators: []validator.Int64{
							int64validator.AtMost(5),
						},
					},
					"ip_address": schema.StringAttribute{
						Computed:            true,
						Description:         "Public IP address of the database. It will be 'Pending' if no address has been assigned yet.",
						MarkdownDescription: "Public IP address of the database. It will be 'Pending' if no address has been assigned yet.",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"password": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						Sensitive:           true,
						Description:         "Password for the admin user.",
						MarkdownDescription: "Password for the admin user.",
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
								Optional: true,
								Computed: true,
						    Description:         "Set to true, when the given target should be excluded.",
						    MarkdownDescription: "Set to true, when the given target should be excluded.",
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
								},
							},
							"source": schema.StringAttribute{
								Optional: true,
								Computed: true,
						    Description:         "UUID of the source database.",
						    MarkdownDescription: "UUID of the source database.",
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
							"target_lsn": schema.StringAttribute{
								Optional: true,
								Computed: true,
						    Description:         "LSN of the write-ahead log location up to which recovery will proceed. target_* parameters are mutually exclusive.",
						    MarkdownDescription: "LSN of the write-ahead log location up to which recovery will proceed. target_* parameters are mutually exclusive.",
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
							"target_name": schema.StringAttribute{
								Optional: true,
								Computed: true,
						    Description:         "Named restore point (created with pg_create_restore_point()) to which recovery will proceed. target_* parameters are mutually exclusive.",
						    MarkdownDescription: "Named restore point (created with pg_create_restore_point()) to which recovery will proceed. target_* parameters are mutually exclusive.",
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
							"target_time": schema.StringAttribute{
								Optional: true,
								Computed: true,
						    Description:         "Time stamp up to which recovery will proceed, expressed in RFC 3339 format. target_* parameters are mutually exclusive.",
						    MarkdownDescription: "Time stamp up to which recovery will proceed, expressed in RFC 3339 format. target_* parameters are mutually exclusive.",
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
							"target_xid": schema.StringAttribute{
								Optional: true,
								Computed: true,
						    Description:         "Transaction ID up to which recovery will proceed. target_* parameters are mutually exclusive.",
						    MarkdownDescription: "Transaction ID up to which recovery will proceed. target_* parameters are mutually exclusive.",
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
								Optional:            true,
								Computed:            true,
								Description:         "Duration in days for which backups should be stored.",
								MarkdownDescription: "Duration in days for which backups should be stored.",
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
										Optional:            true,
										Computed:            true,
										Description:         "Hour when the full backup should start. If this value is omitted, a random hour between 1am and 5am will be generated.",
										MarkdownDescription: "Hour when the full backup should start. If this value is omitted, a random hour between 1am and 5am will be generated.",
										Validators: []validator.Int64{
											int64validator.Between(0, 23),
										},
										PlanModifiers: []planmodifier.Int64{
											int64planmodifier.UseStateForUnknown(),
										},
									},
									"minute": schema.Int64Attribute{
										Optional:            true,
										Computed:            true,
										Description:         "Minute when the full backup should start. If this value is omitted, a random minute will be generated.",
										MarkdownDescription: "Minute when the full backup should start. If this value is omitted, a random minute will be generated.",
										Validators: []validator.Int64{
											int64validator.Between(0, 59),
										},
										PlanModifiers: []planmodifier.Int64{
											int64planmodifier.UseStateForUnknown(),
										},
									},
								},
								Optional:            true,
								Computed:            true,
								Description:         "Schedules for the backup policy.",
								MarkdownDescription: "Schedules for the backup policy.",
							},
						},
						Optional:            true,
						Computed:            true,
						Description:         "Scheduled backups policy for the database.",
						MarkdownDescription: "Scheduled backups policy for the database.",
					},
					"type": schema.StringAttribute{
						Required: true,
						Description:         "Type of the database. Currently only supports 'postgresql'.",
						MarkdownDescription: "Type of the database. Currently only supports 'postgresql'.",
					},
					"version": schema.StringAttribute{
						Required:            true,
						Description:         "Minor version of PostgreSQL.",
						MarkdownDescription: "Minor version of PostgreSQL.",
					},
				},
				Required: true,
			},
			"created_at": schema.StringAttribute{
				Computed:            true,
				Description:         "Date when the database was created.",
				MarkdownDescription: "Date when the database was created.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_by": schema.StringAttribute{
				Computed:            true,
				Description:         "Initial creator of the database.",
				MarkdownDescription: "Initial creator of the database.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Fulltext description of the database.",
				MarkdownDescription: "Fulltext description of the database.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(0, 128),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Default: stringdefault.StaticString(""),
			},
			"last_modified_at": schema.StringAttribute{
				Computed:            true,
				Description:         "Date when the database was last modified.",
				MarkdownDescription: "Date when the database was last modified.",
			},
			"last_modified_by": schema.StringAttribute{
				Computed:            true,
				Description:         "User who last changed the database.",
				MarkdownDescription: "User who last changed the database.",
			},
			"name": schema.StringAttribute{
				Required:            true,
				Description:         "Name of the database.",
				MarkdownDescription: "Name of the database.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 64),
					stringvalidator.RegexMatches(regexp.MustCompile("^(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?$"), ""),
				},
			},
			"service_config": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"disksize": schema.Int64Attribute{
						Required:            true,
						Description:         "Disksize in GB.",
						MarkdownDescription: "Disksize in GB.",
						Validators: []validator.Int64{
							int64validator.Between(5, 500),
						},
					},
					"flavor": schema.StringAttribute{
						Required:            true,
						Description:         "VM flavor to use.",
						MarkdownDescription: "VM flavor to use.",
					},
					"maintenance_window": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"day_of_week": schema.Int64Attribute{
								Optional:            true,
								Computed:            true,
								Description:         "Day of week as a cron time (0=Sun, 1=Mon, ..., 6=Sat). If omitted, a random day will be used.",
								MarkdownDescription: "Day of week as a cron time (0=Sun, 1=Mon, ..., 6=Sat). If omitted, a random day will be used.",
								PlanModifiers: []planmodifier.Int64{
									int64planmodifier.UseStateForUnknown(),
								},
							},
							"start_hour": schema.Int64Attribute{
								Optional:            true,
								Computed:            true,
								Description:         "Hour when the maintenance window starts. If omitted, a random hour between 20 and 4 will be used.",
								MarkdownDescription: "Hour when the maintenance window starts. If omitted, a random hour between 20 and 4 will be used.",
								PlanModifiers: []planmodifier.Int64{
									int64planmodifier.UseStateForUnknown(),
								},
							},
							"start_minute": schema.Int64Attribute{
								Optional:            true,
								Computed:            true,
								Description:         "Minute when the maintenance window starts. If omitted, a random minute will be used.",
								MarkdownDescription: "Minute when the maintenance window starts. If omitted, a random minute will be used.",
								PlanModifiers: []planmodifier.Int64{
									int64planmodifier.UseStateForUnknown(),
								},
							},
						},
						Optional:            true,
						Computed:            true,
						Description:         "Maintenance window. This will be a time window for updates and maintenance. If omitted, a random window will be generated.",
						MarkdownDescription: "Maintenance window. This will be a time window for updates and maintenance. If omitted, a random window will be generated.",
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},
					},
					"region": schema.StringAttribute{
						Required:            true,
						Description:         "Region for the database.",
						MarkdownDescription: "Region for the database.",
					},
					"remote_ips": schema.ListAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						Computed:            true,
						Description:         "List of IP addresses, that should be allowed to connect to the database.",
						MarkdownDescription: "List of IP addresses, that should be allowed to connect to the database.",
						PlanModifiers: []planmodifier.List{
							listplanmodifier.UseStateForUnknown(),
						},
					},
					"type": schema.StringAttribute{
						Required: true,
					},
				},
				Required: true,
			},
			"status": schema.StringAttribute{
				Computed: true,
				Description:         "Overall status of the database.",
				MarkdownDescription: "Overall status of the database.",
			},
			"phase": schema.StringAttribute{
				Computed: true,
				Description:         "Detailed status of the database.",
				MarkdownDescription: "Detailed status of the database.",
			},
			"resource_status": schema.StringAttribute{
				Computed: true,
				Description:         "Sync status of the database.",
				MarkdownDescription: "Sync status of the database.",
			},
			"uuid": schema.StringAttribute{
				Computed:            true,
				Description:         "UUID of the database.",
				MarkdownDescription: "UUID of the database.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}
