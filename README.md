# htpasswd Terraform provider

[Documentation](https://registry.terraform.io/providers/loafoe/htpasswd/latest/docs)

## ⚠️ Breaking Change in 1.6.0+

Version 1.6.0 introduces a change in the salt generation algorithm. This means
that **generated password hashes may differ** from those created by earlier
versions, even with the same input. If you upgrade and run `terraform plan`,
you may see resources marked for replacement. Review carefully before applying.

## Overview

This is a Terraform provider to generate htpasswd-compatible password hashes
(`apr1`, `bcrypt`, `sha256`, `sha512`) for use with Apache, nginx, and other
web servers. It works without shelling out to local tools, making it Terraform
Cloud friendly.

## Features

* **Managed resource** (`htpasswd_password`) - Password hashes stored in state
* **Ephemeral resource** (`htpasswd_password`) - Password hashes generated
  without storing in state (requires Terraform 1.10+ or OpenTofu 1.8+)

## Using the provider

To install this provider, copy and paste this code into your Terraform
configuration, then run `terraform init`.

```terraform
terraform {
  required_providers {
    htpasswd = {
      source = "loafoe/htpasswd"
    }
  }
}
```

## Requirements

| Feature | Terraform | OpenTofu |
|---------|-----------|----------|
| Managed resources | 1.0+ | 1.0+ |
| Ephemeral resources | 1.10+ | 1.8+ |

## Development requirements

* [Terraform](https://www.terraform.io/downloads.html) 1.10 or newer
* [Go](https://golang.org/doc/install) 1.25 or newer (to build the provider
  plugin)

## Issues

If you have an issue: report it on the
[issue tracker](https://github.com/loafoe/terraform-provider-htpasswd/issues)

## LICENSE

License is MIT
