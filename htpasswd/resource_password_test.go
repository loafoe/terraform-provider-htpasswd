package htpasswd

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourcePassword_Complete(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders, // Use the existing testAccProviders
		Steps: []resource.TestStep{
			{
				Config: testAccResourcePasswordConfig("secret123", "saltySalt"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("htpasswd_password.test", "password", "secret123"),
					resource.TestCheckResourceAttr("htpasswd_password.test", "salt", "saltySalt"),
					resource.TestCheckResourceAttrSet("htpasswd_password.test", "sha256"),
					resource.TestCheckResourceAttrSet("htpasswd_password.test", "sha512"),
					resource.TestCheckResourceAttrSet("htpasswd_password.test", "apr1"),
					resource.TestCheckResourceAttrSet("htpasswd_password.test", "bcrypt"),
				),
			},
		},
	})
}

func testAccResourcePasswordConfig(password, salt string) string {
	return fmt.Sprintf(`
resource "htpasswd_password" "test" {
	password = "%s"
	salt     = "%s"
}
`, password, salt)
}
