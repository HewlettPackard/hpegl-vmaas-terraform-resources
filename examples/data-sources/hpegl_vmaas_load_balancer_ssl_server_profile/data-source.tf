(C) Copyright 2022 Hewlett Packard Enterprise Development LP

data "hpegl_vmaas_load_balancer_ssl_server_profile" "tf_ssl_server" {
  lb_id = data.hpegl_vmaas_load_balancer.tf_lb.id
  name = "default-balanced-server-ssl-profile"
}