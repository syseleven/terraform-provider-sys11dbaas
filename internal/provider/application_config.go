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

var _ basetypes.ObjectTypable = ApplicationConfigType{}

type ApplicationConfigType struct {
	basetypes.ObjectType
}

func (t ApplicationConfigType) Equal(o attr.Type) bool {
	other, ok := o.(ApplicationConfigType)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t ApplicationConfigType) String() string {
	return "ApplicationConfigType"
}

func (t ApplicationConfigType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributes := in.Attributes()

	hostnameAttribute, ok := attributes["hostname"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`hostname is missing from object`)

		return nil, diags
	}

	hostnameVal, ok := hostnameAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`hostname expected to be basetypes.StringValue, was: %T`, hostnameAttribute))
	}

	instancesAttribute, ok := attributes["instances"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`instances is missing from object`)

		return nil, diags
	}

	instancesVal, ok := instancesAttribute.(basetypes.Int64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`instances expected to be basetypes.Int64Value, was: %T`, instancesAttribute))
	}

	ipAddressAttribute, ok := attributes["ip_address"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`ip_address is missing from object`)

		return nil, diags
	}

	ipAddressVal, ok := ipAddressAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`ip_address expected to be basetypes.StringValue, was: %T`, ipAddressAttribute))
	}

	passwordAttribute, ok := attributes["password"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`password is missing from object`)

		return nil, diags
	}

	passwordVal, ok := passwordAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`password expected to be basetypes.StringValue, was: %T`, passwordAttribute))
	}

	recoveryAttribute, ok := attributes["recovery"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`recovery is missing from object`)

		return nil, diags
	}

	recoveryVal, ok := recoveryAttribute.(basetypes.ObjectValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`recovery expected to be basetypes.ObjectValue, was: %T`, recoveryAttribute))
	}

	scheduledBackupsAttribute, ok := attributes["scheduled_backups"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`scheduled_backups is missing from object`)

		return nil, diags
	}

	scheduledBackupsVal, ok := scheduledBackupsAttribute.(basetypes.ObjectValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`scheduled_backups expected to be basetypes.ObjectValue, was: %T`, scheduledBackupsAttribute))
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

	versionAttribute, ok := attributes["version"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`version is missing from object`)

		return nil, diags
	}

	versionVal, ok := versionAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`version expected to be basetypes.StringValue, was: %T`, versionAttribute))
	}

	if diags.HasError() {
		return nil, diags
	}

	return ApplicationConfigValue{
		Hostname:              hostnameVal,
		Instances:             instancesVal,
		IpAddress:             ipAddressVal,
		Password:              passwordVal,
		Recovery:              recoveryVal,
		ScheduledBackups:      scheduledBackupsVal,
		ApplicationConfigType: typeVal,
		Version:               versionVal,
		state:                 attr.ValueStateKnown,
	}, diags
}

func NewApplicationConfigValueNull() ApplicationConfigValue {
	return ApplicationConfigValue{
		state: attr.ValueStateNull,
	}
}

func NewApplicationConfigValueUnknown() ApplicationConfigValue {
	return ApplicationConfigValue{
		state: attr.ValueStateUnknown,
	}
}

func NewApplicationConfigValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (ApplicationConfigValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]

		if !ok {
			diags.AddError(
				"Missing ApplicationConfigValue Attribute Value",
				"While creating a ApplicationConfigValue value, a missing attribute value was detected. "+
					"A ApplicationConfigValue must contain values for all attributes, even if null or unknown. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("ApplicationConfigValue Attribute Name (%s) Expected Type: %s", name, attributeType.String()),
			)

			continue
		}

		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError(
				"Invalid ApplicationConfigValue Attribute Type",
				"While creating a ApplicationConfigValue value, an invalid attribute value was detected. "+
					"A ApplicationConfigValue must use a matching attribute type for the value. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("ApplicationConfigValue Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
					fmt.Sprintf("ApplicationConfigValue Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
			)
		}
	}

	for name := range attributes {
		_, ok := attributeTypes[name]

		if !ok {
			diags.AddError(
				"Extra ApplicationConfigValue Attribute Value",
				"While creating a ApplicationConfigValue value, an extra attribute value was detected. "+
					"A ApplicationConfigValue must not contain values beyond the expected attribute types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Extra ApplicationConfigValue Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewApplicationConfigValueUnknown(), diags
	}

	hostnameAttribute, ok := attributes["hostname"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`hostname is missing from object`)

		return NewApplicationConfigValueUnknown(), diags
	}

	hostnameVal, ok := hostnameAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`hostname expected to be basetypes.StringValue, was: %T`, hostnameAttribute))
	}

	instancesAttribute, ok := attributes["instances"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`instances is missing from object`)

		return NewApplicationConfigValueUnknown(), diags
	}

	instancesVal, ok := instancesAttribute.(basetypes.Int64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`instances expected to be basetypes.Int64Value, was: %T`, instancesAttribute))
	}

	ipAddressAttribute, ok := attributes["ip_address"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`ip_address is missing from object`)

		return NewApplicationConfigValueUnknown(), diags
	}

	ipAddressVal, ok := ipAddressAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`ip_address expected to be basetypes.StringValue, was: %T`, ipAddressAttribute))
	}

	passwordAttribute, ok := attributes["password"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`password is missing from object`)

		return NewApplicationConfigValueUnknown(), diags
	}

	passwordVal, ok := passwordAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`password expected to be basetypes.StringValue, was: %T`, passwordAttribute))
	}

	recoveryAttribute, ok := attributes["recovery"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`recovery is missing from object`)

		return NewApplicationConfigValueUnknown(), diags
	}

	recoveryVal, ok := recoveryAttribute.(basetypes.ObjectValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`recovery expected to be basetypes.ObjectValue, was: %T`, recoveryAttribute))
	}

	scheduledBackupsAttribute, ok := attributes["scheduled_backups"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`scheduled_backups is missing from object`)

		return NewApplicationConfigValueUnknown(), diags
	}

	scheduledBackupsVal, ok := scheduledBackupsAttribute.(basetypes.ObjectValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`scheduled_backups expected to be basetypes.ObjectValue, was: %T`, scheduledBackupsAttribute))
	}

	typeAttribute, ok := attributes["type"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`type is missing from object`)

		return NewApplicationConfigValueUnknown(), diags
	}

	typeVal, ok := typeAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`type expected to be basetypes.StringValue, was: %T`, typeAttribute))
	}

	versionAttribute, ok := attributes["version"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`version is missing from object`)

		return NewApplicationConfigValueUnknown(), diags
	}

	versionVal, ok := versionAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`version expected to be basetypes.StringValue, was: %T`, versionAttribute))
	}

	if diags.HasError() {
		return NewApplicationConfigValueUnknown(), diags
	}

	return ApplicationConfigValue{
		Hostname:              hostnameVal,
		Instances:             instancesVal,
		IpAddress:             ipAddressVal,
		Password:              passwordVal,
		Recovery:              recoveryVal,
		ScheduledBackups:      scheduledBackupsVal,
		ApplicationConfigType: typeVal,
		Version:               versionVal,
		state:                 attr.ValueStateKnown,
	}, diags
}

