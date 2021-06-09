
# minimal instance creation
resource "hpegl_vmaas_instance" "minimal_instance" {
  name               = "shihad_tf_minimal"
  cloud_id           = data.hpegl_vmaas_cloud.cloud.id
  group_id           = data.hpegl_vmaas_group.default_group.id
  layout_id          = data.hpegl_vmaas_layout.vmware_centos.id
  plan_id            = data.hpegl_vmaas_plan.g1_small.id
  instance_type_code = data.hpegl_vmaas_layout.vmware_centos.instance_type_code
  network {
    id = data.hpegl_vmaas_network.blue_net.id
  }

  volume {
    name         = "root_vol"
    size         = 5
    datastore_id = data.hpegl_vmaas_datastore.c_3par.id
    root         = true
  }

  config {
    resource_pool_id = data.hpegl_vmaas_resourcePool.cl_resourcePool.id
    no_agent         = true
    vm_folder        = "group-v140"
  }
}

# create instance will all possible options
resource "hpegl_vmaas_instance" "tf_instance" {
  name               = "shihad_tf_new"
  cloud_id           = data.hpegl_vmaas_cloud.cloud.id
  group_id           = data.hpegl_vmaas_group.default_group.id
  layout_id          = data.hpegl_vmaas_layout.vmware_centos.id
  plan_id            = data.hpegl_vmaas_plan.g1_small.id
  instance_type_code = data.hpegl_vmaas_layout.vmware_centos.instance_type_code
  network {
    id = data.hpegl_vmaas_network.blue_net.id
  }
  network {
    id = data.hpegl_vmaas_network.green_net.id
  }

  volume {
    name         = "root_vol"
    size         = 5
    datastore_id = data.hpegl_vmaas_datastore.c_3par.id
    root         = true
  }

  volume {
    name         = "Local_vol"
    size         = 5
    datastore_id = data.hpegl_vmaas_datastore.c_3par.id
    root         = false
  }

  labels = ["test_label"]
  tags = {
    key  = "value"
    name = "data"
    some = "fdsfs"
  }

  config {
    resource_pool_id = data.hpegl_vmaas_resourcePool.cl_resourcePool.id
    template_id      = data.hpegl_vmaas_template.vanilla.id
    no_agent         = true
    vm_folder        = "group-v140"
    create_user      = true
    asset_tag        = "vm_tag"
  }
  hostname = "hpegl_tf_host"
  scale    = 2
  evars = {
    proxy = "http://some:proxy"
  }
  power_schedule_id = data.hpegl_vmaas_powerSchedule.weekday.id
}


# Clone a instance from an existing instance
resource "hpegl_vmaas_instance" "tf_instance_clone" {
  name               = "shihad_tf_clone"
  cloud_id           = data.hpegl_vmaas_cloud.cloud.id
  group_id           = data.hpegl_vmaas_group.default_group.id
  layout_id          = data.hpegl_vmaas_layout.vmware.id
  plan_id            = data.hpegl_vmaas_plan.g1_small.id
  instance_type_code = data.hpegl_vmaas_layout.vmware.instance_type_code
  network {
    id = data.hpegl_vmaas_network.blue_net.id
  }
  network {
    id = data.hpegl_vmaas_network.green_net.id
  }

  volume {
    name         = "root_vol"
    size         = 5
    datastore_id = data.hpegl_vmaas_datastore.c_3par.id
    root         = true
  }
  labels = ["test"]
  tags = {
    name = "value"
    data = "data"
  }
  config {
    resource_pool_id = data.hpegl_vmaas_resourcePool.cl_resourcePool.id
    no_agent         = true
    vm_folder        = "group-v140"
  }

  scale = 1
  clone {
    source_instance_id = hpegl_vmaas_instance.tf_instance.id
  }
}
