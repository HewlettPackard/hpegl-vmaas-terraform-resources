# (C) Copyright 2022 Hewlett Packard Enterprise Development LP

resource "hpegl_vmaas_load_balancer" "tf_load_balancer" {
  name  = "lab-1"       
  description  = "creating load balancer for test"
  enabled      =    true        
  visibility   = "public"
  resource_permission {
    all = true
  }
  config {
    admin_state = true
    tier1  = "/infra/tier-1s/26cdb82e-0057-4461-ad4d-cddd61d77b1f"
  }
}
