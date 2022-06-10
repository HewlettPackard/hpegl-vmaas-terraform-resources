# (C) Copyright 2022 Hewlett Packard Enterprise Development LP

resource hpegl_vmaas_load_balancer test_lb_pool {
  name  =  "loadbalancer pool"       
  description  = "creating load balancer pool"
  minActive     = data.hpegl_vmaas_lb_pool_minActive.tf_minActive.minActive
  vipBalance = data.hpegl_vmaas_lb_pool_vipBalance.tf_vipBalance.vipBalance
  config{
    snatTranslationType = data.hpegl_vmaas_lb_pool_snatTranslationType.tf_snatTranslationType.snatTranslationType
    passiveMonitorPath = 136
    activeMonitorPaths = 133
    tcpMultiplexing = false
    tcpMultiplexingNumber = 6 
    snatIpAddress = ""
    memberGroup {
        name = "pushpa"
        path = ""
        ipRevisionFilter = "IPV4" 
        port = 80
    }  
  }
}
