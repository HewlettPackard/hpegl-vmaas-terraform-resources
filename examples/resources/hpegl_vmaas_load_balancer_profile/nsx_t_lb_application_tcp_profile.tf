# (C) Copyright 2022 Hewlett Packard Enterprise Development LP

# APPLICATION Profile for TCP service
resource "hpegl_vmaas_load_balancer_profile" "tf_lb_profile" {
  lb_id = data.hpegl_vmaas_load_balancer.lb.id  
  name  =  "APPLICATION-TCP-Profile"       
  description  = "creating LB Profile"
  service_type     = "LBFastTcpProfile"
  config {
    profile_type = "application-profile"
    idle_timeout = 1800
    connection_close_timeout = 8
    ha_flow_mirroring = false
    tags {
        tag = "tag1"
        scope = "scope1"
    }
  }
}