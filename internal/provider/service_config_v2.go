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

var _ basetypes.ObjectTypable = ServiceConfigTypeV2{}

type ServiceConfigTypeV2 struct {
	basetypes.ObjectType
}

func (t ServiceConfigTypeV2) Equal(o attr.Type) bool {
	other, ok := o.(ServiceConfigType)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t ServiceConfigTypeV2) String() string {
	return "ServiceConfigTypeV2"
}

func (t ServiceConfigTypeV2) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
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

	return ServiceConfigValueV2{
		Disksize:          disksizeVal,
		Flavor:            flavorVal,
		MaintenanceWindow: maintenanceWindowVal,
		Region:            regionVal,
		ServiceConfigType: typeVal,
		state:             attr.ValueStateKnown,
	}, diags
}

func NewServiceConfigValueV2Null() ServiceConfigValueV2 {
	return ServiceConfigValueV2{
		state: attr.ValueStateNull,
	}
}

func NewServiceConfigValueV2Unknown() ServiceConfigValueV2 {
	return ServiceConfigValueV2{
		state: attr.ValueStateUnknown,
	}
}

func NewServiceConfigValueV2(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (ServiceConfigValueV2, diag.Diagnostics) {
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
		return NewServiceConfigValueV2Unknown(), diags
	}

	disksizeAttribute, ok := attributes["disksize"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`disksize is missing from object`)

		return NewServiceConfigValueV2Unknown(), diags
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

		return NewServiceConfigValueV2Unknown(), diags
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

		return NewServiceConfigValueV2Unknown(), diags
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

		return NewServiceConfigValueV2Unknown(), diags
	}

	regionVal, ok := regionAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`region expected to be basetypes.StringValue, was: %T`, regionAttribute))
	}

	typeAttribute, ok := attributes["type"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`type is missing from object`)

		return NewServiceConfigValueV2Unknown(), diags
	}

	typeVal, ok := typeAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`type expected to be basetypes.StringValue, was: %T`, typeAttribute))
	}

	if diags.HasError() {
		return NewServiceConfigValueV2Unknown(), diags
	}

	return ServiceConfigValueV2{
		Disksize:          disksizeVal,
		Flavor:            flavorVal,
		MaintenanceWindow: maintenanceWindowVal,
		Region:            regionVal,
		ServiceConfigType: typeVal,
		state:             attr.ValueStateKnown,
	}, diags
}

func NewServiceConfigValueV2Must(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) ServiceConfigValueV2 {
	object, diags := NewServiceConfigValueV2(attributeTypes, attributes)

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

func (t ServiceConfigTypeV2) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewServiceConfigValueV2Null(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewServiceConfigValueV2Unknown(), nil
	}

	if in.IsNull() {
		return NewServiceConfigValueV2Null(), nil
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

	return NewServiceConfigValueV2Must(ServiceConfigValueV2{}.AttributeTypes(ctx), attributes), nil
}

func (t ServiceConfigTypeV2) ValueType(ctx context.Context) attr.Value {
	return ServiceConfigValueV2{}
}

var _ basetypes.ObjectValuable = ServiceConfigValueV2{}

type ServiceConfigValueV2 struct {
	Disksize          basetypes.Int64Value  `tfsdk:"disksize"`
	Flavor            basetypes.StringValue `tfsdk:"flavor"`
	MaintenanceWindow basetypes.ObjectValue `tfsdk:"maintenance_window"`
	Region            basetypes.StringValue `tfsdk:"region"`
	ServiceConfigType basetypes.StringValue `tfsdk:"type"`
	state             attr.ValueState
}

func (v ServiceConfigValueV2) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 6)

	var val tftypes.Value
	var err error

	attrTypes["disksize"] = basetypes.Int64Type{}.TerraformType(ctx)
	attrTypes["flavor"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["maintenance_window"] = basetypes.ObjectType{
		AttrTypes: MaintenanceWindowValue{}.AttributeTypes(ctx),
	}.TerraformType(ctx)
	attrTypes["region"] = basetypes.StringType{}.TerraformType(ctx)
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

func (v ServiceConfigValueV2) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v ServiceConfigValueV2) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v ServiceConfigValueV2) String() string {
	return "ServiceConfigValueV2"
}

func (v ServiceConfigValueV2) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
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

	if diags.HasError() {
		return types.ObjectUnknown(map[string]attr.Type{
			"disksize": basetypes.Int64Type{},
			"flavor":   basetypes.StringType{},
			"maintenance_window": basetypes.ObjectType{
				AttrTypes: MaintenanceWindowValue{}.AttributeTypes(ctx),
			},
			"region": basetypes.StringType{},
			"type":   basetypes.StringType{},
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
			"type":   basetypes.StringType{},
		},
		map[string]attr.Value{
			"disksize":           v.Disksize,
			"flavor":             v.Flavor,
			"maintenance_window": maintenanceWindow,
			"region":             v.Region,
			"type":               v.ServiceConfigType,
		})

	return objVal, diags
}

