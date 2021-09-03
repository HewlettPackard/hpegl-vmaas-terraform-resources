// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package acceptancetest

import (
	"fmt"
	"testing"

	"github.com/spf13/viper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceLayout(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		IsUnitTest: false,
		PreCheck:   func() { testAccPreCheck(t) },
		Providers:  testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceLayoutConfig(),
				Check: resource.ComposeTestCheckFunc(
					validateDataSourceID("data.hpegl_vmaas_layout." + viper.GetString("vmaas.data_source_layout.layoutLocalName")),
				),
			},
		},
	})
}

func testAccDataSourceLayoutConfig() string {
	return providerStanza + fmt.Sprintf(`
data "hpegl_vmaas_layout" "%s" {
	name               = "%s"
	instance_type_code = "%s"
}
`, viper.GetString("vmaas.data_source_layout.layoutLocalName"),
		viper.GetString("vmaas.data_source_layout.layoutName"),
		viper.GetString("vmaas.data_source_layout.layoutInstanceTypeCode"))
}
