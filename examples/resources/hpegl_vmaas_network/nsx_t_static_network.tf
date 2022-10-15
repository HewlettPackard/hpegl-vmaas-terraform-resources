# (C) Copyright 2021 Hewlett Packard Enterprise Development LP

resource "hpegl_vmaas_network" "test_net" {
  name         = "tf_nsx_t_static_network"
  description  = "Static Network create using tf"
  scope_id     = data.hpegl_vmaas_transport_zone.tf_zone.provider_id
  cidr         = "168.72.10.1/18"
  gateway      = "168.72.10.1"
  primary_dns  = "8.8.8.8"
  scan_network = false
  pool_id = 7
  active       = true
  resource_permissions {
    all = true
  }
  group {
    id = "shared"
  }
  static_network {
    config {
      vlan_ids = "0,3-5"
      connected_gateway = data.hpegl_vmaas_router.tier1_router.provider_id
    }
  }
}