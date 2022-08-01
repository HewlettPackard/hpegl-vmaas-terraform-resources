(C) Copyright 2022 Hewlett Packard Enterprise Development LP

data "hpegl_vmaas_load_balancer_profile" "tf_udp_profile" {
  lb_id = data.hpegl_vmaas_load_balancer.tf_lb.id
  name = "default-udp-lb-app-profile"
}