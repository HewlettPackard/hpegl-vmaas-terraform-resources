# (C) Copyright 2022 Hewlett Packard Enterprise Development LP

resource "hpegl_vmaas_load_balancer_virtual_server" "tf_lb_virtual_server" {
  lb_id = 248
  name  =  "virtual"       
  description  = "creating load balancer virtual server"
  vip_address     = "10.11.12.13"
  type = "http"
  vip_port = "80"
  vip_host_name = "host VS"
  pool = data.hpegl_vmaas_load_balancer_pool.tf_pool.id
  ssl_client_cert = 22
  ssl_server_cert = 22
  config{
    persistence = "COOKIE"
    persistence_profile = data.hpegl_vmaas_load_balancer_persistence_profile.tf_cookie_profile.id
    application_profile = data.hpegl_vmaas_load_balancer_application_profile.tf_app_profile.id
    ssl_client_profile = data.hpegl_vmaas_load_balancer_ssl_client_profile.tf_ssl_client.id
    ssl_server_profile = data.hpegl_vmaas_load_balancer_ssl_server_profile.tf_ssl_server.id
  }
}
