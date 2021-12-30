package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/kraihn/terraform-provider-mailcow/internal/plan_modifiers"
	"github.com/kraihn/terraform-provider-mailcow/internal/validators"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type resourceAliasType struct{}

func (r resourceAliasType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:     types.Int64Type,
				Computed: true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					tfsdk.UseStateForUnknown(),
				},
			},
			"alias": {
				Type:     types.StringType,
				Required: true,
			},
			"goto_addresses": {
				Type: types.ListType{
					ElemType: types.StringType,
				},
				Required: true,
				Validators: []tfsdk.AttributeValidator{
					validators.ListNotEmptyValidator{},
				},
			},
			"active": {
				Type:     types.BoolType,
				Optional: true,
				Computed: true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					plan_modifiers.DefaultBool(true),
				},
			},
		},
	}, nil
}

func (r resourceAliasType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return resourceAlias{
		p: *(p.(*provider)),
	}, nil
}

type resourceAlias struct {
	p provider
}

func (r resourceAlias) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}

	var plan Alias
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := r.p.client.HostURL + "/api/v1/add/alias"

	active := "0"
	if plan.Active.Value {
		active = "1"
	}

	var destinations []string
	for _, destination := range plan.GotoAddresses {
		destinations = append(destinations, destination.Value)
	}

	data := []byte(`{
		"active": "` + active + `",
		"address": "` + plan.Alias.Value + `",
		"goto": "` + strings.Join(destinations, ",") + `"
	}`)

	request, _ := http.NewRequest("POST", url, bytes.NewBuffer(data))
	value, err := r.p.client.DoRequest(request)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read, got error: %s", err))
		return
	}

	type test struct {
		Message []string `json:"msg"`
	}
	var outcome []test
	json.Unmarshal(value, &outcome)

	id, _ := strconv.ParseInt(outcome[0].Message[2], 10, 64)
	result := Alias{
		ID:            types.Int64{Value: id},
		Alias:         plan.Alias,
		GotoAddresses: plan.GotoAddresses,
		Active:        plan.Active,
	}

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r resourceAlias) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var state Alias
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	alias, err := r.p.client.GetAlias(state.ID.Value)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read, got error: %s", err))
		return
	}

	log.Printf("arg %s", state.ID.Value)
	var destinations []types.String
	for _, destination := range strings.Split(alias.GoTo, ",") {
		destinations = append(destinations, types.String{Value: destination})
	}

	state.ID = types.Int64{Value: state.ID.Value}
	state.Alias = types.String{Value: alias.Address}
	state.GotoAddresses = destinations
	state.Active = types.Bool{Value: alias.Active == 1}
	//state.Active = types.Bool{Value: alias.Active}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r resourceAlias) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var plan Alias
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := r.p.client.HostURL + "/api/v1/edit/alias"

	id := strconv.FormatInt(plan.ID.Value, 10)
	active := "0"
	if plan.Active.Value {
		active = "1"
	}

	var destinations []string
	for _, destination := range plan.GotoAddresses {
		destinations = append(destinations, destination.Value)
	}

	data := []byte(`{
		"attr": {
			"active": "` + active + `",
			"address": "` + plan.Alias.Value + `",
			"goto": "` + strings.Join(destinations, ",") + `"
		},
		"items": [
			"` + id + `"
		]
	}`)

	request, _ := http.NewRequest("POST", url, bytes.NewBuffer(data))
	_, err := r.p.client.DoRequest(request)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read, got error: %s", err))
		return
	}

	result := plan

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r resourceAlias) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var state Alias
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.p.client.DeleteAlias(state.ID.Value)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read, got error: %s", err))
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r resourceAlias) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)
}
