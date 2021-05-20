// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

provider "hpegl" {
  vmaas {
    location   = "location"
    space_name = "space_name"
  }
  iam_token = "iam-token"
}
