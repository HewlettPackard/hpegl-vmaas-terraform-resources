# (C) Copyright 2022 Hewlett Packard Enterprise Development LP

data hpegl_vmaas_lb_pool_vipBalance tf_vipBalance {
  vipBalance = "ROUND_ROBIN"
}

data hpegl_vmaas_lb_pool_minActive tf_minActive {
  minActive = 1
}

data hpegl_vmaas_lb_pool_snatTranslationType tf_snatTranslationType {
  snatTranslationType = "LBSnatAutoMap"
}

