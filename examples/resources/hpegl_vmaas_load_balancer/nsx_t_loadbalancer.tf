# (C) Copyright 2022 Hewlett Packard Enterprise Development LP

resource hpegl_vmaas_load_balancer tf_load_balancer {
  name  = "lab-nsxt"        
  description  = ""creating load balancer for test"
  type     = "nsx-t"
  networkServerId = data.hpegl_vmaas_load_balancer_networkServerId.tf_networkServerId.serverID  
  enabled      =    true        
  visibility   = "private"
  config {
     adminState = true
      size = "SMALL"
      loglevel = "DEBUG"
      tier1  = "/infra/tier-1s/26cdb82e-0057-4461-ad4d-cddd61d77b1f"
  }
  resource_permissions {
    all  = true
  }
}
