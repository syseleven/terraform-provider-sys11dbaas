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

var _ basetypes.ObjectTypable = PrivateNetworkConfigTypeV2{}

type PrivateNetworkConfigTypeV2 struct {
	basetypes.ObjectType
}

func (t PrivateNetworkConfigTypeV2) Equal(o attr.Type) bool {
	other, ok := o.(PrivateNetworkConfigTypeV2)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t PrivateNetworkConfigTypeV2) String() string {
	return "PrivateNetworkConfigTypeV2"
}

func (t PrivateNetworkConfigTypeV2) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributes := in.Attributes()

	enabledAttribute, ok := attributes["enabled"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`enabled is missing from object`)

		return nil, diags
	}

	enabledVal, ok := enabledAttribute.(basetypes.BoolValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`enabled expected to be basetypes.BoolValue, was: %T`, enabledAttribute))
	}

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

	sharedSubnetCIDRAttribute, ok := attributes["shared_subnet_cidr"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`shared_subnet_cidr is missing from object`)

		return nil, diags
	}

	sharedSubnetCIDRVal, ok := sharedSubnetCIDRAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`shared_subnet_cidr expected to be basetypes.StringValue, was: %T`, sharedSubnetCIDRAttribute))
	}

	allowedCIDRsAttribute, ok := attributes["allowed_cidrs"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`allowed_cidrs is missing from object`)

		return nil, diags
	}

	allowedCIDRsVal, ok := allowedCIDRsAttribute.(basetypes.ListValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`allowed_cidrs expected to be basetypes.ListValue, was: %T`, allowedCIDRsAttribute))
	}

	sharedSubnetIDAttribute, ok := attributes["shared_subnet_id"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`shared_subnet_id is missing from object`)

		return nil, diags
	}

	sharedSubnetIDVal, ok := sharedSubnetIDAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`shared_subnet_id expected to be basetypes.StringValue, was: %T`, sharedSubnetIDAttribute))
	}

	sharedNetworkIDAttribute, ok := attributes["shared_network_id"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`shared_network_id is missing from object`)

		return nil, diags
	}

	sharedNetworkIDVal, ok := sharedNetworkIDAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`shared_network_id expected to be basetypes.StringValue, was: %T`, sharedNetworkIDAttribute))
	}

	if diags.HasError() {
		return nil, diags
	}

	return PrivateNetworkConfigValueV2{
		Enabled:          enabledVal,
		Hostname:         hostnameVal,
		IPAddress:        ipAddressVal,
		SharedSubnetCIDR: sharedSubnetCIDRVal,
		AllowedCIDRs:     allowedCIDRsVal,
		SharedSubnetID:   sharedSubnetIDVal,
		SharedNetworkID:  sharedNetworkIDVal,
		state:            attr.ValueStateKnown,
	}, diags
}

func NewPrivateNetworkConfigValueV2Null() PrivateNetworkConfigValueV2 {
	return PrivateNetworkConfigValueV2{
		state: attr.ValueStateNull,
	}
}

func NewPrivateNetworkConfigValueV2Unknown() PrivateNetworkConfigValueV2 {
	return PrivateNetworkConfigValueV2{
		state: attr.ValueStateUnknown,
	}
}

