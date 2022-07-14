# (C) Copyright 2022 Hewlett Packard Enterprise Development LP

resource "hpegl_vmaas_load_balancer" "tf_lb" {
  name  = "lab-2"       
  description  = "CREATE load balancer for test"
  enabled      =  true 

  group_access {
    all = true
  }

  config {
    admin_state = true
    size = "SMALL"
    log_level = "INFO"
    tier1_gateways  = data.hpegl_vmaas_tier1_router.tier1_router.provider_id
  }
}

