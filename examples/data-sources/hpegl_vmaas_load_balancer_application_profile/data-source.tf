(C) Copyright 2022 Hewlett Packard Enterprise Development LP

data "hpegl_vmaas_load_balancer_application_profile" "tf_app_profile" {
  lb_id = data.hpegl_vmaas_load_balancer.tf_lb.id
  name = "APPLICATION-HTTP-Profile"
}