func NewPrivateNetworkConfigValueV2(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (PrivateNetworkConfigValueV2, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]

		if !ok {
			diags.AddError(
				"Missing PrivateNetworkConfigValueV2 Attribute Value",
				"While creating a PrivateNetworkConfigValueV2 value, a missing attribute value was detected. "+
					"A PrivateNetworkConfigValueV2 must contain values for all attributes, even if null or unknown. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("PrivateNetworkConfigValueV2 Attribute Name (%s) Expected Type: %s", name, attributeType.String()),
			)

			continue
		}

		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError(
				"Invalid PrivateNetworkConfigValueV2 Attribute Type",
				"While creating a PrivateNetworkConfigValueV2 value, an invalid attribute value was detected. "+
					"A PrivateNetworkConfigValueV2 must use a matching attribute type for the value. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("PrivateNetworkConfigValueV2 Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
					fmt.Sprintf("PrivateNetworkConfigValueV2 Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
			)
		}
	}

	for name := range attributes {
		_, ok := attributeTypes[name]

		if !ok {
			diags.AddError(
				"Extra PrivateNetworkConfigValueV2 Attribute Value",
				"While creating a PrivateNetworkConfigValueV2 value, an extra attribute value was detected. "+
					"A PrivateNetworkConfigValueV2 must not contain values beyond the expected attribute types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Extra PrivateNetworkConfigValueV2 Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewPrivateNetworkConfigValueV2Unknown(), diags
	}

	enabledAttribute, ok := attributes["enabled"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`enabled is missing from object`)

		return NewPrivateNetworkConfigValueV2Unknown(), diags
	}

	enabledVal, ok := enabledAttribute.(basetypes.BoolValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`retention expected to be basetypes.BoolValue, was: %T`, enabledVal))
	}

	hostnameAttribute, ok := attributes["hostname"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`hostname is missing from object`)

		return NewPrivateNetworkConfigValueV2Unknown(), diags
	}

	hostnameVal, ok := hostnameAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`hostname expected to be basetypes.StringValue, was: %T`, hostnameAttribute))
	}

	ipAddressAttribute, ok := attributes["ip_address"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`ip_address is missing from object`)

		return NewPrivateNetworkConfigValueV2Unknown(), diags
	}

	ipAddressVal, ok := ipAddressAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`ip_address expected to be basetypes.StringValue, was: %T`, ipAddressAttribute))
	}

	sharedSubnetCIDRAttribute, ok := attributes["shared_subnet_cidr"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`shared_subnet_cidr is missing from object`)

		return NewPrivateNetworkConfigValueV2Unknown(), diags
	}

	sharedSubnetCIDRVal, ok := sharedSubnetCIDRAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`shared_subnet_cidr expected to be basetypes.StringValue, was: %T`, sharedSubnetCIDRAttribute))
	}

	allowedCIDRsAttribute, ok := attributes["allowed_cidrs"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`allowed_cidrs is missing from object`)

		return NewPrivateNetworkConfigValueV2Unknown(), diags
	}

	allowedCIDRsVal, ok := allowedCIDRsAttribute.(basetypes.ListValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`allowed_cidrs expected to be basetypes.ListValue, was: %T`, allowedCIDRsAttribute))
	}

	sharedSubnetIDAttribute, ok := attributes["shared_subnet_id"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`shared_subnet_id is missing from object`)

		return NewPrivateNetworkConfigValueV2Unknown(), diags
	}

	sharedSubnetIDVal, ok := sharedSubnetIDAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`shared_subnet_id expected to be basetypes.StringValue, was: %T`, sharedSubnetIDAttribute))
	}

	sharedNetworkIDAttribute, ok := attributes["shared_network_id"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`shared_network_id is missing from object`)

		return NewPrivateNetworkConfigValueV2Unknown(), diags
	}

	sharedNetworkIDVal, ok := sharedNetworkIDAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`shared_network_id expected to be basetypes.StringValue, was: %T`, sharedNetworkIDAttribute))
	}

	if diags.HasError() {
		return NewPrivateNetworkConfigValueV2Unknown(), diags
	}

	return PrivateNetworkConfigValueV2{
		Enabled:          enabledVal,
		AllowedCIDRs:     allowedCIDRsVal,
		Hostname:         hostnameVal,
		IPAddress:        ipAddressVal,
		SharedSubnetCIDR: sharedSubnetCIDRVal,
		SharedSubnetID:   sharedSubnetIDVal,
		SharedNetworkID:  sharedNetworkIDVal,
		state:            attr.ValueStateKnown,
	}, diags
}

