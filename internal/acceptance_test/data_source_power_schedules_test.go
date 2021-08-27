// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package acceptancetest

import (
	"fmt"
	"testing"

	"github.com/spf13/viper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourcePowerSchedule(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		IsUnitTest: false,
		PreCheck:   func() { testAccPreCheck(t) },
		Providers:  testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePowerScheduleConfig(),
				Check: resource.ComposeTestCheckFunc(
					validateDataSourceID("data.hpegl_vmaas_power_schedule." +
						viper.GetString("vmaas.data_source_power_schedules_test.powerScheduleLocalName")),
				),
			},
		},
	})
}

func testAccDataSourcePowerScheduleConfig() string {
	return providerStanza + fmt.Sprintf(`
data "hpegl_vmaas_power_schedule" "%s" {
	name = "%s"
}
`, viper.GetString("vmaas.data_source_power_schedules_test.powerScheduleLocalName"),
		viper.GetString("vmaas.data_source_power_schedules_test.powerScheduleName"))
}
