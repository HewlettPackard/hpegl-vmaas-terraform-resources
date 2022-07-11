# (C) Copyright 2022 Hewlett Packard Enterprise Development LP

# HTTPS Monitor
resource "hpegl_vmaas_load_balancer_monitor" "tf_lb_monitor" {
  lb_id = data.hpegl_vmaas_load_balancer.tf_lb.id
  name  =  "HTTPS-MONITOR"       
  description  = "Creating lb monitor for HTTPS"
  type = "LBHttpsMonitorProfile"
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