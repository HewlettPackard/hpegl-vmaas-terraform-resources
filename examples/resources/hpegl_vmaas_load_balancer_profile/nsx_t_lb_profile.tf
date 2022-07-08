# (C) Copyright 2022 Hewlett Packard Enterprise Development LP

resource "hpegl_vmaas_load_balancer_profile" "tf_lb_profile" {
  lb_id = data.hpegl_vmaas_load_balancer.lb.id  
  name  =  "LB-PROFILE"       
  description  = "creating LB Profile"
  service_type     = "LBHttpProfile"
  config {
    profile_type = "application-profile"
    request_header_size = 1024
    response_header_size = 4096
    http_idle_timeout = 10
    fast_tcp_idle_timeout = 1800 
    connection_close_timeout = 8 
    ha_flow_mirroring = true
    ssl_suite = "CUSTOM" 
    cookie_mode = "INSERT"
    cookie_name = "NSXLB"
    cookie_type = "LBSessionCookieTime"
    cookie_fallback = true
    cookie_garbling = true
  }
}
