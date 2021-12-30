package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
)

type allAliasesDataSourceType struct{}

func (t allAliasesDataSourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"aliases": {
				Computed: true,
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"id": {
						Type:     types.Int64Type,
						Computed: true,
					},
					"alias": {
						Type:     types.StringType,
						Computed: true,
					},
					"goto_addresses": {
						Type: types.ListType{
							ElemType: types.StringType,
						},
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

func (r allAliasesDataSourceType) NewDataSource(ctx context.Context, p tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	return allAliasesDataSource{
		p: *(p.(*provider)),
	}, nil
}

type allAliasesDataSourceData struct {
	Aliases []allAliasItem `tfsdk:"aliases"`
}

type allAliasItem struct {
	Active        types.Bool     `tfsdk:"active"`
	Alias         types.String   `tfsdk:"alias"`
	GotoAddresses []types.String `tfsdk:"goto_addresses"`
	ID            types.Int64    `tfsdk:"id"`
}

type allAliasesDataSource struct {
	p provider
}

func (d allAliasesDataSource) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	var data allAliasesDataSourceData
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	aliases, err := d.p.client.GetAllAliases()
	if err != nil {
		resp.Diagnostics.AddError("Client Error - Get All Aliases", fmt.Sprintf("Unable to read, got error: %s", err))
		return
	}

	for _, alias := range *aliases {
		var destinations []types.String
		for _, destination := range strings.Split(alias.GoTo, ",") {
			destinations = append(destinations, types.String{Value: destination})
		}

		m := allAliasItem{
			Active:        types.Bool{Value: alias.Active == 1},
			Alias:         types.String{Value: alias.Address},
			GotoAddresses: destinations,
			ID:            types.Int64{Value: alias.ID},
		}

		data.Aliases = append(data.Aliases, m)
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
