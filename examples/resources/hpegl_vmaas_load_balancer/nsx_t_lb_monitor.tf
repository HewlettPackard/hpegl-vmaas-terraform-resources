# (C) Copyright 2022 Hewlett Packard Enterprise Development LP

resource "hpegl_vmaas_load_balancer_monitor" "tf_lb_monitor" {
  name  =  "Test-Monitor"       
  description  = "Creating lb monitor"
  monitor_type = "data.hpegl_vmaas_load_balancer_monitor.tf_monitorType.monitor_type"
  monitor_timeout = 15
  monitor_interval =  30  
  send_version      = "data.hpegl_vmaas_load_balancer_monitor.tf_sendVersion.send_version "
  send_type = "GET"          
  receive_code = "200,300,301,302,304,307"
  monitor_destination = "data.hpegl_vmaas_load_balancer_monitor.tf_monitorDestination.monitor_destination"
  monitor_reverse = false
  monitor_transparent = false
  monitor_adaptive = false
  fall_count = 30
  rise_count = 30
  alias_port = 80
}


