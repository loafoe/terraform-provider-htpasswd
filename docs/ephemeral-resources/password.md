# htpasswd_password (Ephemeral Resource)

Generates htpasswd compatible password hashes without storing the password in
state.

Ephemeral resources are a special type of resource that do not persist any data
in Terraform state. This makes them ideal for generating sensitive password
hashes that should not be stored in state files.

## Example Usage

```hcl
ephemeral "htpasswd_password" "hash" {
  password = var.password
  salt     = "abcdefgh"
}

resource "some_resource" "example" {
  password_hash = ephemeral.htpasswd_password.hash.apr1
}
```

### Using with random_password

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

ephemeral "htpasswd_password" "hash" {
  password = random_password.password.result
  salt     = random_password.salt.result
}
```

## Argument reference

The following arguments are supported:

* `password` - (Required, Sensitive) The password string to hash.
* `salt` - (Optional) Salt for apr1 and sha512 hash generation.
  Must be exactly 8 characters or empty. Valid characters are the crypt-style
  base64 alphabet: `./0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz`

## Attribute reference

In addition to all arguments above, the following attributes are exported:

* `apr1` - (Computed) The APR1-MD5 hash of the password.
* `bcrypt` - (Computed) The bcrypt hash of the password.
* `sha256` - (Computed) The SHA-256 hash of the password (hex encoded).
* `sha512` - (Computed) The SHA-512 crypt hash of the password.

## When to use Ephemeral vs Resource

Use the **ephemeral resource** (`ephemeral "htpasswd_password"`) when:

* You want to avoid storing password hashes in Terraform state
* The hash is being passed directly to another resource that stores it
* Security policies require minimizing sensitive data in state

Use the **managed resource** (`resource "htpasswd_password"`) when:

* You need the hash value to persist across Terraform runs
* You're outputting the hash for use outside of Terraform
* You need Terraform to track changes to the hash over time
