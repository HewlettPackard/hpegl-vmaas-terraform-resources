# (C) Copyright 2022 Hewlett Packard Enterprise Development LP

data "hpegl_vmaas_load_balancer_monitor" "tf_monitorType" {
  monitor_type = "LBTcpMonitorProfile
}

data "hpegl_vmaas_load_balancer_monitor" "tf_sendVersion" {
  send_version = "HTTP_VERSION_1_1"
}

data "hpegl_vmaas_load_balancer_monitor" "tf_monitorDestination" {
  monitor_destination = "https=//qa.test.com"
}