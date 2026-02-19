package htpasswd

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourcePassword_Complete(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourcePasswordConfig("1", "secret123", "saltySal"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("htpasswd_password.test_1", "password", "secret123"),
					resource.TestCheckResourceAttr("htpasswd_password.test_1", "salt", "saltySal"),
					resource.TestCheckResourceAttrSet("htpasswd_password.test_1", "sha256"),
					resource.TestCheckResourceAttrSet("htpasswd_password.test_1", "sha512"),
					resource.TestCheckResourceAttrSet("htpasswd_password.test_1", "apr1"),
					resource.TestCheckResourceAttrSet("htpasswd_password.test_1", "bcrypt"),
					// Check apr1 format: should start with $apr1$saltySal$
					resource.TestMatchResourceAttr("htpasswd_password.test_1", "apr1",
						regexp.MustCompile(`^\$apr1\$saltySal\$.+`)),
					// Check sha512 format: should start with $6$saltySal$
					resource.TestMatchResourceAttr("htpasswd_password.test_1", "sha512",
						regexp.MustCompile(`^\$6\$saltySal\$.+`)),
					// Check bcrypt format: should start with $2a$ or $2b$
					resource.TestMatchResourceAttr("htpasswd_password.test_1", "bcrypt",
						regexp.MustCompile(`^\$2[ab]\$\d+\$.+`)),
				),
			},
			{
				Config: testAccResourcePasswordConfig("2", "1234567890abcdefghijklmnopqrstuvwxyz", "12341234"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("htpasswd_password.test_2", "password", "1234567890abcdefghijklmnopqrstuvwxyz"),
					resource.TestCheckResourceAttr("htpasswd_password.test_2", "salt", "12341234"),
					resource.TestCheckResourceAttrSet("htpasswd_password.test_2", "sha256"),
					resource.TestCheckResourceAttrSet("htpasswd_password.test_2", "sha512"),
					resource.TestCheckResourceAttrSet("htpasswd_password.test_2", "apr1"),
					resource.TestCheckResourceAttrSet("htpasswd_password.test_2", "bcrypt"),
					// Check apr1 format: should start with $apr1$12341234$
					resource.TestMatchResourceAttr("htpasswd_password.test_2", "apr1",
						regexp.MustCompile(`^\$apr1\$12341234\$.+`)),
					// Check sha512 format: should start with $6$12341234$
					resource.TestMatchResourceAttr("htpasswd_password.test_2", "sha512",
						regexp.MustCompile(`^\$6\$12341234\$.+`)),
					// Check bcrypt format: should start with $2a$ or $2b$
					resource.TestMatchResourceAttr("htpasswd_password.test_2", "bcrypt",
						regexp.MustCompile(`^\$2[ab]\$\d+\$.+`)),
				),
			},
		},
	})
}

func TestAccResourcePassword_SHA512Regression(t *testing.T) {
	// CRITICAL REGRESSION TEST:
	// This test prevents regression of the SHA-512 hash generation bug
	// reported in version 1.4.0 where the following configuration:
	//   password = "1234567890abcdefghijklmnopqrstuvwxyz"
	//   salt     = "12341234"
	// produced a different SHA-512 hash than OpenSSL and version 1.3.0.
	//
	// Expected output (matches OpenSSL):
	// $6$12341234$b4koNtwY05CUmMhYkmcf9mU6K4QkuHVuVDcQWPpZoLf0dFXUggoBUV1O3MFBnAfApbrDrETCEhDdqyzSBHGvm1
	//
	// DO NOT change this expected value unless you're certain the implementation
	// still matches OpenSSL `passwd -6 -salt 12341234 1234567890abcdefghijklmnopqrstuvwxyz`
	expectedSHA512 := "$6$12341234$b4koNtwY05CUmMhYkmcf9mU6K4QkuHVuVDcQWPpZoLf0dFXUggoBUV1O3MFBnAfApbrDrETCEhDdqyzSBHGvm1"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourcePasswordConfig("regression", "1234567890abcdefghijklmnopqrstuvwxyz", "12341234"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("htpasswd_password.test_regression", "password", "1234567890abcdefghijklmnopqrstuvwxyz"),
					resource.TestCheckResourceAttr("htpasswd_password.test_regression", "salt", "12341234"),
					// Critical: Verify exact SHA-512 output to prevent regression
					resource.TestCheckResourceAttr("htpasswd_password.test_regression", "sha512", expectedSHA512),
					// Also verify it has the correct format
					resource.TestMatchResourceAttr("htpasswd_password.test_regression", "sha512",
						regexp.MustCompile(`^\$6\$12341234\$.+`)),
				),
			},
		},
	})
}

func testAccResourcePasswordConfig(postfix, password, salt string) string {
	return fmt.Sprintf(`
resource "htpasswd_password" "test_%s" {
	password = "%s"
	salt     = "%s"
}
`, postfix, password, salt)
}