func (v ServiceConfigValueV2) Equal(o attr.Value) bool {
	other, ok := o.(ServiceConfigValueV2)

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

	if !v.ServiceConfigType.Equal(other.ServiceConfigType) {
		return false
	}

	return true
}

func (v ServiceConfigValueV2) Type(ctx context.Context) attr.Type {
	return ServiceConfigTypeV2{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v ServiceConfigValueV2) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"disksize": basetypes.Int64Type{},
		"flavor":   basetypes.StringType{},
		"maintenance_window": basetypes.ObjectType{
			AttrTypes: MaintenanceWindowValue{}.AttributeTypes(ctx),
		},
		"region": basetypes.StringType{},
		"type":   basetypes.StringType{},
	}
}

var _ basetypes.ObjectTypable = MaintenanceWindowTypeV2{}

type MaintenanceWindowTypeV2 struct {
	basetypes.ObjectType
}

func (t MaintenanceWindowTypeV2) Equal(o attr.Type) bool {
	other, ok := o.(MaintenanceWindowTypeV2)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t MaintenanceWindowTypeV2) String() string {
	return "MaintenanceWindowTypeV2"
}

func (t MaintenanceWindowTypeV2) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
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

	return MaintenanceWindowValueV2{
		DayOfWeek:   dayOfWeekVal,
		StartHour:   startHourVal,
		StartMinute: startMinuteVal,
		state:       attr.ValueStateKnown,
	}.ToObjectValue(ctx)
}

func NewMaintenanceWindowValueV2Null() MaintenanceWindowValueV2 {
	return MaintenanceWindowValueV2{
		state: attr.ValueStateNull,
	}
}

func NewMaintenanceWindowValueV2Unknown() MaintenanceWindowValueV2 {
	return MaintenanceWindowValueV2{
		state: attr.ValueStateUnknown,
	}
}

func NewMaintenanceWindowValueV2(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (MaintenanceWindowValueV2, diag.Diagnostics) {
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
		return NewMaintenanceWindowValueV2Unknown(), diags
	}

	dayOfWeekAttribute, ok := attributes["day_of_week"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`day_of_week is missing from object`)

		return NewMaintenanceWindowValueV2Unknown(), diags
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

		return NewMaintenanceWindowValueV2Unknown(), diags
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

		return NewMaintenanceWindowValueV2Unknown(), diags
	}

	startMinuteVal, ok := startMinuteAttribute.(basetypes.Int64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`start_minute expected to be basetypes.Int64Value, was: %T`, startMinuteAttribute))
	}

	if diags.HasError() {
		return NewMaintenanceWindowValueV2Unknown(), diags
	}

	return MaintenanceWindowValueV2{
		DayOfWeek:   dayOfWeekVal,
		StartHour:   startHourVal,
		StartMinute: startMinuteVal,
		state:       attr.ValueStateKnown,
	}, diags
}

func NewMaintenanceWindowValueV2Must(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) MaintenanceWindowValueV2 {
	object, diags := NewMaintenanceWindowValueV2(attributeTypes, attributes)

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

func (t MaintenanceWindowTypeV2) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewMaintenanceWindowValueV2Null(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewMaintenanceWindowValueV2Unknown(), nil
	}

	if in.IsNull() {
		return NewMaintenanceWindowValueV2Null(), nil
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

	return NewMaintenanceWindowValueV2Must(MaintenanceWindowValueV2{}.AttributeTypes(ctx), attributes), nil
}

func (t MaintenanceWindowTypeV2) ValueType(ctx context.Context) attr.Value {
	return MaintenanceWindowValueV2{}
}

var _ basetypes.ObjectValuable = MaintenanceWindowValueV2{}

type MaintenanceWindowValueV2 struct {
	DayOfWeek   basetypes.Int64Value `tfsdk:"day_of_week"`
	StartHour   basetypes.Int64Value `tfsdk:"start_hour"`
	StartMinute basetypes.Int64Value `tfsdk:"start_minute"`
	state       attr.ValueState
}

func (v MaintenanceWindowValueV2) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
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

func (v MaintenanceWindowValueV2) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v MaintenanceWindowValueV2) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v MaintenanceWindowValueV2) String() string {
	return "MaintenanceWindowValueV2"
}

func (v MaintenanceWindowValueV2) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
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

func (v MaintenanceWindowValueV2) Equal(o attr.Value) bool {
	other, ok := o.(MaintenanceWindowValueV2)

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

func (v MaintenanceWindowValueV2) Type(ctx context.Context) attr.Type {
	return MaintenanceWindowTypeV2{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v MaintenanceWindowValueV2) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"day_of_week":  basetypes.Int64Type{},
		"start_hour":   basetypes.Int64Type{},
		"start_minute": basetypes.Int64Type{},
	}
}

func (v MaintenanceWindowValueV2) ToDBaaSSdkObject(ctx context.Context) (*sys11dbaassdk.MaintenanceWindowV2, diag.Diagnostics) {

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

	return &sys11dbaassdk.MaintenanceWindowV2{
		DayOfWeek:   dayOfWeek,
		StartHour:   startHour,
		StartMinute: startMinute,
	}, diag.Diagnostics{}
}
