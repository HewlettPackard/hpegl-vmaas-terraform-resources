// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package acceptancetest

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceResourcePool(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		IsUnitTest: false,
		PreCheck:   func() { testAccPreCheck(t) },
		Providers:  testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceResourcePool,
				Check: resource.ComposeTestCheckFunc(
					validateDataSourceID("data.hpegl_vmaas_resource_pool.cl_resource_pool"),
				),
			},
		},
	})
}

const testAccDataSourceResourcePool = providerStanza + `
data "hpegl_vmaas_resource_pool" "cl_resource_pool" {
	cloud_id = data.hpegl_vmaas_cloud.cloud.id
	name     = "Cluster"
}
`
