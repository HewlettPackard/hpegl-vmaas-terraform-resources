// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package acceptancetest

import (
	"fmt"
	"testing"

	"github.com/spf13/viper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDataStore(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		IsUnitTest: false,
		PreCheck:   func() { testAccPreCheck(t) },
		Providers:  testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDataStoreConfig(),
				Check: resource.ComposeTestCheckFunc(
					validateDataSourceID("data.hpegl_vmaas_datastore." +
						viper.GetString("vmaas.data_source_datastores.dataStoreLocalName")),
				),
			},
		},
	})
}

func testAccDataSourceDataStoreConfig() string {
	return providerStanza + fmt.Sprintf(`
	data "hpegl_vmaas_datastore" "%s" {
		cloud_id = %d
		name = "%s"
	}
`, viper.GetString("vmaas.data_source_datastores.dataStoreLocalName"),
		viper.GetInt("vmaas.data_source_datastores.cloudID"),
		viper.GetString("vmaas.data_source_datastores.dataStoreName"))
}