func NewApplicationConfigValueMust(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) ApplicationConfigValue {
	object, diags := NewApplicationConfigValue(attributeTypes, attributes)

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

		panic("NewApplicationConfigValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}

func (t ApplicationConfigType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewApplicationConfigValueNull(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewApplicationConfigValueUnknown(), nil
	}

	if in.IsNull() {
		return NewApplicationConfigValueNull(), nil
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

	return NewApplicationConfigValueMust(ApplicationConfigValue{}.AttributeTypes(ctx), attributes), nil
}

func (t ApplicationConfigType) ValueType(ctx context.Context) attr.Value {
	return ApplicationConfigValue{}
}

var _ basetypes.ObjectValuable = ApplicationConfigValue{}

type ApplicationConfigValue struct {
	Hostname              basetypes.StringValue `tfsdk:"hostname"`
	Instances             basetypes.Int64Value  `tfsdk:"instances"`
	IpAddress             basetypes.StringValue `tfsdk:"ip_address"`
	Password              basetypes.StringValue `tfsdk:"password"`
	Recovery              basetypes.ObjectValue `tfsdk:"recovery"`
	ScheduledBackups      basetypes.ObjectValue `tfsdk:"scheduled_backups"`
	ApplicationConfigType basetypes.StringValue `tfsdk:"type"`
	Version               basetypes.StringValue `tfsdk:"version"`
	state                 attr.ValueState
}

func (v ApplicationConfigValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 8)

	var val tftypes.Value
	var err error

	attrTypes["hostname"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["instances"] = basetypes.Int64Type{}.TerraformType(ctx)
	attrTypes["ip_address"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["password"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["recovery"] = basetypes.ObjectType{
		AttrTypes: RecoveryValue{}.AttributeTypes(ctx),
	}.TerraformType(ctx)
	attrTypes["scheduled_backups"] = basetypes.ObjectType{
		AttrTypes: ScheduledBackupsValue{}.AttributeTypes(ctx),
	}.TerraformType(ctx)
	attrTypes["type"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["version"] = basetypes.StringType{}.TerraformType(ctx)

	objectType := tftypes.Object{AttributeTypes: attrTypes}

	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 8)

		val, err = v.Hostname.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["hostname"] = val

		val, err = v.Instances.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["instances"] = val

		val, err = v.IpAddress.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["ip_address"] = val

		val, err = v.Password.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["password"] = val

		val, err = v.Recovery.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["recovery"] = val

		val, err = v.ScheduledBackups.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["scheduled_backups"] = val

		val, err = v.ApplicationConfigType.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["type"] = val

		val, err = v.Version.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["version"] = val

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

func (v ApplicationConfigValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v ApplicationConfigValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v ApplicationConfigValue) String() string {
	return "ApplicationConfigValue"
}

func (v ApplicationConfigValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	var recovery basetypes.ObjectValue

	if v.Recovery.IsNull() {
		recovery = types.ObjectNull(
			RecoveryValue{}.AttributeTypes(ctx),
		)
	}

	if v.Recovery.IsUnknown() {
		recovery = types.ObjectUnknown(
			RecoveryValue{}.AttributeTypes(ctx),
		)
	}

	if !v.Recovery.IsNull() && !v.Recovery.IsUnknown() {
		recovery = types.ObjectValueMust(
			RecoveryValue{}.AttributeTypes(ctx),
			v.Recovery.Attributes(),
		)
	}

	var scheduledBackups basetypes.ObjectValue

	if v.ScheduledBackups.IsNull() {
		scheduledBackups = types.ObjectNull(
			ScheduledBackupsValue{}.AttributeTypes(ctx),
		)
	}

	if v.ScheduledBackups.IsUnknown() {
		scheduledBackups = types.ObjectUnknown(
			ScheduledBackupsValue{}.AttributeTypes(ctx),
		)
	}

	if !v.ScheduledBackups.IsNull() && !v.ScheduledBackups.IsUnknown() {
		scheduledBackups = types.ObjectValueMust(
			ScheduledBackupsValue{}.AttributeTypes(ctx),
			v.ScheduledBackups.Attributes(),
		)
	}

	objVal, diags := types.ObjectValue(
		map[string]attr.Type{
			"hostname":   basetypes.StringType{},
			"instances":  basetypes.Int64Type{},
			"ip_address": basetypes.StringType{},
			"password":   basetypes.StringType{},
			"recovery": basetypes.ObjectType{
				AttrTypes: RecoveryValue{}.AttributeTypes(ctx),
			},
			"scheduled_backups": basetypes.ObjectType{
				AttrTypes: ScheduledBackupsValue{}.AttributeTypes(ctx),
			},
			"type":    basetypes.StringType{},
			"version": basetypes.StringType{},
		},
		map[string]attr.Value{
			"hostname":          v.Hostname,
			"instances":         v.Instances,
			"ip_address":        v.IpAddress,
			"password":          v.Password,
			"recovery":          recovery,
			"scheduled_backups": scheduledBackups,
			"type":              v.ApplicationConfigType,
			"version":           v.Version,
		})

	return objVal, diags
}

func (v ApplicationConfigValue) Equal(o attr.Value) bool {
	other, ok := o.(ApplicationConfigValue)

	if !ok {
		return false
	}

	if v.state != other.state {
		return false
	}

	if v.state != attr.ValueStateKnown {
		return true
	}

	if !v.Hostname.Equal(other.Hostname) {
		return false
	}

	if !v.Instances.Equal(other.Instances) {
		return false
	}

	if !v.IpAddress.Equal(other.IpAddress) {
		return false
	}

	if !v.Password.Equal(other.Password) {
		return false
	}

	if !v.Recovery.Equal(other.Recovery) {
		return false
	}

	if !v.ScheduledBackups.Equal(other.ScheduledBackups) {
		return false
	}

	if !v.ApplicationConfigType.Equal(other.ApplicationConfigType) {
		return false
	}

	if !v.Version.Equal(other.Version) {
		return false
	}

	return true
}

func (v ApplicationConfigValue) Type(ctx context.Context) attr.Type {
	return ApplicationConfigType{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v ApplicationConfigValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"hostname":   basetypes.StringType{},
		"instances":  basetypes.Int64Type{},
		"ip_address": basetypes.StringType{},
		"password":   basetypes.StringType{},
		"recovery": basetypes.ObjectType{
			AttrTypes: RecoveryValue{}.AttributeTypes(ctx),
		},
		"scheduled_backups": basetypes.ObjectType{
			AttrTypes: ScheduledBackupsValue{}.AttributeTypes(ctx),
		},
		"type":    basetypes.StringType{},
		"version": basetypes.StringType{},
	}
}

var _ basetypes.ObjectTypable = RecoveryType{}

type RecoveryType struct {
	basetypes.ObjectType
}

func (t RecoveryType) Equal(o attr.Type) bool {
	other, ok := o.(RecoveryType)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t RecoveryType) String() string {
	return "RecoveryType"
}

func (t RecoveryType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributes := in.Attributes()

	exclusiveAttribute, ok := attributes["exclusive"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`exclusive is missing from object`)

		return nil, diags
	}

	exclusiveVal, ok := exclusiveAttribute.(basetypes.BoolValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`exclusive expected to be basetypes.BoolValue, was: %T`, exclusiveAttribute))
	}

	sourceAttribute, ok := attributes["source"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`source is missing from object`)

		return nil, diags
	}

	sourceVal, ok := sourceAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`source expected to be basetypes.StringValue, was: %T`, sourceAttribute))
	}

	targetLsnAttribute, ok := attributes["target_lsn"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`target_lsn is missing from object`)

		return nil, diags
	}

	targetLsnVal, ok := targetLsnAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`target_lsn expected to be basetypes.StringValue, was: %T`, targetLsnAttribute))
	}

	targetNameAttribute, ok := attributes["target_name"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`target_name is missing from object`)

		return nil, diags
	}

	targetNameVal, ok := targetNameAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`target_name expected to be basetypes.StringValue, was: %T`, targetNameAttribute))
	}

	targetTimeAttribute, ok := attributes["target_time"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`target_time is missing from object`)

		return nil, diags
	}

	targetTimeVal, ok := targetTimeAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`target_time expected to be basetypes.StringValue, was: %T`, targetTimeAttribute))
	}

	targetXidAttribute, ok := attributes["target_xid"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`target_xid is missing from object`)

		return nil, diags
	}

	targetXidVal, ok := targetXidAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`target_xid expected to be basetypes.StringValue, was: %T`, targetXidAttribute))
	}

	if diags.HasError() {
		return nil, diags
	}

	return RecoveryValue{
		Exclusive:  exclusiveVal,
		Source:     sourceVal,
		TargetLsn:  targetLsnVal,
		TargetName: targetNameVal,
		TargetTime: targetTimeVal,
		TargetXid:  targetXidVal,
		state:      attr.ValueStateKnown,
	}, diags
}

