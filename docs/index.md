# htpasswd provider

The htpassswd provider has convenience resource which helps generate output that is related to the Apache htpasswd 
password file format. As an example it can generate `apr1` hashed passwords for use by nginx without needing to shell
out to local tools or binaries. This also makes it Terraform Cloud friendly.

You can also use to create a stable bcrypt hash of the password across Terraform runs.

## Configuring the provider

```hcl
provider "htpasswd" {
}
```

## Argument Reference

No arguments are needed
