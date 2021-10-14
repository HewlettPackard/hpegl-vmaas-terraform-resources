// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package acceptancetest

import (
	"fmt"
	"testing"

	pkgutils "github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/utils"
	"github.com/spf13/viper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceNetworkDomain(t *testing.T) {
	pkgutils.SkipAcc(t, "vmaas.datasource.network_domain")

	resource.ParallelTest(t, resource.TestCase{
		IsUnitTest: false,
		PreCheck:   func() { testAccPreCheck(t) },
		Providers:  testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNetworkDomainConfig(),
				Check:  validateDataSourceID("data.hpegl_vmaas_network_domain.tf_domain"),
			},
		},
	})
}

func testAccDataSourceNetworkDomainConfig() string {
	return providerStanza + fmt.Sprintf(`
	data "hpegl_vmaas_network_domain" "tf_domain" {
		name = "%s"
	  }
`, viper.GetString("vmaas.datasource.network_domain.name"))
}
