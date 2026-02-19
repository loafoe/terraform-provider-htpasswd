package htpasswd

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAccEphemeralPassword_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("1.10.0"))),
		},
		Steps: []resource.TestStep{
			{
				Config:             testAccEphemeralPasswordBasicConfig("test1", "secret123", "saltySal"),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccEphemeralPassword_WithLocalVariable(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("1.10.0"))),
		},
		Steps: []resource.TestStep{
			{
				Config:             testAccEphemeralPasswordWithLocalConfig("test2", "1234567890abcdefghijklmnopqrstuvwxyz", "12341234"),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccEphemeralPasswordBasicConfig(name, password, salt string) string {
	return fmt.Sprintf(`
ephemeral "htpasswd_password" "%s" {
  password = "%s"
  salt     = "%s"
}
`, name, password, salt)
}

func testAccEphemeralPasswordWithLocalConfig(name, password, salt string) string {
	return fmt.Sprintf(`
ephemeral "htpasswd_password" "%s" {
  password = "%s"
  salt     = "%s"
}

locals {
  apr1_hash   = ephemeral.htpasswd_password.%s.apr1
  bcrypt_hash = ephemeral.htpasswd_password.%s.bcrypt
  sha256_hash = ephemeral.htpasswd_password.%s.sha256
  sha512_hash = ephemeral.htpasswd_password.%s.sha512
}
`, name, password, salt, name, name, name, name)
}