func NewRecoveryValueNull() RecoveryValue {
	return RecoveryValue{
		state: attr.ValueStateNull,
	}
}

func NewRecoveryValueUnknown() RecoveryValue {
	return RecoveryValue{
		state: attr.ValueStateUnknown,
	}
}

func NewRecoveryValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (RecoveryValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]

		if !ok {
			diags.AddError(
				"Missing RecoveryValue Attribute Value",
				"While creating a RecoveryValue value, a missing attribute value was detected. "+
					"A RecoveryValue must contain values for all attributes, even if null or unknown. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("RecoveryValue Attribute Name (%s) Expected Type: %s", name, attributeType.String()),
			)

			continue
		}

		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError(
				"Invalid RecoveryValue Attribute Type",
				"While creating a RecoveryValue value, an invalid attribute value was detected. "+
					"A RecoveryValue must use a matching attribute type for the value. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("RecoveryValue Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
					fmt.Sprintf("RecoveryValue Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
			)
		}
	}

	for name := range attributes {
		_, ok := attributeTypes[name]

		if !ok {
			diags.AddError(
				"Extra RecoveryValue Attribute Value",
				"While creating a RecoveryValue value, an extra attribute value was detected. "+
					"A RecoveryValue must not contain values beyond the expected attribute types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Extra RecoveryValue Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewRecoveryValueUnknown(), diags
	}

	exclusiveAttribute, ok := attributes["exclusive"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`exclusive is missing from object`)

		return NewRecoveryValueUnknown(), diags
	}

	exclusiveVal, ok := exclusiveAttribute.(basetypes.BoolValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`exclusive expected to be basetypes.BoolValue, was: %T`, exclusiveAttribute))
	}

	sourceAttribute, ok := attributes["source"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`source is missing from object`)

		return NewRecoveryValueUnknown(), diags
	}

	sourceVal, ok := sourceAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`source expected to be basetypes.StringValue, was: %T`, sourceAttribute))
	}

	targetLsnAttribute, ok := attributes["target_lsn"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`target_lsn is missing from object`)

		return NewRecoveryValueUnknown(), diags
	}

	targetLsnVal, ok := targetLsnAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`target_lsn expected to be basetypes.StringValue, was: %T`, targetLsnAttribute))
	}

	targetNameAttribute, ok := attributes["target_name"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`target_name is missing from object`)

		return NewRecoveryValueUnknown(), diags
	}

	targetNameVal, ok := targetNameAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`target_name expected to be basetypes.StringValue, was: %T`, targetNameAttribute))
	}

	targetTimeAttribute, ok := attributes["target_time"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`target_time is missing from object`)

		return NewRecoveryValueUnknown(), diags
	}

	targetTimeVal, ok := targetTimeAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`target_time expected to be basetypes.StringValue, was: %T`, targetTimeAttribute))
	}

	targetXidAttribute, ok := attributes["target_xid"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`target_xid is missing from object`)

		return NewRecoveryValueUnknown(), diags
	}

	targetXidVal, ok := targetXidAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`target_xid expected to be basetypes.StringValue, was: %T`, targetXidAttribute))
	}

	if diags.HasError() {
		return NewRecoveryValueUnknown(), diags
	}

	return RecoveryValue{
		Exclusive:  exclusiveVal,
		Source:     sourceVal,
		TargetLsn:  targetLsnVal,
		TargetName: targetNameVal,
		TargetTime: targetTimeVal,
		TargetXid:  targetXidVal,
		state:      attr.ValueStateKnown,
	}, diags
}

