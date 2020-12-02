package htpasswd

import (
	"context"

	"github.com/hashicorp/go-cty/cty"

	"github.com/johnaoss/htpasswd/apr1"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePassword() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePasswordRead,
		Schema: map[string]*schema.Schema{
			"password": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"salt": &schema.Schema{
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validateSalt,
				Default:          "",
			},
			"apr1": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourcePasswordRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	password := d.Get("password").(string)

	if d.IsNewResource() || d.HasChange("password") {
		salt := d.Get("salt").(string)
		apr1Hash, err := apr1.Hash(password, salt)
		if err != nil {
			return diag.FromErr(err)
		}
		d.Set("apr1", apr1Hash)
	}
	return diags
}

func validateSalt(i interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics
	if s, ok := i.(string); ok {
		if !(len(s) == 0 || len(s) == 8) {
			diags = append(diags, diag.Errorf("must be 8 chars exactly")...)
		}
	} else {
		diags = append(diags, diag.Errorf("not a string")...)
	}
	return diags
}
