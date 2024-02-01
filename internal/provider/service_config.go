package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	sys11dbaassdk "github.com/syseleven/sys11dbaas-sdk"
)

var _ basetypes.ObjectTypable = ServiceConfigType{}

type ServiceConfigType struct {
	basetypes.ObjectType
}

func (t ServiceConfigType) Equal(o attr.Type) bool {
	other, ok := o.(ServiceConfigType)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t ServiceConfigType) String() string {
	return "ServiceConfigType"
}

func (t ServiceConfigType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributes := in.Attributes()

	disksizeAttribute, ok := attributes["disksize"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`disksize is missing from object`)

		return nil, diags
	}

	disksizeVal, ok := disksizeAttribute.(basetypes.Int64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`disksize expected to be basetypes.Int64Value, was: %T`, disksizeAttribute))
	}

	flavorAttribute, ok := attributes["flavor"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`flavor is missing from object`)

		return nil, diags
	}

	flavorVal, ok := flavorAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`flavor expected to be basetypes.StringValue, was: %T`, flavorAttribute))
	}

	maintenanceWindowAttribute, ok := attributes["maintenance_window"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`maintenance_window is missing from object`)

		return nil, diags
	}

	maintenanceWindowVal, ok := maintenanceWindowAttribute.(basetypes.ObjectValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`maintenance_window expected to be basetypes.ObjectValue, was: %T`, maintenanceWindowAttribute))
	}

	regionAttribute, ok := attributes["region"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`region is missing from object`)

		return nil, diags
	}

	regionVal, ok := regionAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`region expected to be basetypes.StringValue, was: %T`, regionAttribute))
	}

	remoteIpsAttribute, ok := attributes["remote_ips"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`remote_ips is missing from object`)

		return nil, diags
	}

	remoteIpsVal, ok := remoteIpsAttribute.(basetypes.ListValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`remote_ips expected to be basetypes.ListValue, was: %T`, remoteIpsAttribute))
	}

	typeAttribute, ok := attributes["type"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`type is missing from object`)

		return nil, diags
	}

	typeVal, ok := typeAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`type expected to be basetypes.StringValue, was: %T`, typeAttribute))
	}

	if diags.HasError() {
		return nil, diags
	}

	return ServiceConfigValue{
		Disksize:          disksizeVal,
		Flavor:            flavorVal,
		MaintenanceWindow: maintenanceWindowVal,
		Region:            regionVal,
		RemoteIps:         remoteIpsVal,
		ServiceConfigType: typeVal,
		state:             attr.ValueStateKnown,
	}, diags
}

func NewServiceConfigValueNull() ServiceConfigValue {
	return ServiceConfigValue{
		state: attr.ValueStateNull,
	}
}

func NewServiceConfigValueUnknown() ServiceConfigValue {
	return ServiceConfigValue{
		state: attr.ValueStateUnknown,
	}
}

func NewServiceConfigValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (ServiceConfigValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]

		if !ok {
			diags.AddError(
				"Missing ServiceConfigValue Attribute Value",
				"While creating a ServiceConfigValue value, a missing attribute value was detected. "+
					"A ServiceConfigValue must contain values for all attributes, even if null or unknown. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("ServiceConfigValue Attribute Name (%s) Expected Type: %s", name, attributeType.String()),
			)

			continue
		}

		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError(
				"Invalid ServiceConfigValue Attribute Type",
				"While creating a ServiceConfigValue value, an invalid attribute value was detected. "+
					"A ServiceConfigValue must use a matching attribute type for the value. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("ServiceConfigValue Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
					fmt.Sprintf("ServiceConfigValue Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
			)
		}
	}

	for name := range attributes {
		_, ok := attributeTypes[name]

		if !ok {
			diags.AddError(
				"Extra ServiceConfigValue Attribute Value",
				"While creating a ServiceConfigValue value, an extra attribute value was detected. "+
					"A ServiceConfigValue must not contain values beyond the expected attribute types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Extra ServiceConfigValue Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewServiceConfigValueUnknown(), diags
	}

	disksizeAttribute, ok := attributes["disksize"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`disksize is missing from object`)

		return NewServiceConfigValueUnknown(), diags
	}

	disksizeVal, ok := disksizeAttribute.(basetypes.Int64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`disksize expected to be basetypes.Int64Value, was: %T`, disksizeAttribute))
	}

	flavorAttribute, ok := attributes["flavor"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`flavor is missing from object`)

		return NewServiceConfigValueUnknown(), diags
	}

	flavorVal, ok := flavorAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`flavor expected to be basetypes.StringValue, was: %T`, flavorAttribute))
	}

	maintenanceWindowAttribute, ok := attributes["maintenance_window"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`maintenance_window is missing from object`)

		return NewServiceConfigValueUnknown(), diags
	}

	maintenanceWindowVal, ok := maintenanceWindowAttribute.(basetypes.ObjectValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`maintenance_window expected to be basetypes.ObjectValue, was: %T`, maintenanceWindowAttribute))
	}

	regionAttribute, ok := attributes["region"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`region is missing from object`)

		return NewServiceConfigValueUnknown(), diags
	}

	regionVal, ok := regionAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`region expected to be basetypes.StringValue, was: %T`, regionAttribute))
	}

	remoteIpsAttribute, ok := attributes["remote_ips"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`remote_ips is missing from object`)

		return NewServiceConfigValueUnknown(), diags
	}

	remoteIpsVal, ok := remoteIpsAttribute.(basetypes.ListValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`remote_ips expected to be basetypes.ListValue, was: %T`, remoteIpsAttribute))
	}

	typeAttribute, ok := attributes["type"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`type is missing from object`)

		return NewServiceConfigValueUnknown(), diags
	}

	typeVal, ok := typeAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`type expected to be basetypes.StringValue, was: %T`, typeAttribute))
	}

	if diags.HasError() {
		return NewServiceConfigValueUnknown(), diags
	}

	return ServiceConfigValue{
		Disksize:          disksizeVal,
		Flavor:            flavorVal,
		MaintenanceWindow: maintenanceWindowVal,
		Region:            regionVal,
		RemoteIps:         remoteIpsVal,
		ServiceConfigType: typeVal,
		state:             attr.ValueStateKnown,
	}, diags
}

