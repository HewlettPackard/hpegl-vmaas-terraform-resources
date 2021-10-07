// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package acceptancetest

import (
	"fmt"
	"testing"

	pkgutils "github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/utils"
	"github.com/spf13/viper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceRouter(t *testing.T) {
	pkgutils.SkipAcc(t, "vmaas.datasource.router")
	resource.ParallelTest(t, resource.TestCase{
		IsUnitTest: false,
		PreCheck:   func() { testAccPreCheck(t) },
		Providers:  testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRouterConfig(),
				Check: resource.ComposeTestCheckFunc(
					validateDataSourceID("data.hpegl_vmaas_router.tier0_router"),
				),
			},
		},
	})
}

func testAccDataSourceRouterConfig() string {
	return fmt.Sprintf(`%s
data "hpegl_vmaas_router" "tier0_router" {
	name = "%s"
}
`, providerStanza,
		viper.GetString("vmaas.datasource.router.name"))
}
