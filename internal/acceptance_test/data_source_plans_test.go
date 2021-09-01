// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package acceptancetest

import (
	"fmt"
	"testing"

	"github.com/spf13/viper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourcePlan(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		IsUnitTest: false,
		PreCheck:   func() { testAccPreCheck(t) },
		Providers:  testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePlanConfig(),
				Check: resource.ComposeTestCheckFunc(
					validateDataSourceID("data.hpegl_vmaas_plan." + viper.GetString("vmaas.data_source_plans.planLocalName")),
				),
			},
		},
	})
}

func testAccDataSourcePlanConfig() string {
	return providerStanza + fmt.Sprintf(`
data "hpegl_vmaas_plan" "%s" {
	name = "%s"
}
`, viper.GetString("vmaas.data_source_plans.planLocalName"), viper.GetString("vmaas.data_source_plans.planName"))
}
