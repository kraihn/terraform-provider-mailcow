package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

type Alias struct {
	Active        types.Bool     `tfsdk:"active"`
	Alias         types.String   `tfsdk:"alias"`
	GotoAddresses []types.String `tfsdk:"goto_addresses"`
	ID            types.Int64    `tfsdk:"id"`
}

type Domain struct {
	Active               types.Bool   `tfsdk:"active"`
	Aliases              types.Int64  `tfsdk:"aliases"`
	Description          types.String `tfsdk:"description"`
	Domain               types.String `tfsdk:"domain"`
	MailboxDefaultSizeMB types.Int64  `tfsdk:"mailbox_default_size"`
	MailboxMaxSizeMB     types.Int64  `tfsdk:"mailbox_max_size"`
	Mailboxes            types.Int64  `tfsdk:"mailboxes"`
	Password             types.String `tfsdk:"password"`
	QuotaMB              types.Int64  `tfsdk:"quota"`
}
