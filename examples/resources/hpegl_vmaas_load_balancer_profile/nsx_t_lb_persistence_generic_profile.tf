# (C) Copyright 2022 Hewlett Packard Enterprise Development LP

# PERSISTENCE Profile for GENERIC service
resource "hpegl_vmaas_load_balancer_profile" "tf_lb_profile" {
  lb_id = data.hpegl_vmaas_load_balancer.lb.id  
  name  =  "PERSISTENCE-GENERIC-Profile"       
  description  = "creating LB Profile"
  service_type     = "LBGenericPersistenceProfile"
  config {
    profile_type = "persistence-profile"
    share_persistence = false
    ha_persistence_mirroring = false
    persistence_entry_timeout - 300
    tags {
        tag = "tag1"
        scope = "scope1"
    }s
  }
}