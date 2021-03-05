# htpasswd_password

Generates hashes of provided password string

## Example Usage

```hcl
resource "random_password" "password" {
  length = 30
}

data "htpasswd_password" "hash" {
  password = random_password.password.result
  salt     = substr(sha512(random_password.password.result), 0, 8)
}

output "apr1_hash" {
  value = data.htpasswd_password.hash.apr1
}
```

## Attributes Reference

The following attributes are exported:

* `password` - (Required) The password string
* `salt` - (Optional) Salt for apr1 hash generation. Must 8-charachter string or empty. Default: `""`
* `apr1` - (Computed) The apr1 hash of the password
