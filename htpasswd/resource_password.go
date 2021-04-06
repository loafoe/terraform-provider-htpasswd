package htpasswd

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-cty/cty"
	"golang.org/x/crypto/bcrypt"

	"github.com/johnaoss/htpasswd/apr1"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePassword() *schema.Resource {
	return &schema.Resource{
		CreateContext: datasourcePasswordCreate,
		ReadContext:   repopulateHashes,
		Delete:        schema.RemoveFromState,
		Schema: map[string]*schema.Schema{
			"password": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"salt": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				ValidateDiagFunc: validateSalt,
				Default:          "",
			},
			"apr1": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bcrypt": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func datasourcePasswordCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	password := d.Get("password").(string)

	bcryptHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return diag.FromErr(err)
	}
	_ = d.Set("bcrypt", string(bcryptHash))
	d.SetId(fmt.Sprintf("PW%x", string(bcryptHash)))
	return repopulateHashes(ctx, d, m)
}

func repopulateHashes(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	id := d.Id()
	var bcryptString string
	_, _ = fmt.Sscanf(id, "PW%x", &bcryptString)

	password := d.Get("password").(string)

	salt := d.Get("salt").(string)
	apr1Hash, err := apr1.Hash(password, salt)
	if err != nil {
		return diag.FromErr(err)
	}
	_ = d.Set("apr1", apr1Hash)
	_ = d.Set("bcrypt", bcryptString)
	return diags
}

func validateSalt(i interface{}, _ cty.Path) diag.Diagnostics {
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
