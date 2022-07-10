# (C) Copyright 2022 Hewlett Packard Enterprise Development LP

# PERSISTENCE Profile for COOKIE service
resource "hpegl_vmaas_load_balancer_profile" "tf_lb_profile" {
  lb_id = data.hpegl_vmaas_load_balancer.lb.id  
  name  =  "PERSISTENCE-COOKIE-Profile"       
  description  = "creating LB Profile"
  service_type     = "LBCookiePersistenceProfile"
  config {
    profile_type = "persistence-profile"
    cookie_name = "cookie1"
    cookie_fallback = true
    cookie_garbling = true
    cookie_mode = "INSERT"
    cookie_type = "sessioncookie"
    cookie_domain = "domain1"
    cookie_path = "http://cookie.com"
    max_idle_time = 60
    max_cookie_age = 2
    share_persistence = true
    tags {
        tag = "tag1"
        scope = "scope1"
    }s
  }
}