// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package acceptancetest

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceGroup(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		IsUnitTest: false,
		PreCheck:   func() { testAccPreCheck(t) },
		Providers:  testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGroup,
				Check: resource.ComposeTestCheckFunc(
					validateDataSourceID("data.hpegl_vmaas_group.dev_group"),
				),
			},
		},
	})
}

const testAccDataSourceGroup = providerStanza + `
data "hpegl_vmaas_group" "dev_group" {
	name = "test_dev_group"
}
`
