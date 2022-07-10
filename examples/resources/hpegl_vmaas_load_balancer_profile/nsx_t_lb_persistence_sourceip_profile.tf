# (C) Copyright 2022 Hewlett Packard Enterprise Development LP

# PERSISTENCE Profile for SOURCEIP service
resource "hpegl_vmaas_load_balancer_profile" "tf_lb_profile" {
  lb_id = data.hpegl_vmaas_load_balancer.lb.id  
  name  =  "PERSISTENCE-SOURCEIP-Profile"       
  description  = "creating LB Profile"
  service_type     = "LBSourceIpPersistenceProfile"
  config {
    profile_type = "persistence-profile"
    share_persistence = false
    ha_persistence_mirroring = false
    persistence_entry_timeout - 300
    purge_entries_when_full = true
    tags {
        tag = "tag1"
        scope = "scope1"
    }s
  }
}