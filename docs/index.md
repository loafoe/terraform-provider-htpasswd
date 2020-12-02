# htpasswd provider

The htpassswd provider has convenience data sources which help generate output that is related to the Apache htpasswd 
password file format. As an example it can generate `apr1` hashed passwords for use by nginx without needing to shell
out to local tools or binaries. This also makes it Terraform Cloud friendly.

## Configuring the provider

```hcl
provider "htpassswd" {
}
```

## Argument Reference

No arguments are needed
