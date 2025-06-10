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

var _ basetypes.ObjectTypable = ApplicationConfigTypeV2{}

type ApplicationConfigTypeV2 struct {
	basetypes.ObjectType
}

func (t ApplicationConfigTypeV2) Equal(o attr.Type) bool {
	other, ok := o.(ApplicationConfigTypeV2)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t ApplicationConfigTypeV2) String() string {
	return "ApplicationConfigTypeV2"
}

func (t ApplicationConfigTypeV2) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributes := in.Attributes()

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

	return ApplicationConfigValueV2{
		Instances:             instancesVal,
		Password:              passwordVal,
		Recovery:              recoveryVal,
		ScheduledBackups:      scheduledBackupsVal,
		ApplicationConfigType: typeVal,
		Version:               versionVal,
		state:                 attr.ValueStateKnown,
	}, diags
}

func NewApplicationConfigValueV2Null() ApplicationConfigValueV2 {
	return ApplicationConfigValueV2{
		state: attr.ValueStateNull,
	}
}

func NewApplicationConfigValueV2Unknown() ApplicationConfigValueV2 {
	return ApplicationConfigValueV2{
		state: attr.ValueStateUnknown,
	}
}

func NewApplicationConfigValueV2(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (ApplicationConfigValueV2, diag.Diagnostics) {
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
		return NewApplicationConfigValueV2Unknown(), diags
	}

	instancesAttribute, ok := attributes["instances"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`instances is missing from object`)

		return NewApplicationConfigValueV2Unknown(), diags
	}

	instancesVal, ok := instancesAttribute.(basetypes.Int64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`instances expected to be basetypes.Int64Value, was: %T`, instancesAttribute))
	}

	passwordAttribute, ok := attributes["password"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`password is missing from object`)

		return NewApplicationConfigValueV2Unknown(), diags
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

		return NewApplicationConfigValueV2Unknown(), diags
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

		return NewApplicationConfigValueV2Unknown(), diags
	}

	scheduledBackupsVal, ok := scheduledBackupsAttribute.(basetypes.ObjectValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`scheduled_backups expected to be basetypes.ObjectValue, was: %T`, scheduledBackupsAttribute))
	}

	privateNetworkingAttribute, ok := attributes["private_networking"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`private_networking is missing from object`)

		return NewApplicationConfigValueV2Unknown(), diags
	}

	privateNetworkingVal, ok := privateNetworkingAttribute.(basetypes.ObjectValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`private_networking expected to be basetypes.ObjectValue, was: %T`, privateNetworkingAttribute))
	}

	publicNetworkingAttribute, ok := attributes["public_networking"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`public_networking is missing from object`)

		return NewApplicationConfigValueV2Unknown(), diags
	}

	publicNetworkingVal, ok := publicNetworkingAttribute.(basetypes.ObjectValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`public_networking expected to be basetypes.ObjectValue, was: %T`, publicNetworkingAttribute))
	}

	typeAttribute, ok := attributes["type"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`type is missing from object`)

		return NewApplicationConfigValueV2Unknown(), diags
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

		return NewApplicationConfigValueV2Unknown(), diags
	}

	versionVal, ok := versionAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`version expected to be basetypes.StringValue, was: %T`, versionAttribute))
	}

	if diags.HasError() {
		return NewApplicationConfigValueV2Unknown(), diags
	}

	return ApplicationConfigValueV2{
		Instances:             instancesVal,
		Password:              passwordVal,
		Recovery:              recoveryVal,
		ScheduledBackups:      scheduledBackupsVal,
		PrivateNetworking:     privateNetworkingVal,
		PublicNetworking:      publicNetworkingVal,
		ApplicationConfigType: typeVal,
		Version:               versionVal,
		state:                 attr.ValueStateKnown,
	}, diags
}

func NewApplicationConfigValueV2Must(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) ApplicationConfigValueV2 {
	object, diags := NewApplicationConfigValueV2(attributeTypes, attributes)

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

func (t ApplicationConfigTypeV2) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t ApplicationConfigTypeV2) ValueType(ctx context.Context) attr.Value {
	return ApplicationConfigValue{}
}

var _ basetypes.ObjectValuable = ApplicationConfigValueV2{}

