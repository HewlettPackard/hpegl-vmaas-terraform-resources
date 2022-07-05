# (C) Copyright 2022 Hewlett Packard Enterprise Development LP

resource "hpegl_vmaas_lb" "tf_load_balancer" {
  name  = "tf_LB"  
  description  = "creating load balancer for test"
  enabled      =    true      
  resource_permission {
    all = true
  }
  config {
    admin_state = true
  }
}