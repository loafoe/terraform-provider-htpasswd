package htpasswd

import (
	"context"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"

	"github.com/hashicorp/go-cty/cty"
	"golang.org/x/crypto/bcrypt"

	"github.com/johnaoss/htpasswd/apr1"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// sha512Crypt implements the SHA-512 crypt algorithm as specified in
// http://www.akkadia.org/drepper/SHA-crypt.txt
func sha512Crypt(password, salt string) string {
	const rounds = 5000
	const prefix = "$6$"

	// Custom base64 alphabet used by crypt
	const alphabet = "./0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	// Ensure salt is maximum 16 characters
	if len(salt) > 16 {
		salt = salt[:16]
	}

	// Step 1: Compute alternate sum
	h := sha512.New()
	h.Write([]byte(password))
	h.Write([]byte(salt))
	h.Write([]byte(password))
	altResult := h.Sum(nil)

	// Step 2: Compute main sum
	h.Reset()
	h.Write([]byte(password))
	h.Write([]byte(salt))

	// Add altResult for each character in password
	for i := len(password); i > 0; i -= 64 {
		if i > 64 {
			h.Write(altResult)
		} else {
			h.Write(altResult[:i])
		}
	}

	// Add password or altResult based on password length bits
	for i := len(password); i > 0; i >>= 1 {
		if (i & 1) != 0 {
			h.Write(altResult)
		} else {
			h.Write([]byte(password))
		}
	}

	result := h.Sum(nil)

	// Step 3: Compute P sequence
	h.Reset()
	for i := 0; i < len(password); i++ {
		h.Write([]byte(password))
	}
	pBytes := h.Sum(nil)

	// Create P sequence
	p := make([]byte, 0, len(password))
	for i := len(password); i > 0; i -= 64 {
		if i > 64 {
			p = append(p, pBytes...)
		} else {
			p = append(p, pBytes[:i]...)
		}
	}

	// Step 4: Compute S sequence
	h.Reset()
	for i := 0; i < 16+int(result[0]); i++ {
		h.Write([]byte(salt))
	}
	sBytes := h.Sum(nil)

	// Create S sequence
	s := make([]byte, 0, len(salt))
	for i := len(salt); i > 0; i -= 64 {
		if i > 64 {
			s = append(s, sBytes...)
		} else {
			s = append(s, sBytes[:i]...)
		}
	}

	// Step 5: Perform rounds iterations
	for round := 0; round < rounds; round++ {
		h.Reset()

		if (round & 1) != 0 {
			h.Write(p)
		} else {
			h.Write(result)
		}

		if round%3 != 0 {
			h.Write(s)
		}

		if round%7 != 0 {
			h.Write(p)
		}

		if (round & 1) != 0 {
			h.Write(result)
		} else {
			h.Write(p)
		}

		result = h.Sum(nil)
	}

	// Step 6: Encode result using the SHA-512 crypt base64 encoding
	// This follows the exact specification for SHA-512 crypt encoding
	encoded := ""

	// Specific byte reordering for SHA-512 crypt as per specification
	indices := []int{
		0, 21, 42, 22, 43, 1, 44, 2, 23, 3, 24, 45, 25, 46, 4, 47, 5, 26,
		6, 27, 48, 28, 49, 7, 50, 8, 29, 9, 30, 51, 31, 52, 10, 53, 11, 32,
		12, 33, 54, 34, 55, 13, 56, 14, 35, 15, 36, 57, 37, 58, 16, 59, 17, 38,
		18, 39, 60, 40, 61, 19, 62, 20, 41,
	}

	// Process in groups of 3 bytes
	for i := 0; i < len(indices); i += 3 {
		var val int
		if i+2 < len(indices) {
			// Standard 3-byte group
			val = (int(result[indices[i]]) << 16) |
				(int(result[indices[i+1]]) << 8) |
				int(result[indices[i+2]])
			// Encode as 4 characters
			for j := 0; j < 4; j++ {
				encoded += string(alphabet[val&0x3f])
				val >>= 6
			}
		} else {
			// Handle remaining bytes
			for j := i; j < len(indices) && j < i+3; j++ {
				val = (val << 8) | int(result[indices[j]])
			}
			// Encode based on how many bytes we have
			chars := ((len(indices)-i)*8 + 5) / 6
			for j := 0; j < chars; j++ {
				encoded += string(alphabet[val&0x3f])
				val >>= 6
			}
		}
	}

	// Handle the final byte (index 63)
	val := int(result[63])
	encoded += string(alphabet[val&0x3f])
	val >>= 6
	encoded += string(alphabet[val&0x3f])

	return prefix + salt + "$" + encoded
}

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
	sha512hash := sha512Crypt(password, salt)
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
		if len(s) != 0 && len(s) != 8 {
			diags = append(diags, diag.Errorf("Salt must be 8 characters exactly")...)
		}
	} else {
		diags = append(diags, diag.Errorf("Provided salt is not a string")...)
	}
	return diags
}