type ApplicationConfigValueV2 struct {
	Instances             basetypes.Int64Value  `tfsdk:"instances"`
	Password              basetypes.StringValue `tfsdk:"password"`
	Recovery              basetypes.ObjectValue `tfsdk:"recovery"`
	ScheduledBackups      basetypes.ObjectValue `tfsdk:"scheduled_backups"`
	PrivateNetworking     basetypes.ObjectValue `tfsdk:"private_networking"`
	PublicNetworking      basetypes.ObjectValue `tfsdk:"public_networking"`
	ApplicationConfigType basetypes.StringValue `tfsdk:"type"`
	Version               basetypes.StringValue `tfsdk:"version"`
	state                 attr.ValueState
}

func (v ApplicationConfigValueV2) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 8)

	var val tftypes.Value
	var err error

	attrTypes["instances"] = basetypes.Int64Type{}.TerraformType(ctx)
	attrTypes["password"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["recovery"] = basetypes.ObjectType{
		AttrTypes: RecoveryValueV2{}.AttributeTypes(ctx),
	}.TerraformType(ctx)
	attrTypes["scheduled_backups"] = basetypes.ObjectType{
		AttrTypes: ScheduledBackupsValueV2{}.AttributeTypes(ctx),
	}.TerraformType(ctx)
	attrTypes["private_networking"] = basetypes.ObjectType{
		AttrTypes: PrivateNetworkingValueV2{}.AttributeTypes(ctx),
	}.TerraformType(ctx)
	attrTypes["public_networking"] = basetypes.ObjectType{
		AttrTypes: PublicNetworkingValueV2{}.AttributeTypes(ctx),
	}.TerraformType(ctx)
	attrTypes["type"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["version"] = basetypes.StringType{}.TerraformType(ctx)

	objectType := tftypes.Object{AttributeTypes: attrTypes}

	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 8)

		val, err = v.Instances.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["instances"] = val

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

		val, err = v.PrivateNetworking.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["private_networking"] = val

		val, err = v.PublicNetworking.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["public_networking"] = val

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

func (v ApplicationConfigValueV2) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v ApplicationConfigValueV2) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v ApplicationConfigValueV2) String() string {
	return "ApplicationConfigValueV2"
}

func (v ApplicationConfigValueV2) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
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

	var privateNetworking basetypes.ObjectValue

	if v.PrivateNetworking.IsNull() {
		privateNetworking = types.ObjectNull(
			PrivateNetworkingValueV2{}.AttributeTypes(ctx),
		)
	}

	if v.PrivateNetworking.IsUnknown() {
		privateNetworking = types.ObjectUnknown(
			PrivateNetworkingValueV2{}.AttributeTypes(ctx),
		)
	}

	if !v.PrivateNetworking.IsNull() && !v.PrivateNetworking.IsUnknown() {
		privateNetworking = types.ObjectValueMust(
			PrivateNetworkingValueV2{}.AttributeTypes(ctx),
			v.PrivateNetworking.Attributes(),
		)
	}

	var publicNetworking basetypes.ObjectValue

	if v.PublicNetworking.IsNull() {
		publicNetworking = types.ObjectNull(
			PublicNetworkingValueV2{}.AttributeTypes(ctx),
		)
	}

	if v.PublicNetworking.IsUnknown() {
		publicNetworking = types.ObjectUnknown(
			PublicNetworkingValueV2{}.AttributeTypes(ctx),
		)
	}

	if !v.PublicNetworking.IsNull() && !v.PublicNetworking.IsUnknown() {
		publicNetworking = types.ObjectValueMust(
			PublicNetworkingValueV2{}.AttributeTypes(ctx),
			v.PublicNetworking.Attributes(),
		)
	}

	objVal, diags := types.ObjectValue(
		map[string]attr.Type{
			"instances": basetypes.Int64Type{},
			"password":  basetypes.StringType{},
			"recovery": basetypes.ObjectType{
				AttrTypes: RecoveryValueV2{}.AttributeTypes(ctx),
			},
			"scheduled_backups": basetypes.ObjectType{
				AttrTypes: ScheduledBackupsValueV2{}.AttributeTypes(ctx),
			},
			"private_networking": basetypes.ObjectType{
				AttrTypes: PrivateNetworkingValueV2{}.AttributeTypes(ctx),
			},
			"public_networking": basetypes.ObjectType{
				AttrTypes: PublicNetworkingValueV2{}.AttributeTypes(ctx),
			},
			"type":    basetypes.StringType{},
			"version": basetypes.StringType{},
		},
		map[string]attr.Value{
			"instances":          v.Instances,
			"password":           v.Password,
			"recovery":           recovery,
			"scheduled_backups":  scheduledBackups,
			"private_networking": privateNetworking,
			"public_networking":  publicNetworking,
			"type":               v.ApplicationConfigType,
			"version":            v.Version,
		})

	return objVal, diags
}

