# (C) Copyright 2022 Hewlett Packard Enterprise Development LP

data "hpegl_vmaas_lb_profile_service_type" "tf_service_type" {
  service_type = "LBHttpProfile"
}

data "hpegl_vmaas_lb_profile_profileType" "tf_profileType" {
  profileType = "application-profile"
}

data "hpegl_vmaas_lb_profile_sslSuite" "tf_sslSuite" {
  sslSuite = "BALANCED"
}

data "hpegl_vmaas_lb_profile_cookieMode" "tf_cookieMode" {
  cookieMode = "INSERT"
}

data "hpegl_vmaas_lb_profile_cookieType" "tf_cookieType" {
  cookieType = "LBPersistenceCookieTime"
}