func NewPrivateNetworkConfigValueV2Must(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) PrivateNetworkConfigValueV2 {
	object, diags := NewPrivateNetworkConfigValueV2(attributeTypes, attributes)

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

		panic("NewPrivateNetworkConfigValueV2Must received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}

func (t PrivateNetworkConfigTypeV2) ValueFromTerraformV2(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewPrivateNetworkConfigValueV2Null(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewPrivateNetworkConfigValueV2Unknown(), nil
	}

	if in.IsNull() {
		return NewPrivateNetworkConfigValueV2Null(), nil
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

	return NewPrivateNetworkConfigValueV2Must(PrivateNetworkConfigValueV2{}.AttributeTypes(ctx), attributes), nil
}

func (t PrivateNetworkConfigTypeV2) ValueType(ctx context.Context) attr.Value {
	return PrivateNetworkConfigValueV2{}
}

var _ basetypes.ObjectValuable = PrivateNetworkConfigValueV2{}

type PrivateNetworkConfigValueV2 struct {
	Enabled          basetypes.BoolValue   `tfsdk:"enabled"`
	AllowedCIDRs     basetypes.ListValue   `tfsdk:"allowed_cidrs"`
	SharedSubnetCIDR basetypes.StringValue `tfsdk:"shared_subnet_cidr"`
	Hostname         basetypes.StringValue `tfsdk:"hostname"`
	IPAddress        basetypes.StringValue `tfsdk:"ip_address"`
	SharedSubnetID   basetypes.StringValue `tfsdk:"shared_subnet_id"`
	SharedNetworkID  basetypes.StringValue `tfsdk:"shared_network_id"`
	state            attr.ValueState
}

func (v PrivateNetworkConfigValueV2) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 2)

	var val tftypes.Value
	var err error

	attrTypes["enabled"] = basetypes.BoolType{}.TerraformType(ctx)
	attrTypes["hostname"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["ip_address"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["allowed_cidrs"] = basetypes.ListType{ElemType: types.StringType}.TerraformType(ctx)
	attrTypes["shared_subnet_cidr"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["shared_subnet_id"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["shared_network_id"] = basetypes.StringType{}.TerraformType(ctx)

	objectType := tftypes.Object{AttributeTypes: attrTypes}

	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 2)

		val, err = v.Enabled.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["enabled"] = val

		val, err = v.Hostname.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["hostname"] = val

		val, err = v.IPAddress.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["ip_address"] = val

		val, err = v.AllowedCIDRs.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["allowed_cidrs"] = val

		val, err = v.SharedSubnetCIDR.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["shared_subnet_cidr"] = val

		val, err = v.SharedSubnetID.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["shared_subnet_id"] = val

		val, err = v.SharedNetworkID.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["shared_network_id"] = val

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

func (v PrivateNetworkConfigValueV2) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v PrivateNetworkConfigValueV2) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v PrivateNetworkConfigValueV2) String() string {
	return "PrivateNetworkConfigValueV2"
}

func (v PrivateNetworkConfigValueV2) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	objVal, diags := types.ObjectValue(
		map[string]attr.Type{
			"enabled":            basetypes.BoolType{},
			"hostname":           basetypes.StringType{},
			"ip_address":         basetypes.StringType{},
			"allowed_cidrs":      basetypes.ListType{ElemType: types.StringType},
			"shared_subnet_cidr": basetypes.StringType{},
			"shared_subnet_id":   basetypes.StringType{},
			"shared_network_id":  basetypes.StringType{},
		},
		map[string]attr.Value{
			"enabled":            v.Enabled,
			"hostname":           v.Hostname,
			"ip_address":         v.IPAddress,
			"allowed_cidrs":      v.AllowedCIDRs,
			"shared_subnet_cidr": v.SharedSubnetCIDR,
			"shared_subnet_id":   v.SharedSubnetID,
			"shared_network_id":  v.SharedNetworkID,
		})

	return objVal, diags
}

func (v PrivateNetworkConfigValueV2) Equal(o attr.Value) bool {
	other, ok := o.(PrivateNetworkConfigValueV2)

	if !ok {
		return false
	}

	if v.state != other.state {
		return false
	}

	if v.state != attr.ValueStateKnown {
		return true
	}

	if !v.Enabled.Equal(other.Enabled) {
		return false
	}

	if !v.Hostname.Equal(other.Hostname) {
		return false
	}

	if !v.IPAddress.Equal(other.IPAddress) {
		return false
	}

	if !v.AllowedCIDRs.Equal(other.AllowedCIDRs) {
		return false
	}

	if !v.SharedSubnetCIDR.Equal(other.SharedSubnetCIDR) {
		return false
	}

	if !v.SharedSubnetID.Equal(other.SharedSubnetID) {
		return false
	}

	if !v.SharedNetworkID.Equal(other.SharedNetworkID) {
		return false
	}

	return true
}

func (v PrivateNetworkConfigValueV2) Type(ctx context.Context) attr.Type {
	return PrivateNetworkConfigTypeV2{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v PrivateNetworkConfigValueV2) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"enabled":            basetypes.BoolType{},
		"hostname":           basetypes.StringType{},
		"ip_address":         basetypes.StringType{},
		"allowed_cidrs":      basetypes.ListType{ElemType: types.StringType},
		"shared_subnet_cidr": basetypes.StringType{},
		"shared_subnet_id":   basetypes.StringType{},
		"shared_network_id":  basetypes.StringType{},
	}
}

func (v PrivateNetworkConfigValueV2) ToDBaaSSdkResponse(ctx context.Context) (*sys11dbaassdk.PSQLPrivateNetworkConfigResponseV2, diag.Diagnostics) {
	var diags diag.Diagnostics

	var ipList []string

	for _, e := range v.AllowedCIDRs.Elements() {
		ipList = append(ipList, strings.Trim(e.String(), "\""))
	}

	return &sys11dbaassdk.PSQLPrivateNetworkConfigResponseV2{
		Enabled:          v.Enabled.ValueBool(),
		Hostname:         v.Hostname.ValueString(),
		IPAddress:        v.IPAddress.ValueString(),
		AllowedCIDRs:     &ipList,
		SharedSubnetCIDR: v.SharedSubnetCIDR.ValueStringPointer(),
		SharedSubnetID:   v.SharedSubnetCIDR.ValueString(),
		SharedNetworkID:  v.SharedNetworkID.ValueString(),
	}, diags
}

func (v PrivateNetworkConfigValueV2) ToDBaaSSdkRequest(ctx context.Context) (*sys11dbaassdk.PSQLPrivateNetworkConfigRequestV2, diag.Diagnostics) {
	var diags diag.Diagnostics

	var ipList []string

	for _, e := range v.AllowedCIDRs.Elements() {
		ipList = append(ipList, strings.Trim(e.String(), "\""))
	}

	return &sys11dbaassdk.PSQLPrivateNetworkConfigRequestV2{
		Enabled:          v.Enabled.ValueBool(),
		AllowedCIDRs:     &ipList,
		SharedSubnetCIDR: v.SharedSubnetCIDR.ValueStringPointer(),
	}, diags
}

// PublicNetworkConfig

var _ basetypes.ObjectTypable = PublicNetworkConfigTypeV2{}

type PublicNetworkConfigTypeV2 struct {
	basetypes.ObjectType
}

func (t PublicNetworkConfigTypeV2) Equal(o attr.Type) bool {
	other, ok := o.(PublicNetworkConfigTypeV2)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t PublicNetworkConfigTypeV2) String() string {
	return "PublicNetworkConfigTypeV2"
}

func (t PublicNetworkConfigTypeV2) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributes := in.Attributes()

	enabledAttribute, ok := attributes["enabled"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`enabled is missing from object`)

		return nil, diags
	}

	enabledVal, ok := enabledAttribute.(basetypes.BoolValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`enabled expected to be basetypes.BoolValue, was: %T`, enabledAttribute))
	}

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

	allowedCIDRsAttribute, ok := attributes["allowed_cidrs"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`allowed_cidrs is missing from object`)

		return nil, diags
	}

	allowedCIDRsVal, ok := allowedCIDRsAttribute.(basetypes.ListValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`allowed_cidrs expected to be basetypes.ListValue, was: %T`, allowedCIDRsAttribute))
	}

	if diags.HasError() {
		return nil, diags
	}

	return PublicNetworkConfigValueV2{
		Enabled:      enabledVal,
		Hostname:     hostnameVal,
		IPAddress:    ipAddressVal,
		AllowedCIDRs: allowedCIDRsVal,
		state:        attr.ValueStateKnown,
	}, diags
}

