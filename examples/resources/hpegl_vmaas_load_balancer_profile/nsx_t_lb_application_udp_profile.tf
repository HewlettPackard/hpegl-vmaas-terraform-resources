# (C) Copyright 2022 Hewlett Packard Enterprise Development LP

# APPLICATION Profile for UDP service
resource "hpegl_vmaas_load_balancer_profile" "tf_APPLICATION-UDP" {
  lb_id = data.hpegl_vmaas_load_balancer.tf_lb.id  
  name  =  "tf_APPLICATION-UDP"       
  description  = "tf_APPLICATION-UDP updating using tf"
  service_type     = "LBFastUdpProfile"
  config {
    udp_profile{
      profile_type = "application-profile"
      fast_udp_idle_timeout = 30
      ha_flow_mirroring = true
    }
  }
}