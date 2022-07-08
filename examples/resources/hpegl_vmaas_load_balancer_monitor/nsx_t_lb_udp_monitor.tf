# (C) Copyright 2022 Hewlett Packard Enterprise Development LP

# UDP Monitor
resource "hpegl_vmaas_load_balancer_monitor" "tf_lb_monitor" {
  name  =  "UDP-MONITOR"       
  description  = "Creating lb monitor for UDP"
  type = "LBUdpMonitorProfile"
  fall_count = 3
  interval =  5 
  alias_port = 80
  rise_count = 3
  timeout = 15
  request_body = "request body data"
  response_data = "success"
}