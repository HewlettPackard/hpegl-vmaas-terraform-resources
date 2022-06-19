# (C) Copyright 2022 Hewlett Packard Enterprise Development LP

resource "hpegl_vmaas_load_balancer_virtual_server" "tf_lb_virtual_server" {
  vip_name  =  "virtual"       
  description  = "creating load balancer virtual server"
  vip_address     = "10.11.12.13"
  vip_port = "80"
  pool = 266
  ssl_cert = 0 
  ssl_server_cert =  0
  config{
    persistence = "COOKIE"
    application_profile = 2929
  }
}