func NewRecoveryValueMust(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) RecoveryValue {
	object, diags := NewRecoveryValue(attributeTypes, attributes)

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

		panic("NewRecoveryValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}

func (t RecoveryType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewRecoveryValueNull(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewRecoveryValueUnknown(), nil
	}

	if in.IsNull() {
		return NewRecoveryValueNull(), nil
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

	return NewRecoveryValueMust(RecoveryValue{}.AttributeTypes(ctx), attributes), nil
}

func (t RecoveryType) ValueType(ctx context.Context) attr.Value {
	return RecoveryValue{}
}

var _ basetypes.ObjectValuable = RecoveryValue{}

type RecoveryValue struct {
	Exclusive  basetypes.BoolValue   `tfsdk:"exclusive"`
	Source     basetypes.StringValue `tfsdk:"source"`
	TargetLsn  basetypes.StringValue `tfsdk:"target_lsn"`
	TargetName basetypes.StringValue `tfsdk:"target_name"`
	TargetTime basetypes.StringValue `tfsdk:"target_time"`
	TargetXid  basetypes.StringValue `tfsdk:"target_xid"`
	state      attr.ValueState
}

func (v RecoveryValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 6)

	var val tftypes.Value
	var err error

	attrTypes["exclusive"] = basetypes.BoolType{}.TerraformType(ctx)
	attrTypes["source"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["target_lsn"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["target_name"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["target_time"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["target_xid"] = basetypes.StringType{}.TerraformType(ctx)

	objectType := tftypes.Object{AttributeTypes: attrTypes}

	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 6)

		val, err = v.Exclusive.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["exclusive"] = val

		val, err = v.Source.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["source"] = val

		val, err = v.TargetLsn.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["target_lsn"] = val

		val, err = v.TargetName.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["target_name"] = val

		val, err = v.TargetTime.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["target_time"] = val

		val, err = v.TargetXid.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["target_xid"] = val

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

func (v RecoveryValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v RecoveryValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v RecoveryValue) String() string {
	return "RecoveryValue"
}

func (v RecoveryValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	objVal, diags := types.ObjectValue(
		map[string]attr.Type{
			"exclusive":   basetypes.BoolType{},
			"source":      basetypes.StringType{},
			"target_lsn":  basetypes.StringType{},
			"target_name": basetypes.StringType{},
			"target_time": basetypes.StringType{},
			"target_xid":  basetypes.StringType{},
		},
		map[string]attr.Value{
			"exclusive":   v.Exclusive,
			"source":      v.Source,
			"target_lsn":  v.TargetLsn,
			"target_name": v.TargetName,
			"target_time": v.TargetTime,
			"target_xid":  v.TargetXid,
		})

	return objVal, diags
}

func (v RecoveryValue) Equal(o attr.Value) bool {
	other, ok := o.(RecoveryValue)

	if !ok {
		return false
	}

	if v.state != other.state {
		return false
	}

	if v.state != attr.ValueStateKnown {
		return true
	}

	if !v.Exclusive.Equal(other.Exclusive) {
		return false
	}

	if !v.Source.Equal(other.Source) {
		return false
	}

	if !v.TargetLsn.Equal(other.TargetLsn) {
		return false
	}

	if !v.TargetName.Equal(other.TargetName) {
		return false
	}

	if !v.TargetTime.Equal(other.TargetTime) {
		return false
	}

	if !v.TargetXid.Equal(other.TargetXid) {
		return false
	}

	return true
}

func (v RecoveryValue) Type(ctx context.Context) attr.Type {
	return RecoveryType{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v RecoveryValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"exclusive":   basetypes.BoolType{},
		"source":      basetypes.StringType{},
		"target_lsn":  basetypes.StringType{},
		"target_name": basetypes.StringType{},
		"target_time": basetypes.StringType{},
		"target_xid":  basetypes.StringType{},
	}
}

var _ basetypes.ObjectTypable = ScheduledBackupsType{}

type ScheduledBackupsType struct {
	basetypes.ObjectType
}

func (t ScheduledBackupsType) Equal(o attr.Type) bool {
	other, ok := o.(ScheduledBackupsType)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t ScheduledBackupsType) String() string {
	return "ScheduledBackupsType"
}

func (t ScheduledBackupsType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributes := in.Attributes()

	retentionAttribute, ok := attributes["retention"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`retention is missing from object`)

		return nil, diags
	}

	retentionVal, ok := retentionAttribute.(basetypes.Int64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`retention expected to be basetypes.Int64Value, was: %T`, retentionAttribute))
	}

	scheduleAttribute, ok := attributes["schedule"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`schedule is missing from object`)

		return nil, diags
	}

	scheduleVal, ok := scheduleAttribute.(basetypes.ObjectValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`schedule expected to be basetypes.ObjectValue, was: %T`, scheduleAttribute))
	}

	if diags.HasError() {
		return nil, diags
	}

	return ScheduledBackupsValue{
		Retention: retentionVal,
		Schedule:  scheduleVal,
		state:     attr.ValueStateKnown,
	}, diags
}