func NewServiceConfigValueMust(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) ServiceConfigValue {
	object, diags := NewServiceConfigValue(attributeTypes, attributes)

	if diags.HasError() {
		// This could potentially be added to the diag package.
		diagsStrings := make([]string, 0, len(diags))

		for _, diagnostic := range diags {
			diagsStrings = append(diagsStrings, fmt.Sprintf(
				"%s | %s | %s",
				diagnostic.Severity(),
				diagnostic.Summary(),
				diagnostic.Detail()))
		}

		panic("NewServiceConfigValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}

func (t ServiceConfigType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewServiceConfigValueNull(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewServiceConfigValueUnknown(), nil
	}

	if in.IsNull() {
		return NewServiceConfigValueNull(), nil
	}

	attributes := map[string]attr.Value{}

	val := map[string]tftypes.Value{}

	err := in.As(&val)

	if err != nil {
		return nil, err
	}

	for k, v := range val {
		a, err := t.AttrTypes[k].ValueFromTerraform(ctx, v)

		if err != nil {
			return nil, err
		}

		attributes[k] = a
	}

	return NewServiceConfigValueMust(ServiceConfigValue{}.AttributeTypes(ctx), attributes), nil
}

func (t ServiceConfigType) ValueType(ctx context.Context) attr.Value {
	return ServiceConfigValue{}
}

var _ basetypes.ObjectValuable = ServiceConfigValue{}

type ServiceConfigValue struct {
	Disksize          basetypes.Int64Value  `tfsdk:"disksize"`
	Flavor            basetypes.StringValue `tfsdk:"flavor"`
	MaintenanceWindow basetypes.ObjectValue `tfsdk:"maintenance_window"`
	Region            basetypes.StringValue `tfsdk:"region"`
	RemoteIps         basetypes.ListValue   `tfsdk:"remote_ips"`
	ServiceConfigType basetypes.StringValue `tfsdk:"type"`
	state             attr.ValueState
}

func (v ServiceConfigValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 6)

	var val tftypes.Value
	var err error

	attrTypes["disksize"] = basetypes.Int64Type{}.TerraformType(ctx)
	attrTypes["flavor"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["maintenance_window"] = basetypes.ObjectType{
		AttrTypes: MaintenanceWindowValue{}.AttributeTypes(ctx),
	}.TerraformType(ctx)
	attrTypes["region"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["remote_ips"] = basetypes.ListType{
		ElemType: types.StringType,
	}.TerraformType(ctx)
	attrTypes["type"] = basetypes.StringType{}.TerraformType(ctx)

	objectType := tftypes.Object{AttributeTypes: attrTypes}

	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 6)

		val, err = v.Disksize.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["disksize"] = val

		val, err = v.Flavor.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["flavor"] = val

		val, err = v.MaintenanceWindow.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["maintenance_window"] = val

		val, err = v.Region.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["region"] = val

		val, err = v.RemoteIps.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["remote_ips"] = val

		val, err = v.ServiceConfigType.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["type"] = val

		if err := tftypes.ValidateValue(objectType, vals); err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		return tftypes.NewValue(objectType, vals), nil
	case attr.ValueStateNull:
		return tftypes.NewValue(objectType, nil), nil
	case attr.ValueStateUnknown:
		return tftypes.NewValue(objectType, tftypes.UnknownValue), nil
	default:
		panic(fmt.Sprintf("unhandled Object state in ToTerraformValue: %s", v.state))
	}
}

func (v ServiceConfigValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v ServiceConfigValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v ServiceConfigValue) String() string {
	return "ServiceConfigValue"
}

func (v ServiceConfigValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	var maintenanceWindow basetypes.ObjectValue

	if v.MaintenanceWindow.IsNull() {
		maintenanceWindow = types.ObjectNull(
			MaintenanceWindowValue{}.AttributeTypes(ctx),
		)
	}

	if v.MaintenanceWindow.IsUnknown() {
		maintenanceWindow = types.ObjectUnknown(
			MaintenanceWindowValue{}.AttributeTypes(ctx),
		)
	}

	if !v.MaintenanceWindow.IsNull() && !v.MaintenanceWindow.IsUnknown() {
		maintenanceWindow = types.ObjectValueMust(
			MaintenanceWindowValue{}.AttributeTypes(ctx),
			v.MaintenanceWindow.Attributes(),
		)
	}

	remoteIpsVal, d := types.ListValue(types.StringType, v.RemoteIps.Elements())

	diags.Append(d...)

	if d.HasError() {
		return types.ObjectUnknown(map[string]attr.Type{
			"disksize": basetypes.Int64Type{},
			"flavor":   basetypes.StringType{},
			"maintenance_window": basetypes.ObjectType{
				AttrTypes: MaintenanceWindowValue{}.AttributeTypes(ctx),
			},
			"region": basetypes.StringType{},
			"remote_ips": basetypes.ListType{
				ElemType: types.StringType,
			},
			"type": basetypes.StringType{},
		}), diags
	}

	objVal, diags := types.ObjectValue(
		map[string]attr.Type{
			"disksize": basetypes.Int64Type{},
			"flavor":   basetypes.StringType{},
			"maintenance_window": basetypes.ObjectType{
				AttrTypes: MaintenanceWindowValue{}.AttributeTypes(ctx),
			},
			"region": basetypes.StringType{},
			"remote_ips": basetypes.ListType{
				ElemType: types.StringType,
			},
			"type": basetypes.StringType{},
		},
		map[string]attr.Value{
			"disksize":           v.Disksize,
			"flavor":             v.Flavor,
			"maintenance_window": maintenanceWindow,
			"region":             v.Region,
			"remote_ips":         remoteIpsVal,
			"type":               v.ServiceConfigType,
		})

	return objVal, diags
}

func (v ServiceConfigValue) Equal(o attr.Value) bool {
	other, ok := o.(ServiceConfigValue)

	if !ok {
		return false
	}

	if v.state != other.state {
		return false
	}

	if v.state != attr.ValueStateKnown {
		return true
	}

	if !v.Disksize.Equal(other.Disksize) {
		return false
	}

	if !v.Flavor.Equal(other.Flavor) {
		return false
	}

	if !v.MaintenanceWindow.Equal(other.MaintenanceWindow) {
		return false
	}

	if !v.Region.Equal(other.Region) {
		return false
	}

	if !v.RemoteIps.Equal(other.RemoteIps) {
		return false
	}

	if !v.ServiceConfigType.Equal(other.ServiceConfigType) {
		return false
	}

	return true
}

func (v ServiceConfigValue) Type(ctx context.Context) attr.Type {
	return ServiceConfigType{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v ServiceConfigValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"disksize": basetypes.Int64Type{},
		"flavor":   basetypes.StringType{},
		"maintenance_window": basetypes.ObjectType{
			AttrTypes: MaintenanceWindowValue{}.AttributeTypes(ctx),
		},
		"region": basetypes.StringType{},
		"remote_ips": basetypes.ListType{
			ElemType: types.StringType,
		},
		"type": basetypes.StringType{},
	}
}

var _ basetypes.ObjectTypable = MaintenanceWindowType{}

type MaintenanceWindowType struct {
	basetypes.ObjectType
}

func (t MaintenanceWindowType) Equal(o attr.Type) bool {
	other, ok := o.(MaintenanceWindowType)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t MaintenanceWindowType) String() string {
	return "MaintenanceWindowType"
}

func (t MaintenanceWindowType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributes := in.Attributes()

	dayOfWeekAttribute, ok := attributes["day_of_week"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`day_of_week is missing from object`)

		return nil, diags
	}

	dayOfWeekVal, ok := dayOfWeekAttribute.(basetypes.Int64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`day_of_week expected to be basetypes.Int64Value, was: %T`, dayOfWeekAttribute))
	}

	startHourAttribute, ok := attributes["start_hour"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`start_hour is missing from object`)

		return nil, diags
	}

	startHourVal, ok := startHourAttribute.(basetypes.Int64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`start_hour expected to be basetypes.Int64Value, was: %T`, startHourAttribute))
	}

	startMinuteAttribute, ok := attributes["start_minute"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`start_minute is missing from object`)

		return nil, diags
	}

	startMinuteVal, ok := startMinuteAttribute.(basetypes.Int64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`start_minute expected to be basetypes.Int64Value, was: %T`, startMinuteAttribute))
	}

	if diags.HasError() {
		return nil, diags
	}

	return MaintenanceWindowValue{
		DayOfWeek:   dayOfWeekVal,
		StartHour:   startHourVal,
		StartMinute: startMinuteVal,
		state:       attr.ValueStateKnown,
	}.ToObjectValue(ctx)
}

