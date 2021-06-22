// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package acceptancetest

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDataStore(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		IsUnitTest: false,
		PreCheck:   func() { testAccPreCheck(t) },
		Providers:  testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDataStore,
				Check: resource.ComposeTestCheckFunc(
					validateDataSourceID("data.hpegl_vmaas_datastore.glhc_vol10"),
				),
			},
		},
	})
}

const testAccDataSourceDataStore = providerStanza + `
	data "hpegl_vmaas_datastore" "glhc_vol10" {
		cloud_id = 1
		name = "GLHC-Vol10"
	}
`
