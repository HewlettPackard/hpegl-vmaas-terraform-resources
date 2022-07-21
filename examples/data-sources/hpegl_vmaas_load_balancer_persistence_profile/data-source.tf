(C) Copyright 2022 Hewlett Packard Enterprise Development LP

// Cookie based Persistence Profile
data "hpegl_vmaas_load_balancer_persistence_profile" "tf_cookie_profile" {
  lb_id = data.hpegl_vmaas_load_balancer.tf_lb.id
  name = "default-cookie-lb-persistence-profile"
}

// Source IP based Persistence Profile
data "hpegl_vmaas_load_balancer_persistence_profile" "tf_source_profile" {
  lb_id = data.hpegl_vmaas_load_balancer.tf_lb.id
  name = "default-source-ip-lb-persistence-profile"
}