// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package acceptancetest

import (
	"fmt"
	"testing"

	pkgutils "github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/utils"
	"github.com/spf13/viper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceNetworkProxy(t *testing.T) {
	pkgutils.SkipAcc(t, "vmaas.datasource.network_proxy")

	resource.ParallelTest(t, resource.TestCase{
		IsUnitTest: false,
		PreCheck:   func() { testAccPreCheck(t) },
		Providers:  testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNetworkProxyConfig(),
				Check:  validateDataSourceID("data.hpegl_vmaas_network_proxy.tf_proxy"),
			},
		},
	})
}

func testAccDataSourceNetworkProxyConfig() string {
	return providerStanza + fmt.Sprintf(`
	data "hpegl_vmaas_network_proxy" "tf_proxy" {
		name = "%s"
	}
`, viper.GetString("vmaas.datasource.network_proxy.name"))
}
