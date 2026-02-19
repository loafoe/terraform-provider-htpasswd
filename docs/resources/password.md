# htpasswd_password

Generate hashes of provided password string

## Example Usage

```hcl
resource "random_password" "password" {
  length           = 30
  special          = true
  override_special = "!@#%&*()-_=+[]{}<>:?"
}

resource "random_password" "salt" {
  length           = 8
  special          = true
  override_special = "./"
}

resource "htpasswd_password" "hash" {
  password = random_password.password.result
  salt     = random_password.salt.result
}

output "password" {
  value = random_password.password.result
  sensitive = true
}

output "apr1_password" {
  value = random_password.password.result
}

output "apr1_hash" {
  value = htpasswd_password.hash.apr1
}

output "bcrypt_hash" {
  value = htpasswd_password.hash.bcrypt
}

output "sha256_hash" {
  value = htpasswd_password.hash.sha256
}

output "sha512_hash" {
  value = htpasswd_password.hash.sha512
}
```

## Argument reference

The following arguments are supported:

* `password` - (Required) The password string
* `salt` - (Optional) Salt for apr1 and sha512 hash generation.
  Must be exactly 8 characters or empty. Valid characters are the crypt-style
  base64 alphabet: `./0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz`
  Default: `""`

## Attribute reference

In addition to all arguments above, the following attributes are exported:

* `apr1` - (Computed) The apr1 hash of the password
* `bcrypt` - (Computed) the bcrypt hash of the password
* `sha256` - (Computed) the SHA-256 hash of the password
* `sha512` - (Computed) the SHA-512 hash of the password
