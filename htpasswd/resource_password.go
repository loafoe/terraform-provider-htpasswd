package htpasswd

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/hashicorp/go-cty/cty"
	"github.com/tredoe/osutil/v2/userutil/crypt/sha512_crypt"
	"golang.org/x/crypto/bcrypt"

	"github.com/johnaoss/htpasswd/apr1"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePassword() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePasswordCreate,
		ReadContext:   repopulateHashes,
		DeleteContext: resourcePasswordDelete,
		Schema: map[string]*schema.Schema{
			"password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
				ForceNew:  true,
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
			"sha512": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sha256": { // Adding support for SHA-256
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourcePasswordDelete(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	d.SetId("")
	return diags
}

func resourcePasswordCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
	c := sha512_crypt.New()
	sha512hash, err := c.Generate([]byte(password), []byte("$6$"+salt))
	if err != nil {
		return diag.FromErr(err)
	}
	h := sha256.New()
	h.Write([]byte(salt + password))
	sha256Hash := hex.EncodeToString(h.Sum(nil))
	
	_ = d.Set("sha256", sha256Hash)
	_ = d.Set("sha512", sha512hash)
	_ = d.Set("apr1", apr1Hash)
	_ = d.Set("bcrypt", bcryptString)
	return diags
}

func validateSalt(i interface{}, _ cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics
	if s, ok := i.(string); ok {
		if !(len(s) == 0 || len(s) == 8) {
			diags = append(diags, diag.Errorf("Salt must be 8 characters exactly")...)
		}
	} else {
		diags = append(diags, diag.Errorf("Provided salt is not a string")...)
	}
	return diags
}
