package validators

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ListNotEmptyValidator struct {
}

func (v ListNotEmptyValidator) Description(ctx context.Context) string {
	return fmt.Sprintf("list length must be greater than 0")
}

func (v ListNotEmptyValidator) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("list length must be greater than 0")
}

func (v ListNotEmptyValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
	var list types.List
	diags := tfsdk.ValueAs(ctx, req.AttributeConfig, &list)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	if list.Unknown || list.Null {
		return
	}

	listLength := len(list.Elems)

	if listLength < 1 {
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			"Invalid List Length",
			fmt.Sprintf("List length must be greater than 0, got: %d.", listLength),
		)

		return
	}
}
