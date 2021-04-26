# Copyright 2020 Hewlett Packard Enterprise Development LP

# Set-up for terraform >= v0.13
terraform {
  required_providers {
    hpegl = {
      # We are specifying a location that is specific to the service under development
      # In this example it is poc-caas (see "source" below).  The service-specific replacement
      # to poc-caas must be specified in "source" below and also in the Makefile as the
      # value of DUMMY_PROVIDER.
      source  = "terraform.example.com/poc-caas/hpegl"
      version = ">= 0.0.1"
    }
  }
}

provider hpegl {
  caas_api_url = "https://client.greenlake.hpe-gl-intg.com/api/caas/mcaas/v1"
}

resource hpegl_caas_cluster test {
  name         = "tf-test-clus-22"
  blueprint_id = "2b2bb40c-813c-4762-9b49-faaebe1a4e61"
  appliance_id = "3ad9c737-5bb6-430c-9772-3a6f5a7e4015"
  space_id     = "6edb9418-dcda-4517-bf0c-a5d7de9cc60a"
}
