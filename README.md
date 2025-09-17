# terraform-provider-sys11dbaas

[[_TOC_]]

## Generator

This project was initially generated using [terraform-plugin-codegen-openapi](https://github.com/hashicorp/terraform-plugin-codegen-openapi)
and [terraform-plugin-codegen-framework](https://github.com/hashicorp/terraform-plugin-codegen-framework)  
However, i had to fix so much stuff by hand, i don't think those can ever be used again to add something to this provider.
But feel free to try in a branch. I left the provider-spec file in the repo.

## Development

Override the provider in your `~/.terraformrc` file:

```terraform
provider_installation {

  dev_overrides {
      "registry.terraform.io/syseleven/sys11dbaas" = "/home/me/go/bin"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}

```

Then run `go install` in this directory. You should now be able to use the provider.

```terraform
terraform {
  required_providers {
    sys11dbaas = {
      source = "registry.terraform.io/syseleven/sys11dbaas"
    }
  }
}

provider "sys11dbaas" {
 [...]
}
```

Don't forget to run `go install` again after code changes.

## Debug logging

You can enable debug logging by setting the environment variable `SYS11DBAAS_SDK_DEBUG=true` additionally to `TF_LOG=DEBUG`:

```shell
SYS11DBAAS_SDK_DEBUG=true TF_LOG_PROVIDER=DEBUG terraform apply
```

## Generate docs

```bash
make generate-docs
```
