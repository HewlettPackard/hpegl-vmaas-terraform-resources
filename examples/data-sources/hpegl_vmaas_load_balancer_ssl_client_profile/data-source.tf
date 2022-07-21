(C) Copyright 2022 Hewlett Packard Enterprise Development LP

data "hpegl_vmaas_load_balancer_ssl_client_profile" "tf_ssl_client" {
  lb_id = data.hpegl_vmaas_load_balancer.tf_lb.id
  name = "default-balanced-client-ssl-profile"
}