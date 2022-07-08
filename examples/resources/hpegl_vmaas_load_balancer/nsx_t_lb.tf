# (C) Copyright 2022 Hewlett Packard Enterprise Development LP

resource "hpegl_vmaas_load_balancer" "tf_load_balancer" {
  name  = "lab-1"       
  description  = "creating load balancer for test"
  enabled      =    true        
  visibility   = "public"
  resource_permissions {
    all = true
  }
  config {
    admin_state = true
    tier1  = data.hpegl_vmaas_router.tier1_router.provider_id
  }
}
