# (C) Copyright 2022 Hewlett Packard Enterprise Development LP

data hpegl_vmaas_load_balancer_pool tf_vipBalance {
  vip_balance = "ROUND_ROBIN"
}

data hpegl_vmaas_load_balancer_pool tf_minActive {
  min_active = 1
}

data hpegl_vmaas_load_balancer_pool tf_snatTranslationType {
  snat_translation_type = "LBSnatAutoMap"
}