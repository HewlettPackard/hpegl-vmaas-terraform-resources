# (C) Copyright 2022 Hewlett Packard Enterprise Development LP

resource "hpegl_vmaas_load_balancer" "tf_load_balancer" {
  name  = "tf_LB"  
  description  = "creating load balancer for test"
  enabled      =    true      
  resource_permission {
    all = true
  }
  config {
    log_level = "INFO"
    size = "SMALL"
    admin_state = true
    tier1 = data.hpegl_vmaas_lb_tier1.lb.id
  }
}