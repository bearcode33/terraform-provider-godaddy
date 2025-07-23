// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"

	"github.com/bearcode33/terraform-provider-godaddy/internal/godaddy"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure GoDaddyProvider satisfies various provider interfaces.
var _ provider.Provider = &GoDaddyProvider{}

// GoDaddyProvider defines the provider implementation.
type GoDaddyProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// GoDaddyProviderModel describes the provider data model.
type GoDaddyProviderModel struct {
	APIKey      types.String `tfsdk:"api_key"`
	APISecret   types.String `tfsdk:"api_secret"`
	Environment types.String `tfsdk:"environment"`
}

func (p *GoDaddyProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "godaddy"
	resp.Version = p.version
}

func (p *GoDaddyProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Terraform provider for managing GoDaddy domains and DNS records.",
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				MarkdownDescription: "GoDaddy API Key. Can also be set via GODADDY_API_KEY environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
			"api_secret": schema.StringAttribute{
				MarkdownDescription: "GoDaddy API Secret. Can also be set via GODADDY_API_SECRET environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
			"environment": schema.StringAttribute{
				MarkdownDescription: "GoDaddy API environment. Valid values are 'production' (default) and 'test'. Can also be set via GODADDY_ENVIRONMENT environment variable.",
				Optional:            true,
			},
		},
	}
}

func (p *GoDaddyProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data GoDaddyProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get API credentials from config or environment
	apiKey := os.Getenv("GODADDY_API_KEY")
	if !data.APIKey.IsNull() {
		apiKey = data.APIKey.ValueString()
	}

	apiSecret := os.Getenv("GODADDY_API_SECRET")
	if !data.APISecret.IsNull() {
		apiSecret = data.APISecret.ValueString()
	}

	environment := os.Getenv("GODADDY_ENVIRONMENT")
	if !data.Environment.IsNull() {
		environment = data.Environment.ValueString()
	}

	// Validate required fields
	if apiKey == "" {
		resp.Diagnostics.AddError(
			"Missing API Key",
			"The provider requires an API key to be configured. "+
				"Set the api_key attribute in the provider configuration or the GODADDY_API_KEY environment variable.",
		)
	}

	if apiSecret == "" {
		resp.Diagnostics.AddError(
			"Missing API Secret",
			"The provider requires an API secret to be configured. "+
				"Set the api_secret attribute in the provider configuration or the GODADDY_API_SECRET environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create GoDaddy client
	var opts []godaddy.ClientOption
	if environment == "test" {
		opts = append(opts, godaddy.WithTestEnvironment())
	}

	client := godaddy.NewClient(apiKey, apiSecret, opts...)
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *GoDaddyProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewDomainResource,
		NewDNSRecordResource,
	}
}

func (p *GoDaddyProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewDomainDataSource,
		NewDNSRecordsDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &GoDaddyProvider{
			version: version,
		}
	}
}