func NewMaintenanceWindowValueNull() MaintenanceWindowValue {
	return MaintenanceWindowValue{
		state: attr.ValueStateNull,
	}
}

func NewMaintenanceWindowValueUnknown() MaintenanceWindowValue {
	return MaintenanceWindowValue{
		state: attr.ValueStateUnknown,
	}
}

func NewMaintenanceWindowValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (MaintenanceWindowValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]

		if !ok {
			diags.AddError(
				"Missing MaintenanceWindowValue Attribute Value",
				"While creating a MaintenanceWindowValue value, a missing attribute value was detected. "+
					"A MaintenanceWindowValue must contain values for all attributes, even if null or unknown. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("MaintenanceWindowValue Attribute Name (%s) Expected Type: %s, attributes: %v", name, attributeType.String(), attributes),
			)

			continue
		}

		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError(
				"Invalid MaintenanceWindowValue Attribute Type",
				"While creating a MaintenanceWindowValue value, an invalid attribute value was detected. "+
					"A MaintenanceWindowValue must use a matching attribute type for the value. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("MaintenanceWindowValue Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
					fmt.Sprintf("MaintenanceWindowValue Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
			)
		}
	}

	for name := range attributes {
		_, ok := attributeTypes[name]

		if !ok {
			diags.AddError(
				"Extra MaintenanceWindowValue Attribute Value",
				"While creating a MaintenanceWindowValue value, an extra attribute value was detected. "+
					"A MaintenanceWindowValue must not contain values beyond the expected attribute types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Extra MaintenanceWindowValue Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewMaintenanceWindowValueUnknown(), diags
	}

	dayOfWeekAttribute, ok := attributes["day_of_week"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`day_of_week is missing from object`)

		return NewMaintenanceWindowValueUnknown(), diags
	}

	dayOfWeekVal, ok := dayOfWeekAttribute.(basetypes.Int64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`day_of_week expected to be basetypes.Int64Value, was: %T`, dayOfWeekAttribute))
	}

	startHourAttribute, ok := attributes["start_hour"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`start_hour is missing from object`)

		return NewMaintenanceWindowValueUnknown(), diags
	}

	startHourVal, ok := startHourAttribute.(basetypes.Int64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`start_hour expected to be basetypes.Int64Value, was: %T`, startHourAttribute))
	}

	startMinuteAttribute, ok := attributes["start_minute"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`start_minute is missing from object`)

		return NewMaintenanceWindowValueUnknown(), diags
	}

	startMinuteVal, ok := startMinuteAttribute.(basetypes.Int64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`start_minute expected to be basetypes.Int64Value, was: %T`, startMinuteAttribute))
	}

	if diags.HasError() {
		return NewMaintenanceWindowValueUnknown(), diags
	}

	return MaintenanceWindowValue{
		DayOfWeek:   dayOfWeekVal,
		StartHour:   startHourVal,
		StartMinute: startMinuteVal,
		state:       attr.ValueStateKnown,
	}, diags
}

