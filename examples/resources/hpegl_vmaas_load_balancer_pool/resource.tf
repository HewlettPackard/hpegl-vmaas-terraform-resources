# (C) Copyright 2022 Hewlett Packard Enterprise Development LP

resource "hpegl_vmaas_load_balancer_pool" "test_lb_pool" {
  lb_id = data.hpegl_vmaas_load_balancer.tf_lb.id 
  name  =  "TEST-POOL"       
  description  = "creating load balancer pool"
  min_active_members     = 1
  algorithm = "WEIGHTED_ROUND_ROBIN"
  config {
    snat_translation_type = "LBSnatAutoMap"
    active_monitor_paths = data.hpegl_vmaas_active_monitor.tf_lb_active.id
    passive_monitor_path = data.hpegl_vmaas_passive_monitor.tf_lb_passive.id
    tcp_multiplexing = false
    tcp_multiplexing_number = 6 
    snat_ip_address = ""
  }
}