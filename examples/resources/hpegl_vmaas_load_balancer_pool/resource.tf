# (C) Copyright 2022 Hewlett Packard Enterprise Development LP

resource "hpegl_vmaas_load_balancer_pool" "lb_pool" {
  lb_id = data.hpegl_vmaas_load_balancer.tf_lb.id 
  name  =  "tf_POOL"       
  description  = "pool created using tf"
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