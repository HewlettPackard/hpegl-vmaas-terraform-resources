// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package acceptancetest

import (
	"fmt"
	"testing"

	"github.com/spf13/viper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceNetworkInterface(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		IsUnitTest: false,
		PreCheck:   func() { testAccPreCheck(t) },
		Providers:  testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNetworkInterfaceConfig(),
				Check: resource.ComposeTestCheckFunc(
					validateDataSourceID("data.hpegl_vmaas_network_interface." +
						viper.GetString("vmaas.data_source_network_interface_test.networkInterfaceLocalName")),
				),
			},
		},
	})
}

func testAccDataSourceNetworkInterfaceConfig() string {
	return providerStanza + fmt.Sprintf(`
data "hpegl_vmaas_network_interface" "%s"{
	name = "%s"
	cloud_id = %d
  }
`, viper.GetString("vmaas.data_source_network_interface_test.networkInterfaceLocalName"),
		viper.GetString("vmaas.data_source_network_interface_test.networkInterfaceName"),
		viper.GetInt("vmaas.data_source_network_interface_test.networkInterfaceCloudID"))
}