func NewPublicNetworkConfigValueV2Null() PublicNetworkConfigValueV2 {
	return PublicNetworkConfigValueV2{
		state: attr.ValueStateNull,
	}
}

func NewPublicNetworkConfigValueV2Unknown() PublicNetworkConfigValueV2 {
	return PublicNetworkConfigValueV2{
		state: attr.ValueStateUnknown,
	}
}

func NewPublicNetworkConfigValueV2(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (PublicNetworkConfigValueV2, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]

		if !ok {
			diags.AddError(
				"Missing PublicNetworkConfigValueV2 Attribute Value",
				"While creating a PublicNetworkConfigValueV2 value, a missing attribute value was detected. "+
					"A PublicNetworkConfigValueV2 must contain values for all attributes, even if null or unknown. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("PublicNetworkConfigValueV2 Attribute Name (%s) Expected Type: %s", name, attributeType.String()),
			)

			continue
		}

		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError(
				"Invalid PublicNetworkConfigValueV2 Attribute Type",
				"While creating a PublicNetworkConfigValueV2 value, an invalid attribute value was detected. "+
					"A PublicNetworkConfigValueV2 must use a matching attribute type for the value. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("PublicNetworkConfigValueV2 Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
					fmt.Sprintf("PublicNetworkConfigValueV2 Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
			)
		}
	}

	for name := range attributes {
		_, ok := attributeTypes[name]

		if !ok {
			diags.AddError(
				"Extra PublicNetworkConfigValueV2 Attribute Value",
				"While creating a PublicNetworkConfigValueV2 value, an extra attribute value was detected. "+
					"A PublicNetworkConfigValueV2 must not contain values beyond the expected attribute types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Extra PublicNetworkConfigValueV2 Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewPublicNetworkConfigValueV2Unknown(), diags
	}

	enabledAttribute, ok := attributes["enabled"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`enabled is missing from object`)

		return NewPublicNetworkConfigValueV2Unknown(), diags
	}

	enabledVal, ok := enabledAttribute.(basetypes.BoolValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`retention expected to be basetypes.BoolValue, was: %T`, enabledVal))
	}

	hostnameAttribute, ok := attributes["hostname"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`hostname is missing from object`)

		return NewPublicNetworkConfigValueV2Unknown(), diags
	}

	hostnameVal, ok := hostnameAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`hostname expected to be basetypes.StringValue, was: %T`, hostnameAttribute))
	}

	ipAddressAttribute, ok := attributes["ip_address"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`ip_address is missing from object`)

		return NewPublicNetworkConfigValueV2Unknown(), diags
	}

	ipAddressVal, ok := ipAddressAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`ip_address expected to be basetypes.StringValue, was: %T`, ipAddressAttribute))
	}

	allowedCIDRsAttribute, ok := attributes["allowed_cidrs"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`allowed_cidrs is missing from object`)

		return NewPublicNetworkConfigValueV2Unknown(), diags
	}

	allowedCIDRsVal, ok := allowedCIDRsAttribute.(basetypes.ListValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`allowed_cidrs expected to be basetypes.ListValue, was: %T`, allowedCIDRsAttribute))
	}

	if diags.HasError() {
		return NewPublicNetworkConfigValueV2Unknown(), diags
	}

	return PublicNetworkConfigValueV2{
		Enabled:      enabledVal,
		AllowedCIDRs: allowedCIDRsVal,
		Hostname:     hostnameVal,
		IPAddress:    ipAddressVal,
		state:        attr.ValueStateKnown,
	}, diags
}

