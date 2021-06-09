
data "hpegl_vmaas_resourcePool" "cluster" {
  cloud_id = data.hpegl_vmaas_cloud.cloud.id
  name     = "Cluster"
}
