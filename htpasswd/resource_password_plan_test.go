package htpasswd

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccResourcePassword_ReplacementPlans(t *testing.T) {
	for _, tc := range []struct {
		name     string
		password string
		salt     string
		extra    string
	}{
		{name: "password", password: `"new-password"`, salt: `"abcdefgh"`},
		{name: "salt", password: `"secret123"`, salt: `"12345678"`},
		{
			name:     "unknown-password",
			password: "terraform_data.value.output",
			salt:     `"abcdefgh"`,
			extra:    `resource "terraform_data" "value" { input = "new-password" }`,
		},
		{
			name:     "unknown-salt",
			password: `"secret123"`,
			salt:     "terraform_data.value.output",
			extra:    `resource "terraform_data" "value" { input = "12345678" }`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			const address = "htpasswd_password.test_replace"
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{Config: testAccResourcePasswordConfig("replace", "secret123", "abcdefgh")},
					{
						Config: fmt.Sprintf(`
%s
resource "htpasswd_password" "test_replace" {
  password = %s
  salt     = %s
}
`, tc.extra, tc.password, tc.salt),
						ConfigPlanChecks: resource.ConfigPlanChecks{
							PreApply: []plancheck.PlanCheck{
								plancheck.ExpectResourceAction(address, plancheck.ResourceActionDestroyBeforeCreate),
								plancheck.ExpectUnknownValue(address, tfjsonpath.New("apr1")),
								plancheck.ExpectUnknownValue(address, tfjsonpath.New("bcrypt")),
								plancheck.ExpectUnknownValue(address, tfjsonpath.New("sha256")),
								plancheck.ExpectUnknownValue(address, tfjsonpath.New("sha512")),
							},
						},
					},
				},
			})
		})
	}
}
