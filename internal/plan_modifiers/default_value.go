package plan_modifiers

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

type DefaultValue struct {
	Value attr.Value
}

func DefaultBool(v bool) DefaultValue {
	return DefaultValue{Value: types.Bool{Value: v}}
}

func DefaultInt64(v int64) DefaultValue {
	return DefaultValue{Value: types.Int64{Value: v}}
}

func (d DefaultValue) Description(ctx context.Context) string {
	return ""
}

func (d DefaultValue) MarkdownDescription(ctx context.Context) string {
	return ""
}

func (d DefaultValue) Modify(ctx context.Context, req tfsdk.ModifyAttributePlanRequest, res *tfsdk.ModifyAttributePlanResponse) {
	result, _, _ := tftypes.WalkAttributePath(req.Config.Raw, tftypes.NewAttributePathWithSteps(req.AttributePath.Steps()))

	isNull := false
	if result.(tftypes.Value).IsNull() {
		isNull = true
	} else if result.(tftypes.Value).Type().Is(tftypes.List{}) {
		if req.AttributeConfig.(types.List).Null {
			isNull = true
		}
	} else if result.(tftypes.Value).Type().Is(tftypes.Map{}) {
		if req.AttributeConfig.(types.Map).Null {
			isNull = true
		}
	} else if result.(tftypes.Value).Type().Is(tftypes.Set{}) {
		if req.AttributeConfig.(types.Set).Null {
			isNull = true
		}
	}

	if isNull {
		res.AttributePlan = d.Value
	}
}
