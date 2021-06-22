// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package acceptancetest

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceCloud(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		IsUnitTest: false,
		PreCheck:   func() { testAccPreCheck(t) },
		Providers:  testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCloud,
				Check: resource.ComposeTestCheckFunc(
					checkVmaasCloud("data.hpegl_vmaas_cloud.test_cloud"),
				),
			},
		},
	})
}

const testAccDataSourceCloud = providerStanza + `
	data "hpegl_vmaas_cloud" "test_cloud" {
		name = "HPE GreenLake VMaaS Cloud"
	}
`

func checkVmaasCloud(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Data source not found %s", name)
		}

		id := rs.Primary.Attributes["id"]
		if id == "" {
			return fmt.Errorf("Data source %s is not set", name)
		}

		return nil
	}
}