func (v ApplicationConfigValueV2) Equal(o attr.Value) bool {
	other, ok := o.(ApplicationConfigValueV2)

	if !ok {
		return false
	}

	if v.state != other.state {
		return false
	}

	if v.state != attr.ValueStateKnown {
		return true
	}

	if !v.Instances.Equal(other.Instances) {
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

	if !v.PrivateNetworking.Equal(other.PrivateNetworking) {
		return false
	}

	if !v.PublicNetworking.Equal(other.PublicNetworking) {
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

func (v ApplicationConfigValueV2) Type(ctx context.Context) attr.Type {
	return ApplicationConfigType{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v ApplicationConfigValueV2) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"instances": basetypes.Int64Type{},
		"password":  basetypes.StringType{},
		"recovery": basetypes.ObjectType{
			AttrTypes: RecoveryValueV2{}.AttributeTypes(ctx),
		},
		"scheduled_backups": basetypes.ObjectType{
			AttrTypes: ScheduledBackupsValueV2{}.AttributeTypes(ctx),
		},
		"private_networking": basetypes.ObjectType{
			AttrTypes: PrivateNetworkingValueV2{}.AttributeTypes(ctx),
		},
		"public_networking": basetypes.ObjectType{
			AttrTypes: PublicNetworkingValueV2{}.AttributeTypes(ctx),
		},
		"type":    basetypes.StringType{},
		"version": basetypes.StringType{},
	}
}

var _ basetypes.ObjectTypable = RecoveryType{}

type RecoveryTypeV2 struct {
	basetypes.ObjectType
}

func (t RecoveryTypeV2) Equal(o attr.Type) bool {
	other, ok := o.(RecoveryType)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t RecoveryTypeV2) String() string {
	return "RecoveryTypeV2"
}

func (t RecoveryTypeV2) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
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

	return RecoveryValueV2{
		Exclusive:  exclusiveVal,
		Source:     sourceVal,
		TargetLsn:  targetLsnVal,
		TargetName: targetNameVal,
		TargetTime: targetTimeVal,
		TargetXid:  targetXidVal,
		state:      attr.ValueStateKnown,
	}, diags
}

func NewRecoveryValueV2Null() RecoveryValueV2 {
	return RecoveryValueV2{
		state: attr.ValueStateNull,
	}
}

func NewRecoveryValueV2Unknown() RecoveryValueV2 {
	return RecoveryValueV2{
		state: attr.ValueStateUnknown,
	}
}

func NewRecoveryValueV2(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (RecoveryValueV2, diag.Diagnostics) {
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
		return NewRecoveryValueV2Unknown(), diags
	}

	exclusiveAttribute, ok := attributes["exclusive"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`exclusive is missing from object`)

		return NewRecoveryValueV2Unknown(), diags
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

		return NewRecoveryValueV2Unknown(), diags
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

		return NewRecoveryValueV2Unknown(), diags
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

		return NewRecoveryValueV2Unknown(), diags
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

		return NewRecoveryValueV2Unknown(), diags
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

		return NewRecoveryValueV2Unknown(), diags
	}

	targetXidVal, ok := targetXidAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`target_xid expected to be basetypes.StringValue, was: %T`, targetXidAttribute))
	}

	if diags.HasError() {
		return NewRecoveryValueV2Unknown(), diags
	}

	return RecoveryValueV2{
		Exclusive:  exclusiveVal,
		Source:     sourceVal,
		TargetLsn:  targetLsnVal,
		TargetName: targetNameVal,
		TargetTime: targetTimeVal,
		TargetXid:  targetXidVal,
		state:      attr.ValueStateKnown,
	}, diags
}

func NewRecoveryValueV2Must(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) RecoveryValueV2 {
	object, diags := NewRecoveryValueV2(attributeTypes, attributes)

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

func (t RecoveryTypeV2) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t RecoveryTypeV2) ValueType(ctx context.Context) attr.Value {
	return RecoveryValueV2{}
}

var _ basetypes.ObjectValuable = RecoveryValue{}

