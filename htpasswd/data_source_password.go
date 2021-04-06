package htpasswd

import (
	"context"
	"fmt"

	"github.com/johnaoss/htpasswd/apr1"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePassword() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "Please switch to the resource 'htpasswd_password'.",
		ReadContext:        dataSourcePasswordRead,
		Schema: map[string]*schema.Schema{
			"password": {
				Type:     schema.TypeString,
				Required: true,
			},
			"salt": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validateSalt,
				Default:          "",
			},
			"apr1": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourcePasswordRead(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	password := d.Get("password").(string)

	if d.IsNewResource() || d.HasChange("password") {
		salt := d.Get("salt").(string)
		apr1Hash, err := apr1.Hash(password, salt)
		if err != nil {
			return diag.FromErr(err)
		}
		_ = d.Set("apr1", apr1Hash)
	}
	d.SetId(fmt.Sprintf("PW%x", password))
	return diags
}
