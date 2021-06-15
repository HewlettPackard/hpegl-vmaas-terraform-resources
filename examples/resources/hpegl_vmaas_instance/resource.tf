# (C) Copyright 2021 Hewlett Packard Enterprise Development LP

#  Set-up for terraform >= v0.13
terraform {
  required_providers {
    hpegl = {
      source  = "terraform.example.com/vmaas/hpegl"
      version = ">= 0.0.1"
    }
  }
}

provider "hpegl" {
  vmaas {
    location   = "AGENA-DEV1-INTG"
    space_name = "Default"
  }
  iam_token = "<<I_AM Token>>"
}

data "hpegl_vmaas_datastore" "c_3par" {
  cloud_id = data.hpegl_vmaas_cloud.cloud.id
  name     = "Compute-3par-A64G-FC-1TB"
}

data "hpegl_vmaas_network" "blue_net" {
  name = "Blue-Net"
}

data "hpegl_vmaas_network" "green_net" {
  name = "Green-Net"
}

data "hpegl_vmaas_group" "default_group" {
  name = "Default"
}

data "hpegl_vmaas_resource_pool" "cl_resource_pool" {
  cloud_id = data.hpegl_vmaas_cloud.cloud.id
  name     = "Cluster"
}

data "hpegl_vmaas_layout" "vmware" {
  name               = "Vmware VM"
  instance_type_code = "vmware"
}

data "hpegl_vmaas_layout" "vmware_centos" {
  name               = "VMware VM with vanilla CentOS"
  instance_type_code = "glhc-vanilla-centos"
}

data "hpegl_vmaas_cloud" "cloud" {
  name = "HPE GreenLake VMaaS Cloud"
}

data "hpegl_vmaas_plan" "g1_small" {
  name = "G1-Small"
}

data "hpegl_vmaas_power_schedule" "weekday" {
  name = "DEMO_WEEKDAY"
}

data "hpegl_vmaas_template" "vanilla" {
  name = "vanilla-centos7-x86_64-09072020"
}

data "hpegl_vmaas_environment" "dev" {
  name = "Dev"
}

# minimal instance creation
resource "hpegl_vmaas_instance" "minimal_instance" {
  name               = "tf_minimal"
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
    resource_pool_id = data.hpegl_vmaas_resource_pool.cl_resource_pool.id
    vm_folder        = "group-v140"
  }
  environment_code = data.hpegl_vmaas_environment.dev.code
}

# create instance will all possible options
resource "hpegl_vmaas_instance" "tf_instance" {
  name               = "tf_advanced"
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
    resource_pool_id = data.hpegl_vmaas_resource_pool.cl_resource_pool.id
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
  env_prefix        = "tf_test"
  power_schedule_id = data.hpegl_vmaas_power_schedule.weekday.id
  port {
    name = "nginx"
    port = 80
    lb   = "No LB"
  }
  environment_code = data.hpegl_vmaas_environment.dev.code
  # On creating only poweron operation is supported. Upon updation all other
  # lifecycle operations are permitted.
  power= "poweron"
}

# Clone a instance from an existing instance
resource "hpegl_vmaas_instance" "tf_instance_clone" {
  name               = "tf_clone"
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

  volume {
    name         = "Local_vol"
    size         = 5
    datastore_id = data.hpegl_vmaas_datastore.c_3par.id
    root         = false
  }

  config {
    resource_pool_id = data.hpegl_vmaas_resource_pool.cl_resource_pool.id
    template_id      = data.hpegl_vmaas_template.vanilla.id
    no_agent         = true
    vm_folder        = "group-v140"
    # create_user      = true
    # asset_tag        = "vm_tag"
  }
  hostname = "hpegl_tf_host_clone"
  scale    = 2
  evars = {
    proxy = "http://some:proxy"
  }
  # power_schedule_id = data.hpegl_vmaas_powerSchedule.weekday.id
  clone {
    source_instance_id = hpegl_vmaas_instance.tf_instance.id
  }
}
