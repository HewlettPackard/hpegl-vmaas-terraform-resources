# (C) Copyright 2022 Hewlett Packard Enterprise Development LP

resource "hpegl_vmaas_load_balancer" "tf_load_balancer" {
  name  = "TEST-LB" 
  type = "nsx-t" 
  network_server_id = 1     
  description  = "creating load balancer for test"  
  enabled      =    true        
  visibility   = "private"
resource_permission {
all =true
}
  config {
     admin_state = true
  }
}