func NewMaintenanceWindowValueMust(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) MaintenanceWindowValue {
	object, diags := NewMaintenanceWindowValue(attributeTypes, attributes)

	if diags.HasError() {
		// This could potentially be added to the diag package.
		diagsStrings := make([]string, 0, len(diags))

		for _, diagnostic := range diags {
			diagsStrings = append(diagsStrings, fmt.Sprintf(
				"%s | %s | %s",
				diagnostic.Severity(),
				diagnostic.Summary(),
				diagnostic.Detail()))
		}

		panic("NewMaintenanceWindowValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}

func (t MaintenanceWindowType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewMaintenanceWindowValueNull(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewMaintenanceWindowValueUnknown(), nil
	}

	if in.IsNull() {
		return NewMaintenanceWindowValueNull(), nil
	}

	attributes := map[string]attr.Value{}

	val := map[string]tftypes.Value{}

	err := in.As(&val)

	if err != nil {
		return nil, err
	}

	for k, v := range val {
		a, err := t.AttrTypes[k].ValueFromTerraform(ctx, v)

		if err != nil {
			return nil, err
		}

		attributes[k] = a
	}

	return NewMaintenanceWindowValueMust(MaintenanceWindowValue{}.AttributeTypes(ctx), attributes), nil
}

func (t MaintenanceWindowType) ValueType(ctx context.Context) attr.Value {
	return MaintenanceWindowValue{}
}

var _ basetypes.ObjectValuable = MaintenanceWindowValue{}

type MaintenanceWindowValue struct {
	DayOfWeek   basetypes.Int64Value `tfsdk:"day_of_week"`
	StartHour   basetypes.Int64Value `tfsdk:"start_hour"`
	StartMinute basetypes.Int64Value `tfsdk:"start_minute"`
	state       attr.ValueState
}

func (v MaintenanceWindowValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 3)

	var val tftypes.Value
	var err error

	attrTypes["day_of_week"] = basetypes.Int64Type{}.TerraformType(ctx)
	attrTypes["start_hour"] = basetypes.Int64Type{}.TerraformType(ctx)
	attrTypes["start_minute"] = basetypes.Int64Type{}.TerraformType(ctx)

	objectType := tftypes.Object{AttributeTypes: attrTypes}

	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 3)

		val, err = v.DayOfWeek.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["day_of_week"] = val

		val, err = v.StartHour.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["start_hour"] = val

		val, err = v.StartMinute.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["start_minute"] = val

		if err := tftypes.ValidateValue(objectType, vals); err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		return tftypes.NewValue(objectType, vals), nil
	case attr.ValueStateNull:
		return tftypes.NewValue(objectType, nil), nil
	case attr.ValueStateUnknown:
		return tftypes.NewValue(objectType, tftypes.UnknownValue), nil
	default:
		panic(fmt.Sprintf("unhandled Object state in ToTerraformValue: %s", v.state))
	}
}

