# (C) Copyright 2022 Hewlett Packard Enterprise Development LP

resource hpegl_vmaas_load_balancer_profile tf_lb_profile {
  name  =  "LoadBalancer profile"       
  description  = "creating LB Profile"
  service_type     = "data.hpegl_vmaas_load_balancer_profile_service_type.tf_service_type.service_type"
  config{
    profile_type = data.hpegl_vmaas_lb_profile_profileType.tf_profileType.profile_type
    request_header_size = 1024
    response_header_size = 4096
    http_idle_timeout = 10
    fast_tcp_idle_timeout = 1800 
    connection_close_timeout = 8 
    ha_flow_mirroring = true
    ssl_suite = data.hpegl_vmaas_load_balancer_profile_sslSuite.tf_sslSuite.ssl_suite 
    cookie_mode = data.hpegl_vmaas_load_balancer_profile_cookieMode.tf_cookieMode.cookie_mode
    cookie_name = "NSXLB"
    cookie_type = data.hpegl_vmaas_load_balancer_profile_cookieType.tf_cookieType.cookie_type
    cookie_fallback = true
    cookie_garbling = true
  }
}
