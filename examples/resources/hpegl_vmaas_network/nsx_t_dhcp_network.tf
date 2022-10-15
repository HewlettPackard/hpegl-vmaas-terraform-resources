# (C) Copyright 2021 Hewlett Packard Enterprise Development LP

resource "hpegl_vmaas_network" "dhcp_net" {
  name         = "tf_nsx_t_dhcp_network"
  description  = "DHCP Network create using tf"
  display_name = "tf_nsx_t_dhcp_network"
  scope_id     = data.hpegl_vmaas_transport_zone.tf_zone.provider_id
  cidr         = "193.3.0.1/24" 
  primary_dns = "8.8.8.8"
  secondary_dns = "8.8.8.8"
  scan_network = false
  active       = true
  allow_ip_override = false
  appliance_url_proxy_bypass = true
  group {
    id = "shared"
  }
  resource_permissions {
    all = true
  }
  dhcp_network{
    dhcp_enabled = true
    config {
      dhcp_type = "dhcpLocal"
      dhcp_server = "/infra/dhcp-server-configs/3b2124e4-fad5-4df9-8644-5acb69b1efac"
      dhcp_lease_time = "86400"
      dhcp_range = "192.168.1.0/24"
      dhcp_server_address = "193.3.0.0/24"
      vlan_ids = "0,3-5"
      connected_gateway = data.hpegl_vmaas_router.tier1_router.provider_id
    }
  }
}
