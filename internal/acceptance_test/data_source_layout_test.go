// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package acceptancetest

import (
	"fmt"
	"testing"

	pkgutils "github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/utils"
	"github.com/spf13/viper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceLayout(t *testing.T) {
	pkgutils.SkipAcc(t, "vmaas.datasource.layout")
	resource.ParallelTest(t, resource.TestCase{
		IsUnitTest: false,
		PreCheck:   func() { testAccPreCheck(t) },
		Providers:  testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceLayoutConfig(),
				Check: resource.ComposeTestCheckFunc(
					validateDataSourceID("data.hpegl_vmaas_layout.vmware"),
				),
			},
		},
	})
}

func testAccDataSourceLayoutConfig() string {
	return fmt.Sprintf(`%s
data "hpegl_vmaas_layout" "vmware" {
	name               = "%s"
	instance_type_code = "%s"
}
`, providerStanza,
		viper.GetString("vmaas.datasource.layout.name"),
		viper.GetString("vmaas.datasource.layout.instance_type_code"))
}
