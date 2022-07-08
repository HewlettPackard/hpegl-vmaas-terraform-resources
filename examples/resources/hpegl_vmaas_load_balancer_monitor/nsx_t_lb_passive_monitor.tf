# (C) Copyright 2022 Hewlett Packard Enterprise Development LP

# PASSIVE Monitor
resource "hpegl_vmaas_load_balancer_monitor" "tf_lb_monitor" {
  name  =  "PASSIVE-MONITOR"       
  description  = "Creating lb monitor for PASSIVE"
  type = "LBPassiveMonitorProfile"
  timeout = 15
  data_length = 56
  max_fail = 5
}