func (v MaintenanceWindowValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v MaintenanceWindowValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v MaintenanceWindowValue) String() string {
	return "MaintenanceWindowValue"
}

func (v MaintenanceWindowValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	objVal, diags := types.ObjectValue(
		map[string]attr.Type{
			"day_of_week":  basetypes.Int64Type{},
			"start_hour":   basetypes.Int64Type{},
			"start_minute": basetypes.Int64Type{},
		},
		map[string]attr.Value{
			"day_of_week":  v.DayOfWeek,
			"start_hour":   v.StartHour,
			"start_minute": v.StartMinute,
		})

	return objVal, diags
}

func (v MaintenanceWindowValue) Equal(o attr.Value) bool {
	other, ok := o.(MaintenanceWindowValue)

	if !ok {
		return false
	}

	if v.state != other.state {
		return false
	}

	if v.state != attr.ValueStateKnown {
		return true
	}

	if !v.DayOfWeek.Equal(other.DayOfWeek) {
		return false
	}

	if !v.StartHour.Equal(other.StartHour) {
		return false
	}

	if !v.StartMinute.Equal(other.StartMinute) {
		return false
	}

	return true
}

func (v MaintenanceWindowValue) Type(ctx context.Context) attr.Type {
	return MaintenanceWindowType{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v MaintenanceWindowValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"day_of_week":  basetypes.Int64Type{},
		"start_hour":   basetypes.Int64Type{},
		"start_minute": basetypes.Int64Type{},
	}
}

func (v MaintenanceWindowValue) ToDBaaSSdkObject(ctx context.Context) (*sys11dbaassdk.MaintenanceWindow, diag.Diagnostics) {

	var dayOfWeek *int
	dayOfWeek = nil
	if !v.DayOfWeek.IsNull() && !v.DayOfWeek.IsUnknown() {
		dayOfWeek = sys11dbaassdk.Int64ToIntPtr(v.DayOfWeek.ValueInt64())
	}

	var startHour *int
	startHour = nil
	if !v.StartHour.IsNull() && !v.StartHour.IsUnknown() {
		startHour = sys11dbaassdk.Int64ToIntPtr(v.StartHour.ValueInt64())
	}

	var startMinute *int
	startMinute = nil
	if !v.StartMinute.IsNull() && !v.StartMinute.IsUnknown() {
		startMinute = sys11dbaassdk.Int64ToIntPtr(v.StartMinute.ValueInt64())
	}

	return &sys11dbaassdk.MaintenanceWindow{
		DayOfWeek:   dayOfWeek,
		StartHour:   startHour,
		StartMinute: startMinute,
	}, diag.Diagnostics{}
}
