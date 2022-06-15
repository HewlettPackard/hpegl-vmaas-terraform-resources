# (C) Copyright 2022 Hewlett Packard Enterprise Development LP

resource hpegl_vmaas_load_balancer_pool test_lb_pool {
  name  =  "loadbalancer pool"       
  description  = "creating load balancer pool"
  min_active     = "data.hpegl_vmaas_load_balancer_pool_minActive.tf_minActive.min_active"
  vip_balance = "data.hpegl_vmaas_load_balancer_pool_vipBalance.tf_vipBalance.vip_balance"
  config {
    snat_translation_type = "data.hpegl_vmaas_load_balancer_pool_snatTranslationType.tf_snatTranslationType.snat_translation_type"
    passive_monitor_path = 136
    active_monitor_paths = 133
    tcp_multiplexing = false
    tcp_multiplexing_number = 6 
    snat_ip_address = ""
  }
}
