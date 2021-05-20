# (C) Copyright 2021 Hewlett Packard Enterprise Development LP

resource "hpegl_vmaas_instance" "test" {
  name          = "test"
  cloud_id      = 1
  group_id      = 1
  plan_id       = 1
  instance_type = "test"
  networks      = [1]
  volumes {
    size         = 5
    datastore_id = "test"

  }
  volumes {
    size         = 10
    datastore_id = "test2"

  }
  labels = ["test"]
  tags = {
    name = "value"
    data = "data"
  }
  config {
    vmware_resource_pool = "test"
  }

  copies = 1

}
