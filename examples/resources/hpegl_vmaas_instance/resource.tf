# (C) Copyright 2021 Hewlett Packard Enterprise Development LP

resource "hpegl_vmaas_instance" "tf_instance" {
  name          = "instance_tf"
  cloud_id      = data.hpegl_vmaas_cloud.cloud.id
  group_id      = data.hpegl_vmaas_group.default_group.id
  layout_id     = data.hpegl_vmaas_layout.vmware.id
  plan_id       = data.hpegl_vmaas_plan.g1_small.id
  instance_code = data.hpegl_vmaas_layout.vmware.instance_code
  networks {
    id = data.hpegl_vmaas_network.blue_net.id
  }

  volumes {
    name         = "root"
    size         = 5
    datastore_id = "auto"
  }

  labels = ["test"]
  tags = {
    name = "vmaas"
    data = "test_vm"
  }
  config {
    resource_pool_id = data.hpegl_vmaas_resourcePool.cluster.id
    template         = "apache-centos7-x86_64-09072020"
  }

  copies = 1
}
