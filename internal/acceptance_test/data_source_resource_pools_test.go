// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package acceptancetest

import (
	"fmt"
	"testing"

	"github.com/spf13/viper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceResourcePool(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		IsUnitTest: false,
		PreCheck:   func() { testAccPreCheck(t) },
		Providers:  testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceResourcePoolConfig(),
				Check: resource.ComposeTestCheckFunc(
					validateDataSourceID("data.hpegl_vmaas_resource_pool." +
						viper.GetString("vmaas.data_source_resource_pools.resourcePoolLocalName")),
				),
			},
		},
	})
}

func testAccDataSourceResourcePoolConfig() string {
	return providerStanza + fmt.Sprintf(`
data "hpegl_vmaas_resource_pool" "%s" {
	cloud_id = %d
	name     = "%s"
}
`, viper.GetString("vmaas.data_source_resource_pools.resourcePoolLocalName"),
		viper.GetInt("vmaas.data_source_resource_pools.resourcePoolCloudID"),
		viper.GetString("vmaas.data_source_resource_pools.resourcePoolName"))
}
