// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package acceptancetest

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/spf13/viper"
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
					validateDataSourceID("data.hpegl_vmaas_cloud.test_cloud"),
				),
			},
		},
	})
}

var testAccDataSourceCloud = providerStanza + `
	data "hpegl_vmaas_cloud" "test_cloud" {
		name = "` + viper.GetString("vmaas.datasource.cloud.name") + `"
	}
`
