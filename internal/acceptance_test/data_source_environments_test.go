// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package acceptancetest

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceEnvironment(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		IsUnitTest: false,
		PreCheck:   func() { testAccPreCheck(t) },
		Providers:  testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceEnvironment,
				Check: resource.ComposeTestCheckFunc(
					validateDataSourceID("data.hpegl_vmaas_environment.dev"),
				),
			},
		},
	})
}

const testAccDataSourceEnvironment = providerStanza + `
data "hpegl_vmaas_environment" "dev" {
	name = "Dev"
  }
`
