# (C) Copyright 2022 Hewlett Packard Enterprise Development LP

# ICMP Monitor
resource "hpegl_vmaas_load_balancer_monitor" "tf_lb_monitor" {
  name  =  "ICMP-MONITOR"       
  description  = "Creating lb monitor for ICMP"
  type = "LBIcmpMonitorProfile"
  fall_count = 3
  interval =  5  
  alias_port = 80
  rise_count = 3
  timeout = 15
  data_length = 56
}