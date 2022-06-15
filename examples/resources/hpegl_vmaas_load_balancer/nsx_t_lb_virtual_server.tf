# (C) Copyright 2022 Hewlett Packard Enterprise Development LP

resource hpegl_vmaas_load_balancer_virtual_server tf_lb_virtual_server {
  vip_name  =  "load balancer virtual server"       
  description  = "creating load balancer virtual server"
  vip_address     = "10.11.12.13"
  vip_port = "80"
  pool = 40
  ssl_cert = 22 
  ssl_server_cert =  22
  config{
    persistence = data.hpegl_vmaas_load_balancer_virtual_server.tf_persistence.persistence
    persistence_profile = 582
    application_profile = 605
    ssl_client_profile = 631
    ssl_server_profile = 635
  }
}
