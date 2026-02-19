# htpasswd provider

The htpasswd provider has a convenience resource which helps generate output
that is related to the Apache htpasswd password file format. As an example it
can generate `apr1` hashed passwords for use by nginx without needing to shell
out to local tools or binaries. This also makes it Terraform Cloud friendly.

You can also use to create a stable `bcrypt` hash of the password across
Terraform runs. More recent versions also support `SHA-512`

## Configuring the provider

```hcl
provider "htpasswd" {
}
```

## Example usage

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
