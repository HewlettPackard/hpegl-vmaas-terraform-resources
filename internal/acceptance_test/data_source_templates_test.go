// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package acceptancetest

import (
	"fmt"
	"testing"

	pkgutils "github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/utils"
	"github.com/spf13/viper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceTemplate(t *testing.T) {
	pkgutils.SkipAcc(t, "vmaas.datasource.template")
	resource.ParallelTest(t, resource.TestCase{
		IsUnitTest: false,
		PreCheck:   func() { testAccPreCheck(t) },
		Providers:  testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTemplateConfig(),
				Check: resource.ComposeTestCheckFunc(
					validateDataSourceID("data.hpegl_vmaas_template.vanilla"),
				),
			},
		},
	})
}

func testAccDataSourceTemplateConfig() string {
	return fmt.Sprintf(`%s
data "hpegl_vmaas_template" "vanilla" {
	name = "%s"
}
`,
		providerStanza,
		viper.GetString("vmaas.datasource.template.name"))
}
