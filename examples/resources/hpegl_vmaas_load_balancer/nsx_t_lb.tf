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
    tier1 = "/infra/tier-1s/26cdb82e-0057-4461-ad4d-cddd61d77b1f"
  }
}