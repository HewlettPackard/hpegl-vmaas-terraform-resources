# (C) Copyright 2022 Hewlett Packard Enterprise Development LP

resource "hpegl_vmaas_load_balancer_pool" "test_lb_pool" {
  lb_id = data.hpegl_vmaas_lb.lb.id 
  name  =  "TEST-POOL"       
  description  = "creating load balancer pool"
  min_active     = 1
  vip_balance = "ROUND_ROBIN"
  config {
    snat_translation_type = "LBSnatAutoMap"
    passive_monitor_path = 136
    active_monitor_paths = 133
    tcp_multiplexing = false
    tcp_multiplexing_number = 6 
    snat_ip_address = ""
    member_group {
      name = "Pushpa-group"
      path = ""
      ip_revision_filter = "IPV4"
      port = 80
    }
  }
}