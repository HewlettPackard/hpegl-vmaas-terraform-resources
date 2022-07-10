# (C) Copyright 2022 Hewlett Packard Enterprise Development LP

# SSL Profile for CLIENT service
resource "hpegl_vmaas_load_balancer_profile" "tf_lb_profile" {
  lb_id = data.hpegl_vmaas_load_balancer.lb.id  
  name  =  "SSL-CLIENT-Profile"       
  description  = "creating LB Profile"
  service_type     = "LBClientSslProfile"
  config {
    profile_type = "ssl-profile"
    ssl_suite = "Balanced"
    session_cache = true
    session_cache_timeout = 300
    prefer_server_cipher = true
    tags {
        tag = "tag1"
        scope = "scope1"
    }s
  }
}