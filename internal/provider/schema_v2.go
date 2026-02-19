package provider

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
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
)

func DatabaseResourceV2Schema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"application_config": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"instances": schema.Int64Attribute{
						Required:    true,
						Description: "Node count of the database cluster.",
						Validators: []validator.Int64{
							int64validator.AtMost(5),
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
								PlanModifiers: []planmodifier.Object{
									objectplanmodifier.UseStateForUnknown(),
								},
							},
						},
						Optional:    true,
						Computed:    true,
						Description: "Scheduled backups policy for the database.",
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},
					},
					"type": schema.StringAttribute{
						Required:    true,
						Description: "Type of the database. Currently only supports 'postgresql'.",
					},
					"version": schema.StringAttribute{
						Required:    true,
						Description: "Minor version of PostgreSQL.",
					},
					"private_networking": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Optional:    true,
								Computed:    true,
								Description: "Set to true, when private networking should be enabled.",
								Default:     booldefault.StaticBool(true),
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
								},
							},
							"allowed_cidrs": schema.ListAttribute{
								ElementType: types.StringType,
								Optional:    true,
								Computed:    true,
								Description: "List of IP addresses, that should be allowed to connect to the database via private networking.",
								PlanModifiers: []planmodifier.List{
									listplanmodifier.UseStateForUnknown(),
								},
							},
							"hostname": schema.StringAttribute{
								Computed:    true,
								Description: "DNS name of the database in the format uuid.postgresql-private.syseleven.services.",
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
							"ip_address": schema.StringAttribute{
								Computed:    true,
								Description: "Private IP address of the database. It will be 'pending' if no address has been assigned yet.",
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
							"shared_subnet_cidr": schema.StringAttribute{
								Optional:    true,
								Computed:    true,
								Default:     stringdefault.StaticString("10.240.0.0/24"),
								Description: "The subnet cidr for the shared network. Make sure this does not collide with other subnets you already use in your project.",
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
							"shared_subnet_id": schema.StringAttribute{
								Computed:    true,
								Description: "Openstack ID of the shared subnet.",
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
							"shared_network_id": schema.StringAttribute{
								Computed:    true,
								Description: "Openstack ID of the shared network.",
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
						},
						Optional: true,
						Computed: true,
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},
					},
					"public_networking": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Optional:    true,
								Computed:    true,
								Description: "Set to true, when public networking should be enabled.",
								Default:     booldefault.StaticBool(false),
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
								},
							},
							"allowed_cidrs": schema.ListAttribute{
								ElementType: types.StringType,
								Optional:    true,
								Computed:    true,
								Description: "List of IP addresses, that should be allowed to connect to the database via public networking.",
								PlanModifiers: []planmodifier.List{
									listplanmodifier.UseStateForUnknown(),
								},
							},
							"hostname": schema.StringAttribute{
								Computed:    true,
								Description: "DNS name of the database in the format uuid.postgresql.syseleven.services.",
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
							"ip_address": schema.StringAttribute{
								Computed:    true,
								Description: "Public IP address of the database. It will be 'pending' if no address has been assigned yet.",
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
						},
						Optional: true,
						Computed: true,
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},
					},
				},
				Required: true,
			},
			"created_at": schema.StringAttribute{
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
				Computed:    true,
				Description: "Date when the database was last modified.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_modified_by": schema.StringAttribute{
				Computed:    true,
				Description: "User who last changed the database.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"phase": schema.StringAttribute{
				Computed:    true,
				Description: "Detailed status of the database.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"resource_status": schema.StringAttribute{
				Computed:    true,
				Description: "Sync status of the database.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
