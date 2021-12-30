package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type mailboxDataSourceType struct{}

func (t mailboxDataSourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"email": {
				Type:     types.StringType,
				Required: true,
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
		},
	}, nil
}

func (r mailboxDataSourceType) NewDataSource(ctx context.Context, p tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	return mailboxDataSource{
		p: *(p.(*provider)),
	}, nil
}

type mailboxDataSourceData struct {
	Active   types.Bool   `tfsdk:"active"`
	Domain   types.String `tfsdk:"domain"`
	Email    types.String `tfsdk:"email"`
	Name     types.String `tfsdk:"name"`
	Username types.String `tfsdk:"username"`
}

type mailboxDataSource struct {
	p provider
}

func (d mailboxDataSource) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	var data mailboxDataSourceData
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	mailbox, err := d.p.client.GetMailbox(data.Email.Value)
	if err != nil {
		resp.Diagnostics.AddError("Client Error - Get Mailbox", fmt.Sprintf("Unable to read, got error: %s", err))
		return
	}

	data.Active = types.Bool{Value: mailbox.Active == 1}
	data.Domain = types.String{Value: mailbox.Domain}
	data.Name = types.String{Value: mailbox.Name}
	data.Username = types.String{Value: mailbox.Username}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
