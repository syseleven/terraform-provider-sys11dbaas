package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type MaintenanceWindowModel struct {
	DayOfWeek   types.Int64 `tfsdk:"day_of_week"`
	StartHour   types.Int64 `tfsdk:"start_hour"`
	StartMinute types.Int64 `tfsdk:"start_minute"`
}

func (m MaintenanceWindowModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"day_of_week":  types.Int64Type,
		"start_hour":   types.Int64Type,
		"start_minute": types.Int64Type,
	}
}

type ScheduleModel struct {
	Hour   types.Int64 `tfsdk:"hour"`
	Minute types.Int64 `tfsdk:"minute"`
}

func (m ScheduleModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"hour":   types.Int64Type,
		"minute": types.Int64Type,
	}
}

type ScheduledBackupsModel struct {
	Retention types.Int64  `tfsdk:"retention"`
	Schedule  types.Object `tfsdk:"schedule"`
}

func (m ScheduledBackupsModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"retention": types.Int64Type,
		"schedule": types.ObjectType{
			AttrTypes: ScheduleModel{}.AttributeTypes(),
		},
	}
}

type RecoveryModel struct {
	Exclusive  types.Bool   `tfsdk:"exclusive"`
	Source     types.String `tfsdk:"source"`
	TargetLsn  types.String `tfsdk:"target_lsn"`
	TargetName types.String `tfsdk:"target_name"`
	TargetTime types.String `tfsdk:"target_time"`
	TargetXid  types.String `tfsdk:"target_xid"`
}

func (m RecoveryModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"exclusive":   types.BoolType,
		"source":      types.StringType,
		"target_lsn":  types.StringType,
		"target_name": types.StringType,
		"target_time": types.StringType,
		"target_xid":  types.StringType,
	}
}
