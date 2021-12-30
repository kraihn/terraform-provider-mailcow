package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type allDomainsDataSourceType struct{}

func (t allDomainsDataSourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"domains": {
				Computed: true,
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"domain": {
						Type:     types.StringType,
						Computed: true,
					},
					"description": {
						Type:     types.StringType,
						Computed: true,
					},
					"active": {
						Type:     types.BoolType,
						Computed: true,
					},
				}, tfsdk.ListNestedAttributesOptions{}),
			},
		},
	}, nil
}

func (r allDomainsDataSourceType) NewDataSource(ctx context.Context, p tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	return alldomainDataSource{
		p: *(p.(*provider)),
	}, nil
}

type alldomainDataSourceData struct {
	Domains []alldomainItem `tfsdk:"domains"`
}

type alldomainItem struct {
	Active      types.Bool   `tfsdk:"active"`
	Description types.String `tfsdk:"description"`
	DomainName  types.String `tfsdk:"domain"`
}

type alldomainDataSource struct {
	p provider
}

func (d alldomainDataSource) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	var data alldomainDataSourceData
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	domains, err := d.p.client.GetAllDomains()
	if err != nil {
		resp.Diagnostics.AddError("Client Error - Get All Domains", fmt.Sprintf("Unable to read, got error: %s", err))
		return
	}

	for _, domain := range *domains {
		d := alldomainItem{
			Active:      types.Bool{Value: domain.Active == 1},
			Description: types.String{Value: domain.Description},
			DomainName:  types.String{Value: domain.Name},
		}

		data.Domains = append(data.Domains, d)
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
