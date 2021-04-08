package htpasswd

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{},
		ResourcesMap: map[string]*schema.Resource{
			"htpasswd_password": resourcePassword(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"htpasswd_password": dataSourcePassword(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(_ context.Context, _ *schema.ResourceData) (interface{}, diag.Diagnostics) {
	ron := "swanson"
	var diags diag.Diagnostics
	return ron, diags
}
