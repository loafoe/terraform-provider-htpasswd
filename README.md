# htpasswd Terraform provider

[Documentation](https://registry.terraform.io/providers/loafoe/htpasswd/latest/docs)

## Overview

This is a Terraform provider to generate output related to Apache htpasswd file

## Using the provider

**Terraform 1.0+**: To install this provider, copy and paste this code into
your Terraform configuration. Then, run terraform init.

```terraform
terraform {
  required_providers {
    htpasswd = {
      source = "loafoe/htpasswd"
    }
  }
}
```

## Development requirements

* [Terraform](https://www.terraform.io/downloads.html) 1.0 or newer
* [Go](https://golang.org/doc/install) 1.25 or newer (to build the provider
  plugin)

## Issues

If you have an issue: report it on the
[issue tracker](https://github.com/loafoe/terraform-provider-htpasswd/issues)

## LICENSE

License is MIT
