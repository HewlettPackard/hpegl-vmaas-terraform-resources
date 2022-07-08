# (C) Copyright 2022 Hewlett Packard Enterprise Development LP

# HTTP Monitor
resource "hpegl_vmaas_load_balancer_monitor" "tf_lb_monitor" {
  name  =  "HTTP-MONITOR"       
  description  = "Creating lb monitor for HTTP"
  type = "LBHttpMonitorProfile"
  fall_count = 3
  interval =  5  
  alias_port = 80
  rise_count = 3
  timeout = 15
  request_body = "request body data"
  request_method = "GET"
  request_url = "https://test.com"
  request_version  = "1.0" 
  response_data = "success"
  response_status_codes  = "201,200"
}