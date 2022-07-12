 (C) Copyright 2022 Hewlett Packard Enterprise Development LP

data "hpegl_vmaas_passive_monitor" "tf_lb_monitor" {
  lb_id = data.hpegl_vmaas_load_balancer.tf_lb.id
  type = "LBPassiveMonitorProfile"
} 