func NewScheduledBackupsValueNull() ScheduledBackupsValue {
	return ScheduledBackupsValue{
		state: attr.ValueStateNull,
	}
}

func NewScheduledBackupsValueUnknown() ScheduledBackupsValue {
	return ScheduledBackupsValue{
		state: attr.ValueStateUnknown,
	}
}

func NewScheduledBackupsValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (ScheduledBackupsValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]

		if !ok {
			diags.AddError(
				"Missing ScheduledBackupsValue Attribute Value",
				"While creating a ScheduledBackupsValue value, a missing attribute value was detected. "+
					"A ScheduledBackupsValue must contain values for all attributes, even if null or unknown. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("ScheduledBackupsValue Attribute Name (%s) Expected Type: %s", name, attributeType.String()),
			)

			continue
		}

		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError(
				"Invalid ScheduledBackupsValue Attribute Type",
				"While creating a ScheduledBackupsValue value, an invalid attribute value was detected. "+
					"A ScheduledBackupsValue must use a matching attribute type for the value. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("ScheduledBackupsValue Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
					fmt.Sprintf("ScheduledBackupsValue Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
			)
		}
	}

	for name := range attributes {
		_, ok := attributeTypes[name]

		if !ok {
			diags.AddError(
				"Extra ScheduledBackupsValue Attribute Value",
				"While creating a ScheduledBackupsValue value, an extra attribute value was detected. "+
					"A ScheduledBackupsValue must not contain values beyond the expected attribute types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Extra ScheduledBackupsValue Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewScheduledBackupsValueUnknown(), diags
	}

	retentionAttribute, ok := attributes["retention"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`retention is missing from object`)

		return NewScheduledBackupsValueUnknown(), diags
	}

	retentionVal, ok := retentionAttribute.(basetypes.Int64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`retention expected to be basetypes.Int64Value, was: %T`, retentionAttribute))
	}

	scheduleAttribute, ok := attributes["schedule"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`schedule is missing from object`)

		return NewScheduledBackupsValueUnknown(), diags
	}

	scheduleVal, ok := scheduleAttribute.(basetypes.ObjectValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`schedule expected to be basetypes.ObjectValue, was: %T`, scheduleAttribute))
	}

	if diags.HasError() {
		return NewScheduledBackupsValueUnknown(), diags
	}

	return ScheduledBackupsValue{
		Retention: retentionVal,
		Schedule:  scheduleVal,
		state:     attr.ValueStateKnown,
	}, diags
}

