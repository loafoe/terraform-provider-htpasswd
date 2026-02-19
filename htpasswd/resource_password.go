package htpasswd

import (
	"context"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/johnaoss/htpasswd/apr1"
	"golang.org/x/crypto/bcrypt"
)

var _ resource.Resource = &PasswordResource{}

type PasswordResource struct{}

type PasswordModel struct {
	ID       types.String `tfsdk:"id"`
	Password types.String `tfsdk:"password"`
	Salt     types.String `tfsdk:"salt"`
	Apr1     types.String `tfsdk:"apr1"`
	Bcrypt   types.String `tfsdk:"bcrypt"`
	Sha256   types.String `tfsdk:"sha256"`
	Sha512   types.String `tfsdk:"sha512"`
}

func NewPasswordResource() resource.Resource {
	return &PasswordResource{}
}

func (r *PasswordResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_password"
}

func (r *PasswordResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Generates htpasswd compatible password hashes",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Resource identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"password": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "The password to hash",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"salt": schema.StringAttribute{
				Optional:    true,
				Description: "Salt for apr1 and sha512 hashes. Must be exactly 8 characters from the crypt base64 alphabet.",
				Validators: []validator.String{
					&saltValidator{},
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"apr1": schema.StringAttribute{
				Computed:    true,
				Description: "APR1-MD5 hash of the password",
			},
			"bcrypt": schema.StringAttribute{
				Computed:    true,
				Description: "Bcrypt hash of the password",
			},
			"sha256": schema.StringAttribute{
				Computed:    true,
				Description: "SHA-256 hash of the password (hex encoded)",
			},
			"sha512": schema.StringAttribute{
				Computed:    true,
				Description: "SHA-512 crypt hash of the password",
			},
		},
	}
}

func (r *PasswordResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data PasswordModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	password := data.Password.ValueString()
	salt := data.Salt.ValueString()

	bcryptHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		resp.Diagnostics.AddError("Bcrypt Error", fmt.Sprintf("Failed to generate bcrypt hash: %s", err))
		return
	}

	apr1Hash, err := apr1.Hash(password, salt)
	if err != nil {
		resp.Diagnostics.AddError("APR1 Error", fmt.Sprintf("Failed to generate APR1 hash: %s", err))
		return
	}

	sha512hash := sha512Crypt(password, salt)

	h := sha256.New()
	h.Write([]byte(salt + password))
	sha256Hash := hex.EncodeToString(h.Sum(nil))

	data.ID = types.StringValue(fmt.Sprintf("PW%x", string(bcryptHash)))
	data.Bcrypt = types.StringValue(string(bcryptHash))
	data.Apr1 = types.StringValue(apr1Hash)
	data.Sha512 = types.StringValue(sha512hash)
	data.Sha256 = types.StringValue(sha256Hash)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PasswordResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data PasswordModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	password := data.Password.ValueString()
	salt := data.Salt.ValueString()

	id := data.ID.ValueString()
	var bcryptString string
	_, _ = fmt.Sscanf(id, "PW%x", &bcryptString)

	apr1Hash, err := apr1.Hash(password, salt)
	if err != nil {
		resp.Diagnostics.AddError("APR1 Error", fmt.Sprintf("Failed to generate APR1 hash: %s", err))
		return
	}

	sha512hash := sha512Crypt(password, salt)

	h := sha256.New()
	h.Write([]byte(salt + password))
	sha256Hash := hex.EncodeToString(h.Sum(nil))

	data.Bcrypt = types.StringValue(bcryptString)
	data.Apr1 = types.StringValue(apr1Hash)
	data.Sha512 = types.StringValue(sha512hash)
	data.Sha256 = types.StringValue(sha256Hash)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PasswordResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
}

func (r *PasswordResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}

// validSaltChars is the crypt-style base64 alphabet used for APR1/MD5-crypt salts
const validSaltChars = "./0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

type saltValidator struct{}

var _ validator.String = &saltValidator{}

func (v *saltValidator) Description(_ context.Context) string {
	return "salt must be exactly 8 characters from the crypt base64 alphabet"
}

func (v *saltValidator) MarkdownDescription(_ context.Context) string {
	return "salt must be exactly 8 characters from the crypt base64 alphabet"
}

func (v *saltValidator) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	s := req.ConfigValue.ValueString()
	if len(s) == 0 {
		return
	}

	if len(s) != 8 {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Salt Length",
			fmt.Sprintf("Salt must be exactly 8 characters, got %d", len(s)),
		)
	}

	for _, c := range s {
		if !strings.ContainsRune(validSaltChars, c) {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid Salt Character",
				fmt.Sprintf("Salt contains invalid character '%c'; valid characters are: %s", c, validSaltChars),
			)
			break
		}
	}
}

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