func NewPublicNetworkConfigValueV2Must(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) PublicNetworkConfigValueV2 {
	object, diags := NewPublicNetworkConfigValueV2(attributeTypes, attributes)

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

		panic("NewPublicNetworkConfigValueV2Must received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}

func (t PublicNetworkConfigTypeV2) ValueFromTerraformV2(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewPublicNetworkConfigValueV2Null(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewPublicNetworkConfigValueV2Unknown(), nil
	}

	if in.IsNull() {
		return NewPublicNetworkConfigValueV2Null(), nil
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

	return NewPublicNetworkConfigValueV2Must(PublicNetworkConfigValueV2{}.AttributeTypes(ctx), attributes), nil
}

func (t PublicNetworkConfigTypeV2) ValueType(ctx context.Context) attr.Value {
	return PublicNetworkConfigValueV2{}
}

var _ basetypes.ObjectValuable = PublicNetworkConfigValueV2{}

type PublicNetworkConfigValueV2 struct {
	Enabled      basetypes.BoolValue   `tfsdk:"enabled"`
	AllowedCIDRs basetypes.ListValue   `tfsdk:"allowed_cidrs"`
	Hostname     basetypes.StringValue `tfsdk:"hostname"`
	IPAddress    basetypes.StringValue `tfsdk:"ip_address"`
	state        attr.ValueState
}

func (v PublicNetworkConfigValueV2) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 2)

	var val tftypes.Value
	var err error

	attrTypes["enabled"] = basetypes.BoolType{}.TerraformType(ctx)
	attrTypes["hostname"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["ip_address"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["allowed_cidrs"] = basetypes.ListType{
		ElemType: types.StringType,
	}.TerraformType(ctx)

	objectType := tftypes.Object{AttributeTypes: attrTypes}

	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 2)

		val, err = v.Enabled.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["enabled"] = val

		val, err = v.Hostname.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["hostname"] = val

		val, err = v.IPAddress.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["ip_address"] = val

		val, err = v.AllowedCIDRs.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["allowed_cidrs"] = val

		return tftypes.NewValue(objectType, vals), nil
	case attr.ValueStateNull:
		return tftypes.NewValue(objectType, nil), nil
	case attr.ValueStateUnknown:
		return tftypes.NewValue(objectType, tftypes.UnknownValue), nil
	default:
		panic(fmt.Sprintf("unhandled Object state in ToTerraformValue: %s", v.state))
	}
}

func (v PublicNetworkConfigValueV2) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v PublicNetworkConfigValueV2) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v PublicNetworkConfigValueV2) String() string {
	return "PublicNetworkConfigValueV2"
}

