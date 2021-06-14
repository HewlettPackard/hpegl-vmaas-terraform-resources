# (C) Copyright 2021 Hewlett Packard Enterprise Development LP

resource "hpegl_vmaas_snapshot" "tf_instance_snapshot" { 
  instance_id = 1382
  name = "snapshot1"
  description = "Example snapshot resource"
}