func NewScheduledBackupsValueMust(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) ScheduledBackupsValue {
	object, diags := NewScheduledBackupsValue(attributeTypes, attributes)

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

		panic("NewScheduledBackupsValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}

func (t ScheduledBackupsType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewScheduledBackupsValueNull(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewScheduledBackupsValueUnknown(), nil
	}

	if in.IsNull() {
		return NewScheduledBackupsValueNull(), nil
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

	return NewScheduledBackupsValueMust(ScheduledBackupsValue{}.AttributeTypes(ctx), attributes), nil
}

func (t ScheduledBackupsType) ValueType(ctx context.Context) attr.Value {
	return ScheduledBackupsValue{}
}

var _ basetypes.ObjectValuable = ScheduledBackupsValue{}

type ScheduledBackupsValue struct {
	Retention basetypes.Int64Value  `tfsdk:"retention"`
	Schedule  basetypes.ObjectValue `tfsdk:"schedule"`
	state     attr.ValueState
}

func (v ScheduledBackupsValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 2)

	var val tftypes.Value
	var err error

	attrTypes["retention"] = basetypes.Int64Type{}.TerraformType(ctx)
	attrTypes["schedule"] = basetypes.ObjectType{
		AttrTypes: ScheduleValue{}.AttributeTypes(ctx),
	}.TerraformType(ctx)

	objectType := tftypes.Object{AttributeTypes: attrTypes}

	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 2)

		val, err = v.Retention.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["retention"] = val

		val, err = v.Schedule.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["schedule"] = val

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

func (v ScheduledBackupsValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v ScheduledBackupsValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v ScheduledBackupsValue) String() string {
	return "ScheduledBackupsValue"
}

func (v ScheduledBackupsValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	var schedule basetypes.ObjectValue

	if v.Schedule.IsNull() {
		schedule = types.ObjectNull(
			ScheduleValue{}.AttributeTypes(ctx),
		)
	}

	if v.Schedule.IsUnknown() {
		schedule = types.ObjectUnknown(
			ScheduleValue{}.AttributeTypes(ctx),
		)
	}

	if !v.Schedule.IsNull() && !v.Schedule.IsUnknown() {
		schedule = types.ObjectValueMust(
			ScheduleValue{}.AttributeTypes(ctx),
			v.Schedule.Attributes(),
		)
	}

	objVal, diags := types.ObjectValue(
		map[string]attr.Type{
			"retention": basetypes.Int64Type{},
			"schedule": basetypes.ObjectType{
				AttrTypes: ScheduleValue{}.AttributeTypes(ctx),
			},
		},
		map[string]attr.Value{
			"retention": v.Retention,
			"schedule":  schedule,
		})

	return objVal, diags
}

func (v ScheduledBackupsValue) Equal(o attr.Value) bool {
	other, ok := o.(ScheduledBackupsValue)

	if !ok {
		return false
	}

	if v.state != other.state {
		return false
	}

	if v.state != attr.ValueStateKnown {
		return true
	}

	if !v.Retention.Equal(other.Retention) {
		return false
	}

	if !v.Schedule.Equal(other.Schedule) {
		return false
	}

	return true
}

func (v ScheduledBackupsValue) Type(ctx context.Context) attr.Type {
	return ScheduledBackupsType{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v ScheduledBackupsValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"retention": basetypes.Int64Type{},
		"schedule": basetypes.ObjectType{
			AttrTypes: ScheduleValue{}.AttributeTypes(ctx),
		},
	}
}

func (v ScheduledBackupsValue) ToDBaaSSdkObject(ctx context.Context) (*sys11dbaassdk.PSQLScheduledBackups, diag.Diagnostics) {
	var diags diag.Diagnostics
	scheduleObj, d := NewScheduleValue(v.Schedule.AttributeTypes(ctx), v.Schedule.Attributes())
	diags.Append(d...)
	schedule, d := scheduleObj.ToDBaaSSdkObject(ctx)
	diags.Append(d...)

	var retention *int
	retention = nil
	if !v.Retention.IsNull() && !v.Retention.IsUnknown() {
		retention = sys11dbaassdk.Int64ToIntPtr(v.Retention.ValueInt64())
	}

	return &sys11dbaassdk.PSQLScheduledBackups{
		Retention: retention,
		Schedule:  schedule,
	}, diags
}

var _ basetypes.ObjectTypable = ScheduleType{}

type ScheduleType struct {
	basetypes.ObjectType
}

func (t ScheduleType) Equal(o attr.Type) bool {
	other, ok := o.(ScheduleType)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t ScheduleType) String() string {
	return "ScheduleType"
}

func (t ScheduleType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributes := in.Attributes()

	hourAttribute, ok := attributes["hour"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`hour is missing from object`)

		return nil, diags
	}

	hourVal, ok := hourAttribute.(basetypes.Int64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`hour expected to be basetypes.Int64Value, was: %T`, hourAttribute))
	}

	minuteAttribute, ok := attributes["minute"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`minute is missing from object`)

		return nil, diags
	}

	minuteVal, ok := minuteAttribute.(basetypes.Int64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`minute expected to be basetypes.Int64Value, was: %T`, minuteAttribute))
	}

	if diags.HasError() {
		return nil, diags
	}

	return ScheduleValue{
		Hour:   hourVal,
		Minute: minuteVal,
		state:  attr.ValueStateKnown,
	}, diags
}

func NewScheduleValueNull() ScheduleValue {
	return ScheduleValue{
		state: attr.ValueStateNull,
	}
}

func NewScheduleValueUnknown() ScheduleValue {
	return ScheduleValue{
		state: attr.ValueStateUnknown,
	}
}

func NewScheduleValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (ScheduleValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]

		if !ok {
			diags.AddError(
				"Missing ScheduleValue Attribute Value",
				"While creating a ScheduleValue value, a missing attribute value was detected. "+
					"A ScheduleValue must contain values for all attributes, even if null or unknown. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("ScheduleValue Attribute Name (%s) Expected Type: %s", name, attributeType.String()),
			)

			continue
		}

		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError(
				"Invalid ScheduleValue Attribute Type",
				"While creating a ScheduleValue value, an invalid attribute value was detected. "+
					"A ScheduleValue must use a matching attribute type for the value. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("ScheduleValue Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
					fmt.Sprintf("ScheduleValue Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
			)
		}
	}

	for name := range attributes {
		_, ok := attributeTypes[name]

		if !ok {
			diags.AddError(
				"Extra ScheduleValue Attribute Value",
				"While creating a ScheduleValue value, an extra attribute value was detected. "+
					"A ScheduleValue must not contain values beyond the expected attribute types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Extra ScheduleValue Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewScheduleValueUnknown(), diags
	}

	hourAttribute, ok := attributes["hour"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`hour is missing from object`)

		return NewScheduleValueUnknown(), diags
	}

	hourVal, ok := hourAttribute.(basetypes.Int64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`hour expected to be basetypes.Int64Value, was: %T`, hourAttribute))
	}

	minuteAttribute, ok := attributes["minute"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`minute is missing from object`)

		return NewScheduleValueUnknown(), diags
	}

	minuteVal, ok := minuteAttribute.(basetypes.Int64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`minute expected to be basetypes.Int64Value, was: %T`, minuteAttribute))
	}

	if diags.HasError() {
		return NewScheduleValueUnknown(), diags
	}

	return ScheduleValue{
		Hour:   hourVal,
		Minute: minuteVal,
		state:  attr.ValueStateKnown,
	}, diags
}

func NewScheduleValueMust(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) ScheduleValue {
	object, diags := NewScheduleValue(attributeTypes, attributes)

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

		panic("NewScheduleValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}

func (t ScheduleType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewScheduleValueNull(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewScheduleValueUnknown(), nil
	}

	if in.IsNull() {
		return NewScheduleValueNull(), nil
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

	return NewScheduleValueMust(ScheduleValue{}.AttributeTypes(ctx), attributes), nil
}

func (t ScheduleType) ValueType(ctx context.Context) attr.Value {
	return ScheduleValue{}
}

var _ basetypes.ObjectValuable = ScheduleValue{}

type ScheduleValue struct {
	Hour   basetypes.Int64Value `tfsdk:"hour"`
	Minute basetypes.Int64Value `tfsdk:"minute"`
	state  attr.ValueState
}

func (v ScheduleValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 2)

	var val tftypes.Value
	var err error

	attrTypes["hour"] = basetypes.Int64Type{}.TerraformType(ctx)
	attrTypes["minute"] = basetypes.Int64Type{}.TerraformType(ctx)

	objectType := tftypes.Object{AttributeTypes: attrTypes}

	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 2)

		val, err = v.Hour.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["hour"] = val

		val, err = v.Minute.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["minute"] = val

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

func (v ScheduleValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v ScheduleValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v ScheduleValue) String() string {
	return "ScheduleValue"
}

func (v ScheduleValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	objVal, diags := types.ObjectValue(
		map[string]attr.Type{
			"hour":   basetypes.Int64Type{},
			"minute": basetypes.Int64Type{},
		},
		map[string]attr.Value{
			"hour":   v.Hour,
			"minute": v.Minute,
		})

	return objVal, diags
}

func (v ScheduleValue) Equal(o attr.Value) bool {
	other, ok := o.(ScheduleValue)

	if !ok {
		return false
	}

	if v.state != other.state {
		return false
	}

	if v.state != attr.ValueStateKnown {
		return true
	}

	if !v.Hour.Equal(other.Hour) {
		return false
	}

	if !v.Minute.Equal(other.Minute) {
		return false
	}

	return true
}

func (v ScheduleValue) Type(ctx context.Context) attr.Type {
	return ScheduleType{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v ScheduleValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"hour":   basetypes.Int64Type{},
		"minute": basetypes.Int64Type{},
	}
}

func (v ScheduleValue) ToDBaaSSdkObject(ctx context.Context) (*sys11dbaassdk.PSQLScheduledBackupsSchedule, diag.Diagnostics) {

	var hour *int
	hour = nil
	if !v.Hour.IsNull() && !v.Hour.IsUnknown() {
		hour = sys11dbaassdk.Int64ToIntPtr(v.Hour.ValueInt64())
	}

	var minute *int
	minute = nil
	if !v.Minute.IsNull() && !v.Minute.IsUnknown() {
		minute = sys11dbaassdk.Int64ToIntPtr(v.Minute.ValueInt64())
	}

	return &sys11dbaassdk.PSQLScheduledBackupsSchedule{
		Hour:   hour,
		Minute: minute,
	}, diag.Diagnostics{}
}
