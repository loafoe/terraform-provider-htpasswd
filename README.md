# htpasswd Terraform provider

- Documentation: https://registry.terraform.io/providers/loafoe/htpasswd/latest/docs

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

## Overview

This is a terraform provider to generate output related to Apache htpasswd file

# Using the provider

**Terraform 0.13**: To install this provider, copy and paste this code into your Terraform configuration. Then, run terraform init.

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

-	[Terraform](https://www.terraform.io/downloads.html) 0.13.x
-	[Go](https://golang.org/doc/install) 1.15 or newer (to build the provider plugin)

## Issues

- If you have an issue: report it on the [issue tracker](https://github.com/loafoe/terraform-provider-htpasswd/issues)

## LICENSE

License is MIT
