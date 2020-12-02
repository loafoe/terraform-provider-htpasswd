# htpasswd_password

Generates hashes of provided password string

## Example Usage

```hcl
data "htpasswd_password" "nginx_data" {
  password = "SuperSecret!"
}
```

```hcl
output "apr1_hash" {
   value = data.htpasswd_password.nginx_data.apr1
}
```

## Attributes Reference

The following attributes are exported:

* `password` - (Required) The password string
* `salt` - (Optional) Salt for apr1 hash generation. Must 8-charachter string or empty. Default: `""`
* `apr1` - (Computed) The apr1 hash of the password
