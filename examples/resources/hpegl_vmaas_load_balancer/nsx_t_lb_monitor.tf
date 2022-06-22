# (C) Copyright 2022 Hewlett Packard Enterprise Development LP

resource "hpegl_vmaas_load_balancer_monitor" "tf_lb_monitor" {
  name  =  "TEST-MONITOR"       
  description  = "Creating lb monitor"
  monitor_type = "LBTcpMonitorProfile"
  monitor_timeout = 15
  monitor_interval =  30  
  send_version      = 20 
  send_type = "GET"          
  receive_code = "200,300,301,302,304,307"
  monitor_destination = "https://test.com"
  monitor_reverse = false
  monitor_transparent = false
  monitor_adaptive = false
  fall_count = 30
  rise_count = 30
  alias_port = 80
}
