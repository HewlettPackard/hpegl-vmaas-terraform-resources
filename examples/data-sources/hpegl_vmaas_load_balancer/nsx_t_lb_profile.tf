# (C) Copyright 2022 Hewlett Packard Enterprise Development LP

data "hpegl_vmaas_load_balancer_profile" "tf_service_type" {
  service_type = "LBHttpProfile"
}

data "hpegl_vmaas_load_balancer_profile" "tf_profileType" {
  profile_type = "application-profile"
}

data "hpegl_vmaas_load_balancer_profile" "tf_sslSuite" {
  ssl_suite = "BALANCED"
}

data "hpegl_vmaas_load_balancer_profile" "tf_cookieMode" {
  cookie_mode = "INSERT"
}

data "hpegl_vmaas_load_balancer_profile" "tf_cookieType" {
  cookie_type = "LBPersistenceCookieTime"
}