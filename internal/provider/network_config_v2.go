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

var _ basetypes.ObjectTypable = PrivateNetworkingTypeV2{}

type PrivateNetworkingTypeV2 struct {
	basetypes.ObjectType
}

func (t PrivateNetworkingTypeV2) Equal(o attr.Type) bool {
	other, ok := o.(PrivateNetworkingTypeV2)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t PrivateNetworkingTypeV2) String() string {
	return "PrivateNetworkingTypeV2"
}

func (t PrivateNetworkingTypeV2) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
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

	return PrivateNetworkingValueV2{
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

func NewPrivateNetworkingValueV2Null() PrivateNetworkingValueV2 {
	return PrivateNetworkingValueV2{
		state: attr.ValueStateNull,
	}
}

func NewPrivateNetworkingValueV2Unknown() PrivateNetworkingValueV2 {
	return PrivateNetworkingValueV2{
		state: attr.ValueStateUnknown,
	}
}

func NewPrivateNetworkingValueV2(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (PrivateNetworkingValueV2, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]

		if !ok {
			diags.AddError(
				"Missing PrivateNetworkingValueV2 Attribute Value",
				"While creating a PrivateNetworkingValueV2 value, a missing attribute value was detected. "+
					"A PrivateNetworkingValueV2 must contain values for all attributes, even if null or unknown. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("PrivateNetworkingValueV2 Attribute Name (%s) Expected Type: %s", name, attributeType.String()),
			)

			continue
		}

		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError(
				"Invalid PrivateNetworkingValueV2 Attribute Type",
				"While creating a PrivateNetworkingValueV2 value, an invalid attribute value was detected. "+
					"A PrivateNetworkingValueV2 must use a matching attribute type for the value. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("PrivateNetworkingValueV2 Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
					fmt.Sprintf("PrivateNetworkingValueV2 Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
			)
		}
	}

	for name := range attributes {
		_, ok := attributeTypes[name]

		if !ok {
			diags.AddError(
				"Extra PrivateNetworkingValueV2 Attribute Value",
				"While creating a PrivateNetworkingValueV2 value, an extra attribute value was detected. "+
					"A PrivateNetworkingValueV2 must not contain values beyond the expected attribute types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Extra PrivateNetworkingValueV2 Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewPrivateNetworkingValueV2Unknown(), diags
	}

	enabledAttribute, ok := attributes["enabled"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`enabled is missing from object`)

		return NewPrivateNetworkingValueV2Unknown(), diags
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

		return NewPrivateNetworkingValueV2Unknown(), diags
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

		return NewPrivateNetworkingValueV2Unknown(), diags
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

		return NewPrivateNetworkingValueV2Unknown(), diags
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

		return NewPrivateNetworkingValueV2Unknown(), diags
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

		return NewPrivateNetworkingValueV2Unknown(), diags
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

		return NewPrivateNetworkingValueV2Unknown(), diags
	}

	sharedNetworkIDVal, ok := sharedNetworkIDAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`shared_network_id expected to be basetypes.StringValue, was: %T`, sharedNetworkIDAttribute))
	}

	if diags.HasError() {
		return NewPrivateNetworkingValueV2Unknown(), diags
	}

	return PrivateNetworkingValueV2{
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

func NewPrivateNetworkingValueV2Must(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) PrivateNetworkingValueV2 {
	object, diags := NewPrivateNetworkingValueV2(attributeTypes, attributes)

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

		panic("NewPrivateNetworkingValueV2Must received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}

func (t PrivateNetworkingTypeV2) ValueFromTerraformV2(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewPrivateNetworkingValueV2Null(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewPrivateNetworkingValueV2Unknown(), nil
	}

	if in.IsNull() {
		return NewPrivateNetworkingValueV2Null(), nil
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

	return NewPrivateNetworkingValueV2Must(PrivateNetworkingValueV2{}.AttributeTypes(ctx), attributes), nil
}

func (t PrivateNetworkingTypeV2) ValueType(ctx context.Context) attr.Value {
	return PrivateNetworkingValueV2{}
}

var _ basetypes.ObjectValuable = PrivateNetworkingValueV2{}

type PrivateNetworkingValueV2 struct {
	Enabled          basetypes.BoolValue   `tfsdk:"enabled"`
	AllowedCIDRs     basetypes.ListValue   `tfsdk:"allowed_cidrs"`
	SharedSubnetCIDR basetypes.StringValue `tfsdk:"shared_subnet_cidr"`
	Hostname         basetypes.StringValue `tfsdk:"hostname"`
	IPAddress        basetypes.StringValue `tfsdk:"ip_address"`
	SharedSubnetID   basetypes.StringValue `tfsdk:"shared_subnet_id"`
	SharedNetworkID  basetypes.StringValue `tfsdk:"shared_network_id"`
	state            attr.ValueState
}

func (v PrivateNetworkingValueV2) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
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

func (v PrivateNetworkingValueV2) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v PrivateNetworkingValueV2) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v PrivateNetworkingValueV2) String() string {
	return "PrivateNetworkingValueV2"
}

func (v PrivateNetworkingValueV2) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
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

func (v PrivateNetworkingValueV2) Equal(o attr.Value) bool {
	other, ok := o.(PrivateNetworkingValueV2)

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

func (v PrivateNetworkingValueV2) Type(ctx context.Context) attr.Type {
	return PrivateNetworkingTypeV2{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v PrivateNetworkingValueV2) AttributeTypes(ctx context.Context) map[string]attr.Type {
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

func (v PrivateNetworkingValueV2) ToDBaaSSdkResponse(ctx context.Context) (*sys11dbaassdk.PSQLPrivateNetworkingResponseV2, diag.Diagnostics) {
	var diags diag.Diagnostics

	var ipList []string

	for _, e := range v.AllowedCIDRs.Elements() {
		ipList = append(ipList, strings.Trim(e.String(), "\""))
	}

	return &sys11dbaassdk.PSQLPrivateNetworkingResponseV2{
		Enabled:          v.Enabled.ValueBool(),
		Hostname:         v.Hostname.ValueString(),
		IPAddress:        v.IPAddress.ValueString(),
		AllowedCIDRs:     &ipList,
		SharedSubnetCIDR: v.SharedSubnetCIDR.ValueStringPointer(),
		SharedSubnetID:   v.SharedSubnetCIDR.ValueString(),
		SharedNetworkID:  v.SharedNetworkID.ValueString(),
	}, diags
}

func (v PrivateNetworkingValueV2) ToDBaaSSdkRequest(ctx context.Context) (*sys11dbaassdk.PSQLPrivateNetworkingRequestV2, diag.Diagnostics) {
	var diags diag.Diagnostics

	var ipList []string

	for _, e := range v.AllowedCIDRs.Elements() {
		ipList = append(ipList, strings.Trim(e.String(), "\""))
	}

	return &sys11dbaassdk.PSQLPrivateNetworkingRequestV2{
		Enabled:          v.Enabled.ValueBool(),
		AllowedCIDRs:     &ipList,
		SharedSubnetCIDR: v.SharedSubnetCIDR.ValueStringPointer(),
	}, diags
}

// PublicNetworking

var _ basetypes.ObjectTypable = PublicNetworkingTypeV2{}

type PublicNetworkingTypeV2 struct {
	basetypes.ObjectType
}

func (t PublicNetworkingTypeV2) Equal(o attr.Type) bool {
	other, ok := o.(PublicNetworkingTypeV2)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t PublicNetworkingTypeV2) String() string {
	return "PublicNetworkingTypeV2"
}

func (t PublicNetworkingTypeV2) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
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

	return PublicNetworkingValueV2{
		Enabled:      enabledVal,
		Hostname:     hostnameVal,
		IPAddress:    ipAddressVal,
		AllowedCIDRs: allowedCIDRsVal,
		state:        attr.ValueStateKnown,
	}, diags
}

func NewPublicNetworkingValueV2Null() PublicNetworkingValueV2 {
	return PublicNetworkingValueV2{
		state: attr.ValueStateNull,
	}
}

func NewPublicNetworkingValueV2Unknown() PublicNetworkingValueV2 {
	return PublicNetworkingValueV2{
		state: attr.ValueStateUnknown,
	}
}

func NewPublicNetworkingValueV2(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (PublicNetworkingValueV2, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]

		if !ok {
			diags.AddError(
				"Missing PublicNetworkingValueV2 Attribute Value",
				"While creating a PublicNetworkingValueV2 value, a missing attribute value was detected. "+
					"A PublicNetworkingValueV2 must contain values for all attributes, even if null or unknown. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("PublicNetworkingValueV2 Attribute Name (%s) Expected Type: %s", name, attributeType.String()),
			)

			continue
		}

		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError(
				"Invalid PublicNetworkingValueV2 Attribute Type",
				"While creating a PublicNetworkingValueV2 value, an invalid attribute value was detected. "+
					"A PublicNetworkingValueV2 must use a matching attribute type for the value. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("PublicNetworkingValueV2 Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
					fmt.Sprintf("PublicNetworkingValueV2 Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
			)
		}
	}

	for name := range attributes {
		_, ok := attributeTypes[name]

		if !ok {
			diags.AddError(
				"Extra PublicNetworkingValueV2 Attribute Value",
				"While creating a PublicNetworkingValueV2 value, an extra attribute value was detected. "+
					"A PublicNetworkingValueV2 must not contain values beyond the expected attribute types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Extra PublicNetworkingValueV2 Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewPublicNetworkingValueV2Unknown(), diags
	}

	enabledAttribute, ok := attributes["enabled"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`enabled is missing from object`)

		return NewPublicNetworkingValueV2Unknown(), diags
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

		return NewPublicNetworkingValueV2Unknown(), diags
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

		return NewPublicNetworkingValueV2Unknown(), diags
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

		return NewPublicNetworkingValueV2Unknown(), diags
	}

	allowedCIDRsVal, ok := allowedCIDRsAttribute.(basetypes.ListValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`allowed_cidrs expected to be basetypes.ListValue, was: %T`, allowedCIDRsAttribute))
	}

	if diags.HasError() {
		return NewPublicNetworkingValueV2Unknown(), diags
	}

	return PublicNetworkingValueV2{
		Enabled:      enabledVal,
		AllowedCIDRs: allowedCIDRsVal,
		Hostname:     hostnameVal,
		IPAddress:    ipAddressVal,
		state:        attr.ValueStateKnown,
	}, diags
}

func NewPublicNetworkingValueV2Must(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) PublicNetworkingValueV2 {
	object, diags := NewPublicNetworkingValueV2(attributeTypes, attributes)

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

		panic("NewPublicNetworkingValueV2Must received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}

func (t PublicNetworkingTypeV2) ValueFromTerraformV2(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewPublicNetworkingValueV2Null(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewPublicNetworkingValueV2Unknown(), nil
	}

	if in.IsNull() {
		return NewPublicNetworkingValueV2Null(), nil
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

	return NewPublicNetworkingValueV2Must(PublicNetworkingValueV2{}.AttributeTypes(ctx), attributes), nil
}

func (t PublicNetworkingTypeV2) ValueType(ctx context.Context) attr.Value {
	return PublicNetworkingValueV2{}
}

var _ basetypes.ObjectValuable = PublicNetworkingValueV2{}

type PublicNetworkingValueV2 struct {
	Enabled      basetypes.BoolValue   `tfsdk:"enabled"`
	AllowedCIDRs basetypes.ListValue   `tfsdk:"allowed_cidrs"`
	Hostname     basetypes.StringValue `tfsdk:"hostname"`
	IPAddress    basetypes.StringValue `tfsdk:"ip_address"`
	state        attr.ValueState
}

func (v PublicNetworkingValueV2) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
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

func (v PublicNetworkingValueV2) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v PublicNetworkingValueV2) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v PublicNetworkingValueV2) String() string {
	return "PublicNetworkingValueV2"
}

func (v PublicNetworkingValueV2) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
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

func (v PublicNetworkingValueV2) Equal(o attr.Value) bool {
	other, ok := o.(PublicNetworkingValueV2)

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

func (v PublicNetworkingValueV2) Type(ctx context.Context) attr.Type {
	return PublicNetworkingTypeV2{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v PublicNetworkingValueV2) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"enabled":       basetypes.BoolType{},
		"hostname":      basetypes.StringType{},
		"ip_address":    basetypes.StringType{},
		"allowed_cidrs": basetypes.ListType{ElemType: types.StringType},
	}
}

func (v PublicNetworkingValueV2) ToDBaaSSdkResponse(ctx context.Context) (*sys11dbaassdk.PSQLPublicNetworkingResponseV2, diag.Diagnostics) {
	var diags diag.Diagnostics

	var ipList []string

	for _, e := range v.AllowedCIDRs.Elements() {
		ipList = append(ipList, strings.Trim(e.String(), "\""))
	}

	return &sys11dbaassdk.PSQLPublicNetworkingResponseV2{
		Enabled:      v.Enabled.ValueBool(),
		Hostname:     v.Hostname.ValueString(),
		IPAddress:    v.IPAddress.ValueString(),
		AllowedCIDRs: &ipList,
	}, diags
}

func (v PublicNetworkingValueV2) ToDBaaSSdkRequest(ctx context.Context) (*sys11dbaassdk.PSQLPublicNetworkingRequestV2, diag.Diagnostics) {
	var diags diag.Diagnostics

	var ipList []string

	for _, e := range v.AllowedCIDRs.Elements() {
		ipList = append(ipList, strings.Trim(e.String(), "\""))
	}

	return &sys11dbaassdk.PSQLPublicNetworkingRequestV2{
		Enabled:      v.Enabled.ValueBool(),
		AllowedCIDRs: &ipList,
	}, diags
}
