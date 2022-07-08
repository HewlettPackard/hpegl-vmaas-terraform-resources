# (C) Copyright 2022 Hewlett Packard Enterprise Development LP

resource "hpegl_vmaas_load_balancer" "tf_lb" {
  name  = "tf_lb"       
  description  = "CREATING load balancer for test"
  enabled      =    true        
  resource_permissions {
    all = true
  }
  config {
    admin_state = true
    size = "SMALL"
    log_level = "INFO"
    tier1  = data.hpegl_vmaas_tier1_router.tier1_router.provider_id
  }
}
