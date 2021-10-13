# (C) Copyright 2021 Hewlett Packard Enterprise Development LP

resource "hpegl_vmaas_network" "test_net" {
  name        = "tf_net_custom"
  group_id    = data.hpegl_vmaas_group.tf_group.id
  scope_id    = "/infra/sites/default/enforcement-points/default/transport-zones/88cd4dc8-0445-4b8e-b260-0f4cd361f4e1"
  dhcp_server = true
  description = "Network created using tf"
  cidr        = "168.72.10.0/18"
  gateway     = "168.72.10.9"
  netmask     = "255.255.255.255"
  active      = true
  config {
    connected_gateway = data.hpegl_vmaas_router.tf_router.provider_id
  }
}
