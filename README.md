![Acceptance workflow](https://github.com/HewlettPackard/hpegl-vmaas-terraform-resources/actions/workflows/acc.yml/badge.svg)

# vmaas-terraform-resources

- [vmaas-terraform-resources](#vmaas-terraform-resources)
    * [Introduction](#terraform-resources-for-hpegl-vmaas)
    * [Requirements](#Requirements)
    * [Usage](#Usage)
    * [Building the resources as provider](#building-the-resources-as-provider)

# Terraform resources for HPEGL VMaaS

Terraform VMaaS resources is a plugin for HPEGL terraform provider that allows the full lifecycle management of HPEGL
VMaaS resources. This provider is maintained by [HPEGL VMaaS resources team](mailTo:glcs.team-avion@hpe.com).

## Requirements

1. Terraform version >= v0.13 [install terraform](https://learn.hashicorp.com/tutorials/terraform/install-cli)
2. A Service Client to authenticate against GreenLake.
3. Terraform basics. [Terraform Introduction](https://www.terraform.io/intro/index.html)

## Usage

See the terraform provider for
hpegl [documentation](https://registry.terraform.io/providers/HewlettPackard/hpegl/latest/docs)
to get started using the provider.

## Building the resources as provider

Clone repository to: $GOPATH/src/github.com/HewlettPackard/hpegl-vmaas-terraform-resources

```bash
$ mkdir -p $GOPATH/src/github.com/HewlettPackard/
$ cd $GOPATH/src/github.com/HewlettPackard
$ git clone https://github.com/HewlettPackard/hpegl-vmaas-terraform-resources.git
```

Enter the provider directory and build resources as provider

```bash
$ cd $GOPATH/src/github.com/HewlettPackard/hpegl-vmaas-terraform-resources
$ make build 
```

Note: For debugging the provider please refer to the
[debugging guide](https://medium.com/@gandharva666/debugging-terraform-using-jetbrains-goland-f9a7e992cb1d)
