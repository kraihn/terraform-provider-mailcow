package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type allMailboxesDataSourceType struct{}

func (t allMailboxesDataSourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"mailboxes": {
				Computed: true,
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"email": {
						Type:     types.StringType,
						Computed: true,
					},
					"username": {
						Type:     types.StringType,
						Computed: true,
					},
					"domain": {
						Type:     types.StringType,
						Computed: true,
					},
					"name": {
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

func (r allMailboxesDataSourceType) NewDataSource(ctx context.Context, p tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	return allmailboxDataSource{
		p: *(p.(*provider)),
	}, nil
}

type allmailboxDataSourceData struct {
	Mailboxes []allmailboxItem `tfsdk:"mailboxes"`
}

type allmailboxItem struct {
	Active   types.Bool   `tfsdk:"active"`
	Domain   types.String `tfsdk:"domain"`
	Email    types.String `tfsdk:"email"`
	Name     types.String `tfsdk:"name"`
	Username types.String `tfsdk:"username"`
}

type allmailboxDataSource struct {
	p provider
}

func (d allmailboxDataSource) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	var data allmailboxDataSourceData
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	mailboxes, err := d.p.client.GetAllMailboxes()
	if err != nil {
		resp.Diagnostics.AddError("Client Error - Get All Mailboxes", fmt.Sprintf("Unable to read, got error: %s", err))
		return
	}

	for _, mailbox := range *mailboxes {
		m := allmailboxItem{
			Active:   types.Bool{Value: mailbox.Active == 1},
			Domain:   types.String{Value: mailbox.Domain},
			Email:    types.String{Value: mailbox.Email},
			Name:     types.String{Value: mailbox.Name},
			Username: types.String{Value: mailbox.Username},
		}

		data.Mailboxes = append(data.Mailboxes, m)
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
