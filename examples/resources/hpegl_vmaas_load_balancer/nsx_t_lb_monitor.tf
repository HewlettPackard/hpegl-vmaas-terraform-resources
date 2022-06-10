# (C) Copyright 2022 Hewlett Packard Enterprise Development LP

resource hpegl_vmaas_load_balancer tf_lb_monitor {
  name  =  "Test-Monitor"       
  description  = "Creating lb monitor"
  monitorType = data.hpegl_vmaas_load_balancer.tf_monitorType.monitorType
  monitorTimeout = 15
  monitorInterval =  30  
  sendVersion      = data.hpegl_vmaas_load_balancer.tf_sendVersion.sendVersion 
  sendType = "GET"          
  receiveCode = "200,300,301,302,304,307"
  monitorDestination = data.hpegl_vmaas_load_balancer.tf_monitorDestination.monitorDestination
  monitorReverse = false
  monitorTransparent = false
  monitorAdaptive = false
  fallCount = 30
  riseCount = 30
  aliasPort = 80
}


