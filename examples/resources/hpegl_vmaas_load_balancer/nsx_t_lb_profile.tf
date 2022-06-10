# (C) Copyright 2022 Hewlett Packard Enterprise Development LP

resource hpegl_vmaas_load_balancer tf_lb_profile {
  name  =  "LoadBalancer profile"       
  description  = "creating LB Profile"
  serviceType     = data.hpegl_vmaas_lb_profile_service_type.tf_service_type.service_type
  config{
    profileType = data.hpegl_vmaas_lb_profile_profileType.tf_profileType.profileType
    requestHeaderSize = 1024
    responseHeaderSize = 4096
    httpIdleTimeout = 10
    httpIdleTimeout = 15
    fastTcpIdleTimeout = 1800 
    connectionCloseTimeout = 8 
    haFlowMirroring = true
    sslSuite = data.hpegl_vmaas_lb_profile_sslSuite.tf_sslSuite.sslSuite 
    cookieMode = data.hpegl_vmaas_lb_profile_cookieMode.tf_cookieMode.cookieMode
    cookieName = "NSXLB"
    cookieType = data.hpegl_vmaas_lb_profile_cookieType.tf_cookieType.cookieType
    cookieFallback = true
    cookieGarbling = true
  }
}