func (v PublicNetworkConfigValueV2) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	objVal, diags := types.ObjectValue(
		map[string]attr.Type{
			"enabled":    basetypes.BoolType{},
			"hostname":   basetypes.StringType{},
			"ip_address": basetypes.StringType{},
			"allowed_cidrs": basetypes.ListType{
				ElemType: types.StringType,
			},
		},
		map[string]attr.Value{
			"enabled":       v.Enabled,
			"hostname":      v.Hostname,
			"ip_address":    v.IPAddress,
			"allowed_cidrs": v.AllowedCIDRs,
		})

	return objVal, diags
}

func (v PublicNetworkConfigValueV2) Equal(o attr.Value) bool {
	other, ok := o.(PublicNetworkConfigValueV2)

	if !ok {
		return false
	}

	if v.state != other.state {
		return false
	}

	if v.state != attr.ValueStateKnown {
		return true
	}

	if !v.Enabled.Equal(other.Enabled) {
		return false
	}

	if !v.Hostname.Equal(other.Hostname) {
		return false
	}

	if !v.IPAddress.Equal(other.IPAddress) {
		return false
	}

	if !v.AllowedCIDRs.Equal(other.AllowedCIDRs) {
		return false
	}

	return true
}

func (v PublicNetworkConfigValueV2) Type(ctx context.Context) attr.Type {
	return PublicNetworkConfigTypeV2{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v PublicNetworkConfigValueV2) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"enabled":       basetypes.BoolType{},
		"hostname":      basetypes.StringType{},
		"ip_address":    basetypes.StringType{},
		"allowed_cidrs": basetypes.ListType{ElemType: types.StringType},
	}
}

func (v PublicNetworkConfigValueV2) ToDBaaSSdkResponse(ctx context.Context) (*sys11dbaassdk.PSQLPublicNetworkConfigResponseV2, diag.Diagnostics) {
	var diags diag.Diagnostics

	var ipList []string

	for _, e := range v.AllowedCIDRs.Elements() {
		ipList = append(ipList, strings.Trim(e.String(), "\""))
	}

	return &sys11dbaassdk.PSQLPublicNetworkConfigResponseV2{
		Enabled:      v.Enabled.ValueBool(),
		Hostname:     v.Hostname.ValueString(),
		IPAddress:    v.IPAddress.ValueString(),
		AllowedCIDRs: &ipList,
	}, diags
}

func (v PublicNetworkConfigValueV2) ToDBaaSSdkRequest(ctx context.Context) (*sys11dbaassdk.PSQLPublicNetworkConfigRequestV2, diag.Diagnostics) {
	var diags diag.Diagnostics

	var ipList []string

	for _, e := range v.AllowedCIDRs.Elements() {
		ipList = append(ipList, strings.Trim(e.String(), "\""))
	}

	return &sys11dbaassdk.PSQLPublicNetworkConfigRequestV2{
		Enabled:      v.Enabled.ValueBool(),
		AllowedCIDRs: &ipList,
	}, diags
}
