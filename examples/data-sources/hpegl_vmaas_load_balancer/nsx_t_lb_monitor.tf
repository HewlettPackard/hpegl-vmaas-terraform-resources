# (C) Copyright 2022 Hewlett Packard Enterprise Development LP

data "hpegl_vmaas_load_balancer" "tf_monitorType" {
  monitorType = "LBTcpMonitorProfile
}

data "hpegl_vmaas_load_balancer" "tf_sendVersion" {
  sendVersion = "HTTP_VERSION_1_1"
}

data "hpegl_vmaas_load_balancer" "tf_monitorDestination" {
  monitorDestination = "https=//qa.test.com"
}