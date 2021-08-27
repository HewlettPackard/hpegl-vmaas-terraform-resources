// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package acceptancetest

import (
	"fmt"
	"testing"

	"github.com/spf13/viper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceTemplate(t *testing.T) {
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
	return providerStanza + fmt.Sprintf(`
data "hpegl_vmaas_template" "%s" {
	name = "%s"
}
`, viper.GetString("vmaas.data_source_templates_test.templateLocalName"),
		viper.GetString("vmaas.data_source_templates_test.templateName"))
}
