package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/kraihn/terraform-provider-mailcow/internal/client"
	"os"
)

func New() tfsdk.Provider {
	return &provider{}
}

type provider struct {
	configured bool
	client     client.Client
	host       string
	apikey     string
}

func (p *provider) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"host": {
				Type:                types.StringType,
				Description:         "The Mailcow server for accessing the API. Can be sourced from MAILCOW_HOST.",
				MarkdownDescription: "The Mailcow server for accessing the API. Can be sourced from `MAILCOW_HOST`.",
				Optional:            true,
				Computed:            true,
			},
			"apikey": {
				Type:                types.StringType,
				Description:         "The Mailcow API Key for accessing the API. Can be sourced from MAILCOW_APIKEY.",
				MarkdownDescription: "The Mailcow API Key for accessing the API. Can be sourced from `MAILCOW_APIKEY`.",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
			},
		},
	}, nil
}

type providerData struct {
	Host   types.String `tfsdk:"host"`
	ApiKey types.String `tfsdk:"apikey"`
}

func (p *provider) Configure(ctx context.Context, req tfsdk.ConfigureProviderRequest, resp *tfsdk.ConfigureProviderResponse) {
	var config providerData
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var host string
	if config.Host.Unknown {
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Cannot use unknown value as host",
		)
	}

	if config.Host.Null {
		host = os.Getenv("MAILCOW_HOST")
	} else {
		host = config.Host.Value
	}

	if host == "" {
		resp.Diagnostics.AddError(
			"Unable to find host",
			"HostURL cannot be an empty string",
		)
		return
	}

	var apiKey string
	if config.ApiKey.Unknown {
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Cannot use unknown value as apiKey",
		)
	}

	if config.ApiKey.Null {
		apiKey = os.Getenv("MAILCOW_APIKEY")
	} else {
		apiKey = config.ApiKey.Value
	}

	if apiKey == "" {
		resp.Diagnostics.AddError(
			"Unable to find apiKey",
			"Username cannot be an empty string",
		)
		return
	}

	c, _ := client.NewClient(&config.Host.Value, &config.ApiKey.Value)

	p.client = *c
	p.configured = true
}

func (p *provider) GetResources(ctx context.Context) (map[string]tfsdk.ResourceType, diag.Diagnostics) {
	return map[string]tfsdk.ResourceType{
		"mailcow_alias":   resourceAliasType{},
		"mailcow_domain":  resourceDomainType{},
	}, nil
}

func (p *provider) GetDataSources(ctx context.Context) (map[string]tfsdk.DataSourceType, diag.Diagnostics) {
	return map[string]tfsdk.DataSourceType{
		"mailcow_all_aliases":   allAliasesDataSourceType{},
		"mailcow_all_domains":   allDomainsDataSourceType{},
		"mailcow_all_mailboxes": allMailboxesDataSourceType{},
		"mailcow_domain":        domainDataSourceType{},
		"mailcow_mailbox":       mailboxDataSourceType{},
	}, nil
}
