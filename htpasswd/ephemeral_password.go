package htpasswd

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/johnaoss/htpasswd/apr1"
	"golang.org/x/crypto/bcrypt"
)

var _ ephemeral.EphemeralResource = &PasswordEphemeral{}

type PasswordEphemeral struct{}

type PasswordEphemeralModel struct {
	Password types.String `tfsdk:"password"`
	Salt     types.String `tfsdk:"salt"`
	Apr1     types.String `tfsdk:"apr1"`
	Bcrypt   types.String `tfsdk:"bcrypt"`
	Sha256   types.String `tfsdk:"sha256"`
	Sha512   types.String `tfsdk:"sha512"`
}

func NewPasswordEphemeral() ephemeral.EphemeralResource {
	return &PasswordEphemeral{}
}

func (r *PasswordEphemeral) Metadata(_ context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_password"
}

func (r *PasswordEphemeral) Schema(_ context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Generates htpasswd compatible password hashes without storing the password in state",
		Attributes: map[string]schema.Attribute{
			"password": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "The password to hash",
			},
			"salt": schema.StringAttribute{
				Optional:    true,
				Description: "Salt for apr1 and sha512 hashes. Must be exactly 8 characters from the crypt base64 alphabet.",
				Validators: []validator.String{
					&saltValidator{},
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

func (r *PasswordEphemeral) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var data PasswordEphemeralModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
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

	data.Bcrypt = types.StringValue(string(bcryptHash))
	data.Apr1 = types.StringValue(apr1Hash)
	data.Sha512 = types.StringValue(sha512hash)

	resp.Diagnostics.Append(resp.Result.Set(ctx, &data)...)
}
