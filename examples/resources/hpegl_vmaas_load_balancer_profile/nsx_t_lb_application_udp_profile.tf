# (C) Copyright 2022 Hewlett Packard Enterprise Development LP

# APPLICATION Profile for UDP service
resource "hpegl_vmaas_load_balancer_profile" "tf_lb_profile" {
  lb_id = data.hpegl_vmaas_load_balancer.lb.id  
  name  =  "APPLICATION-UDP-Profile"       
  description  = "creating LB Profile"
  service_type     = "LBFastUdpProfile"
  config {
    profile_type = "application-profile"
    fast_udp_idle_timeout = 300
    ha_flow_mirroring = false
    tags {
        tag = "tag1"
        scope = "scope1"
    }
  }
}