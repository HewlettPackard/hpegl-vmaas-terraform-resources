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

provider hpegl {
  vmaas_api_url = "https://client.greenlake.hpe-gl-intg.com/api/v1/vmaas/"
  iam_tokem = "iam-token"

}

resource hpegl_vmaas_vm test {
  name         = "Terrform-VM-1"
  zone_id      = "10"
  cloud_name   = "HPE GreenLake VMaaS Cloud"
  site_id      = "4"
  type         = "centos"
  instance_type_code = "centos"
  layout_id = "2"
  resourcepool_id = "1"
  agent_install = "yes"
  plan_id = "10"
  volume_size = "20"
  datastore_id = "1"
  network_id = "5"

}
