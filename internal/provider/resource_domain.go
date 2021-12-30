package provider

import (
	"bytes"
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/kraihn/terraform-provider-mailcow/internal/plan_modifiers"
	"log"
	"net/http"
	"strconv"
)

type resourceDomainType struct{}

func (r resourceDomainType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"domain": {
				Type:     types.StringType,
				Required: true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					tfsdk.RequiresReplace(),
				},
			},
			"description": {
				Type:     types.StringType,
				Required: true,
			},
			"active": {
				Type:     types.BoolType,
				Optional: true,
				Computed: true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					plan_modifiers.DefaultBool(true),
				},
			},
			"quota": {
				Type:     types.Int64Type,
				Required: true,
			},
			"mailboxes": {
				Type:     types.Int64Type,
				Required: true,
			},
			"mailbox_default_size": {
				Type:     types.Int64Type,
				Required: true,
			},
			"mailbox_max_size": {
				Type:     types.Int64Type,
				Required: true,
			},
			"aliases": {
				Type:     types.Int64Type,
				Required: true,
			},
		},
	}, nil
}

func (r resourceDomainType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return resourceDomain{
		p: *(p.(*provider)),
	}, nil
}

type resourceDomain struct {
	p provider
}

func (r resourceDomain) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}

	// Retrieve values from plan
	var plan Domain
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//Generate the URL to access
	url := r.p.client.HostURL + "/api/v1/add/domain"

	active := "0"
	if plan.Active.Value {
		active = "1"
	}

	data := []byte(`{
		"active": "` + active + `",
		"aliases": "` + strconv.FormatInt(plan.Aliases.Value, 10) + `",
		"defquota": "` + strconv.FormatInt(plan.MailboxDefaultSizeMB.Value, 10) + `",
		"description": "` + plan.Description.Value + `",
		"domain": "` + plan.Domain.Value + `",
		"mailboxes": "` + strconv.FormatInt(plan.Mailboxes.Value, 10) + `",
		"maxquota": "` + strconv.FormatInt(plan.MailboxMaxSizeMB.Value, 10) + `",
		"quota": "` + strconv.FormatInt(plan.QuotaMB.Value, 10) + `"
	}`)

	request, _ := http.NewRequest("POST", url, bytes.NewBuffer(data))
	_, err := r.p.client.DoRequest(request)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read, got error: %s", err))
		return
	}

	var result = plan

	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r resourceDomain) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var state Domain
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	domain, err := r.p.client.GetDomain(state.Domain.Value)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read, got error: %s", err))
		return
	}

	state.Domain = types.String{Value: domain.Name}
	state.Description = types.String{Value: domain.Description}
	state.Active = types.Bool{Value: domain.Active == 1}
	state.QuotaMB = types.Int64{Value: domain.QuotaBytes / 1024 / 1024}
	state.Mailboxes = types.Int64{Value: domain.Mailboxes}
	state.MailboxDefaultSizeMB = types.Int64{Value: domain.MailboxDefaultSizeBytes / 1024 / 1024}
	state.MailboxMaxSizeMB = types.Int64{Value: domain.MailboxMaxSizeBytes / 1024 / 1024}
	state.Aliases = types.Int64{Value: domain.Aliases}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r resourceDomain) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	// Get plan values
	var plan Domain
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//Generate the URL to access
	url := r.p.client.HostURL + "/api/v1/edit/domain"

	active := "0"
	if plan.Active.Value {
		active = "1"
	}

	var jsonData = []byte(`{
		"attr": {
			"active": "` + active + `",
			"aliases": "` + strconv.FormatInt(plan.Aliases.Value, 10) + `",
			"defquota": "` + strconv.FormatInt(plan.MailboxDefaultSizeMB.Value, 10) + `",
			"description": "` + plan.Description.Value + `",
			"mailboxes": "` + strconv.FormatInt(plan.Mailboxes.Value, 10) + `",
			"maxquota": "` + strconv.FormatInt(plan.MailboxMaxSizeMB.Value, 10) + `",
			"quota": "` + strconv.FormatInt(plan.QuotaMB.Value, 10) + `"
		},
		"items": ["` + plan.Domain.Value + `"]
	}`)

	log.Printf("out %s", string(jsonData))

	request, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

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

func (r resourceDomain) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var state Domain
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.p.client.DeleteDomain(state.Domain.Value)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read, got error: %s", err))
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r resourceDomain) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("domain_name"), req, resp)
}
