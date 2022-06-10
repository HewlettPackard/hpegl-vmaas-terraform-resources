# (C) Copyright 2022 Hewlett Packard Enterprise Development LP

resource hpegl_vmaas_load_balancer tf_lb_virtual_server {
  vipName  =  "load balancer virtual server"       
  description  = "creating load balancer virtual server"
  vipAddress     = "10.11.12.13"
  vipPort = "80"
  vipProtocol = data.hpegl_vmaas_load_balancer_vipProtocol.tf_vipProtocol.vipProtocol
  pool = 40
  sslCert = 22 
  sslServerCert =  22
  config{
    persistence = data.hpegl_vmaas_load_balancer_persistence.tf_persistence.persistence
    persistenceProfile = 582
    applicationProfile = 605
    sslClientProfile = 631
    sslServerProfile = 635
  }
}
