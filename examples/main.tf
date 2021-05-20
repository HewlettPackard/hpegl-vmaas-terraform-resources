# (C) Copyright 2021 Hewlett Packard Enterprise Development LP

# Set-up for terraform >= v0.13
terraform {
  required_providers {
    hpegl = {
      source  = "terraform.example.com/vmaas/hpegl"
      version = ">= 0.0.1"
    }
  }
}

provider "hpegl" {
  vmaas {
    location   = "location"
    space_name = "space_name"
  }
  iam_token = "iam-token"
}

resource "hpegl_vmaas_vm" "test" {
  name          = "test"
  cloud_id      = 1
  group_id      = 1
  plan_id       = 1
  instance_type = "test"
  networks      = [1]
  volumes {
    size         = 5
    datastore_id = "test"

  }
  volumes {
    size         = 10
    datastore_id = "test2"

  }
  labels = ["test"]
  tags = {
    name = "value"
    data = "data"
  }
  config {
    vmware_resource_pool = "test"
  }

  copies = 1

}
