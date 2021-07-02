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
    iam_service_url = "https://iam.us1.greenlake-hpe.com"
    tenant_id = "<GLC-Tenant-ID>"
    user_id = "<SCM-Client-ID>"
    user_secret = "<SCM-Client-Secret>"
}
