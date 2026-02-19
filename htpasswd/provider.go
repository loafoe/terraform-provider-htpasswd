package htpasswd

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ provider.Provider = &HtpasswdProvider{}
var _ provider.ProviderWithEphemeralResources = &HtpasswdProvider{}

type HtpasswdProvider struct {
	version string
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &HtpasswdProvider{
			version: version,
		}
	}
}

func (p *HtpasswdProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "htpasswd"
	resp.Version = p.version
}

func (p *HtpasswdProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{}
}

func (p *HtpasswdProvider) Configure(_ context.Context, _ provider.ConfigureRequest, _ *provider.ConfigureResponse) {
}

func (p *HtpasswdProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewPasswordResource,
	}
}

func (p *HtpasswdProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

func (p *HtpasswdProvider) EphemeralResources(_ context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{
		NewPasswordEphemeral,
	}
}