type RecoveryValueV2 struct {
	Exclusive  basetypes.BoolValue   `tfsdk:"exclusive"`
	Source     basetypes.StringValue `tfsdk:"source"`
	TargetLsn  basetypes.StringValue `tfsdk:"target_lsn"`
	TargetName basetypes.StringValue `tfsdk:"target_name"`
	TargetTime basetypes.StringValue `tfsdk:"target_time"`
	TargetXid  basetypes.StringValue `tfsdk:"target_xid"`
	state      attr.ValueState
}

func (v RecoveryValueV2) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
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

func (v RecoveryValueV2) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v RecoveryValueV2) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v RecoveryValueV2) String() string {
	return "RecoveryValueV2"
}

func (v RecoveryValueV2) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
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

func (v RecoveryValueV2) Equal(o attr.Value) bool {
	other, ok := o.(RecoveryValueV2)

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

func (v RecoveryValueV2) Type(ctx context.Context) attr.Type {
	return RecoveryType{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v RecoveryValueV2) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"exclusive":   basetypes.BoolType{},
		"source":      basetypes.StringType{},
		"target_lsn":  basetypes.StringType{},
		"target_name": basetypes.StringType{},
		"target_time": basetypes.StringType{},
		"target_xid":  basetypes.StringType{},
	}
}

var _ basetypes.ObjectTypable = ScheduledBackupsTypeV2{}

type ScheduledBackupsTypeV2 struct {
	basetypes.ObjectType
}

func (t ScheduledBackupsTypeV2) Equal(o attr.Type) bool {
	other, ok := o.(ScheduledBackupsTypeV2)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t ScheduledBackupsTypeV2) String() string {
	return "ScheduledBackupsTypeV2"
}

func (t ScheduledBackupsTypeV2) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
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

func NewScheduledBackupsValueV2Null() ScheduledBackupsValueV2 {
	return ScheduledBackupsValueV2{
		state: attr.ValueStateNull,
	}
}

func NewScheduledBackupsValueV2Unknown() ScheduledBackupsValueV2 {
	return ScheduledBackupsValueV2{
		state: attr.ValueStateUnknown,
	}
}

func NewScheduledBackupsValueV2(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (ScheduledBackupsValueV2, diag.Diagnostics) {
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
		return NewScheduledBackupsValueV2Unknown(), diags
	}

	retentionAttribute, ok := attributes["retention"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`retention is missing from object`)

		return NewScheduledBackupsValueV2Unknown(), diags
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

		return NewScheduledBackupsValueV2Unknown(), diags
	}

	scheduleVal, ok := scheduleAttribute.(basetypes.ObjectValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`schedule expected to be basetypes.ObjectValue, was: %T`, scheduleAttribute))
	}

	if diags.HasError() {
		return NewScheduledBackupsValueV2Unknown(), diags
	}

	return ScheduledBackupsValueV2{
		Retention: retentionVal,
		Schedule:  scheduleVal,
		state:     attr.ValueStateKnown,
	}, diags
}

func NewScheduledBackupsValueV2Must(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) ScheduledBackupsValueV2 {
	object, diags := NewScheduledBackupsValueV2(attributeTypes, attributes)

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

func (t ScheduledBackupsTypeV2) ValueFromTerraformV2(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t ScheduledBackupsTypeV2) ValueType(ctx context.Context) attr.Value {
	return ScheduledBackupsValue{}
}

var _ basetypes.ObjectValuable = ScheduledBackupsValueV2{}

type ScheduledBackupsValueV2 struct {
	Retention basetypes.Int64Value  `tfsdk:"retention"`
	Schedule  basetypes.ObjectValue `tfsdk:"schedule"`
	state     attr.ValueState
}

func (v ScheduledBackupsValueV2) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
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

func (v ScheduledBackupsValueV2) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v ScheduledBackupsValueV2) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v ScheduledBackupsValueV2) String() string {
	return "ScheduledBackupsValueV2"
}

func (v ScheduledBackupsValueV2) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
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

func (v ScheduledBackupsValueV2) Equal(o attr.Value) bool {
	other, ok := o.(ScheduledBackupsValueV2)

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

func (v ScheduledBackupsValueV2) Type(ctx context.Context) attr.Type {
	return ScheduledBackupsTypeV2{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v ScheduledBackupsValueV2) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"retention": basetypes.Int64Type{},
		"schedule": basetypes.ObjectType{
			AttrTypes: ScheduleValue{}.AttributeTypes(ctx),
		},
	}
}

