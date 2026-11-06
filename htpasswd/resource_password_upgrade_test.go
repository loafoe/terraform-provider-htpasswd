package htpasswd

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccResourcePassword_LegacyHashUpgrade(t *testing.T) {
	for _, salt := range []string{"abcdefgh", "ab!@x?yz"} {
		t.Run(salt, func(t *testing.T) {
			const address = "htpasswd_password.test_upgrade"
			config := testAccResourcePasswordConfig("upgrade", "synthetic-upgrade-check", salt)
			original := make(map[string]string)
			capture := func(state *terraform.State) error {
				instance := state.RootModule().Resources[address].Primary
				for _, name := range []string{"id", "apr1", "bcrypt", "sha256", "sha512"} {
					original[name] = instance.Attributes[name]
				}
				return nil
			}
			preserved := func(state *terraform.State) error {
				instance := state.RootModule().Resources[address].Primary
				for name, value := range original {
					if instance.Attributes[name] != value {
						return fmt.Errorf("%s changed during legacy_hash update", name)
					}
				}
				return nil
			}

			steps := []resource.TestStep{
				{
					ExternalProviders: map[string]resource.ExternalProvider{
						"htpasswd": {Source: "loafoe/htpasswd", VersionConstraint: "= 1.5.0"},
					},
					Config: config,
					Check:  capture,
				},
			}
			for _, legacy := range []string{"true", "false", "null"} {
				steps = append(steps, resource.TestStep{
					ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
					ConfigPlanChecks: resource.ConfigPlanChecks{
						PreApply: []plancheck.PlanCheck{
							plancheck.ExpectKnownValue(address, tfjsonpath.New("apr1"), knownvalue.NotNull()),
							plancheck.ExpectKnownValue(address, tfjsonpath.New("bcrypt"), knownvalue.NotNull()),
							plancheck.ExpectKnownValue(address, tfjsonpath.New("sha256"), knownvalue.NotNull()),
							plancheck.ExpectKnownValue(address, tfjsonpath.New("sha512"), knownvalue.NotNull()),
						},
					},
					Config: fmt.Sprintf(`
resource "htpasswd_password" "test_upgrade" {
  password    = "synthetic-upgrade-check"
  salt        = %q
  legacy_hash = %s
}
`, salt, legacy),
					Check: resource.ComposeTestCheckFunc(
						preserved,
						resource.TestCheckResourceAttr(address, "sha1", "{SHA}qcvT9l6SkxoKbnTSePN3opW1SI8="),
					),
				})
			}
			resource.Test(t, resource.TestCase{Steps: steps})
		})
	}
}
