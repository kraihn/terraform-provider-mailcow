package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type domainDataSourceType struct{}

func (t domainDataSourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"domain": {
				Type:                types.StringType,
				Description:         "The @domain.tld part of the email address",
				MarkdownDescription: "",
				Required:            true,
			},
			"description": {
				Type:     types.StringType,
				Computed: true,
			},
			"active": {
				Type:     types.BoolType,
				Computed: true,
			},
		},
	}, nil
}

func (r domainDataSourceType) NewDataSource(ctx context.Context, p tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	return domainDataSource{
		p: *(p.(*provider)),
	}, nil
}

type domainDataSourceData struct {
	Active      types.Bool   `tfsdk:"active"`
	Description types.String `tfsdk:"description"`
	Domain      types.String `tfsdk:"domain"`
}

type domainDataSource struct {
	p provider
}

func (d domainDataSource) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	var data domainDataSourceData
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	domain, err := d.p.client.GetDomain(data.Domain.Value)
	if err != nil {
		resp.Diagnostics.AddError("Client Error - Get Domain", fmt.Sprintf("Unable to read, got error: %s", err))
		return
	}

	data.Active = types.Bool{Value: domain.Active == 1}
	data.Description = types.String{Value: domain.Description}
	data.Domain = types.String{Value: domain.Name}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
