// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package acceptancetest

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceNetworkInterface(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		IsUnitTest: false,
		PreCheck:   func() { testAccPreCheck(t) },
		Providers:  testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNetworkInterface,
				Check: resource.ComposeTestCheckFunc(
					validateDataSourceID("data.hpegl_vmaas_network_interface.e1000"),
				),
			},
		},
	})
}

const testAccDataSourceNetworkInterface = providerStanza + `
data "hpegl_vmaas_network_interface" "e1000"{
	name = "E1000"
	cloud_id = 1
  }
`