func (v ScheduledBackupsValueV2) ToDBaaSSdkObject(ctx context.Context) (*sys11dbaassdk.PSQLScheduledBackupsV2, diag.Diagnostics) {
	var diags diag.Diagnostics
	scheduleObj, d := NewScheduleValueV2(v.Schedule.AttributeTypes(ctx), v.Schedule.Attributes())
	diags.Append(d...)
	schedule, d := scheduleObj.ToDBaaSSdkObject(ctx)
	diags.Append(d...)

	var retention *int
	retention = nil
	if !v.Retention.IsNull() && !v.Retention.IsUnknown() {
		retention = sys11dbaassdk.Int64ToIntPtr(v.Retention.ValueInt64())
	}

	return &sys11dbaassdk.PSQLScheduledBackupsV2{
		Retention: retention,
		Schedule:  schedule,
	}, diags
}

var _ basetypes.ObjectTypable = ScheduleTypeV2{}

type ScheduleTypeV2 struct {
	basetypes.ObjectType
}

func (t ScheduleTypeV2) Equal(o attr.Type) bool {
	other, ok := o.(ScheduleTypeV2)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t ScheduleTypeV2) String() string {
	return "ScheduleTypeV2"
}

func (t ScheduleTypeV2) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
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

func NewScheduleValueV2Null() ScheduleValueV2 {
	return ScheduleValueV2{
		state: attr.ValueStateNull,
	}
}

func NewScheduleValueV2Unknown() ScheduleValueV2 {
	return ScheduleValueV2{
		state: attr.ValueStateUnknown,
	}
}

func NewScheduleValueV2(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (ScheduleValueV2, diag.Diagnostics) {
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
		return NewScheduleValueV2Unknown(), diags
	}

	hourAttribute, ok := attributes["hour"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`hour is missing from object`)

		return NewScheduleValueV2Unknown(), diags
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

		return NewScheduleValueV2Unknown(), diags
	}

	minuteVal, ok := minuteAttribute.(basetypes.Int64Value)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`minute expected to be basetypes.Int64Value, was: %T`, minuteAttribute))
	}

	if diags.HasError() {
		return NewScheduleValueV2Unknown(), diags
	}

	return ScheduleValueV2{
		Hour:   hourVal,
		Minute: minuteVal,
		state:  attr.ValueStateKnown,
	}, diags
}

func NewScheduleValueV2Must(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) ScheduleValueV2 {
	object, diags := NewScheduleValueV2(attributeTypes, attributes)

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

		panic("NewScheduleValueV2Must received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}

func (t ScheduleTypeV2) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewScheduleValueV2Null(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewScheduleValueV2Unknown(), nil
	}

	if in.IsNull() {
		return NewScheduleValueV2Null(), nil
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

	return NewScheduleValueV2Must(ScheduleValueV2{}.AttributeTypes(ctx), attributes), nil
}

func (t ScheduleTypeV2) ValueType(ctx context.Context) attr.Value {
	return ScheduleValueV2{}
}

var _ basetypes.ObjectValuable = ScheduleValue{}

type ScheduleValueV2 struct {
	Hour   basetypes.Int64Value `tfsdk:"hour"`
	Minute basetypes.Int64Value `tfsdk:"minute"`
	state  attr.ValueState
}

func (v ScheduleValueV2) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
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

func (v ScheduleValueV2) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v ScheduleValueV2) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v ScheduleValueV2) String() string {
	return "ScheduleValueV2"
}

func (v ScheduleValueV2) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
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

func (v ScheduleValueV2) Equal(o attr.Value) bool {
	other, ok := o.(ScheduleValueV2)

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

func (v ScheduleValueV2) Type(ctx context.Context) attr.Type {
	return ScheduleType{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v ScheduleValueV2) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"hour":   basetypes.Int64Type{},
		"minute": basetypes.Int64Type{},
	}
}

func (v ScheduleValueV2) ToDBaaSSdkObject(ctx context.Context) (*sys11dbaassdk.PSQLScheduledBackupsScheduleV2, diag.Diagnostics) {

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

	return &sys11dbaassdk.PSQLScheduledBackupsScheduleV2{
		Hour:   hour,
		Minute: minute,
	}, diag.Diagnostics{}
}
