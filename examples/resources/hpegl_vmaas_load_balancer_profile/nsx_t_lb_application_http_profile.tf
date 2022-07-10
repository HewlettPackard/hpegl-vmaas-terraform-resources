# (C) Copyright 2022 Hewlett Packard Enterprise Development LP

# APPLICATION Profile for HTTP service
resource "hpegl_vmaas_load_balancer_profile" "tf_lb_profile" {
  lb_id = data.hpegl_vmaas_load_balancer.lb.id  
  name  =  "APPLICATION-HTTP-Profile"       
  description  = "creating LB Profile"
  service_type     = "LBHttpProfile"
  config {
    profile_type = "application-profile"
    idle_timeout = 15
    request_header_size = 1024
    response_header_size = 4096
    redirection = "http"
    x_forwarded_for = "INSERT"
    request_body_size = 20
    response_timeout = 60
    ntlm_authentication = false
    tags {
        tag = "tag1"
        scope = "scope1"
    }
  }
}