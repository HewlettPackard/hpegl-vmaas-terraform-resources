# (C) Copyright 2022 Hewlett Packard Enterprise Development LP

# ICMP Monitor
resource "hpegl_vmaas_load_balancer_monitor" "tf_lb_monitor" {
  lb_id = data.hpegl_vmaas_load_balancer.tf_lb.id
  name  =  "tf_ICMP_MONITOR"       
  description  = "ICMP_MONITOR create using tf"
  type = "LBIcmpMonitorProfile"
  fall_count = 3
  interval =  5  
  monitor_port = 80
  rise_count = 3
  timeout = 15
  data_length = 56
}