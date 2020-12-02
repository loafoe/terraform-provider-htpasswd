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
		ReadContext: dataSourcePasswordRead,
		Schema: map[string]*schema.Schema{
			"password": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
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
		apr1Hash, err := apr1.Hash(password, "")
		if err != nil {
			return diag.FromErr(err)
		}
		d.Set("apr1", apr1Hash)
	}
	d.SetId(fmt.Sprintf("PW%x", password))

	return diags
}
