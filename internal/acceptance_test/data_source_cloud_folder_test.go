// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package acceptancetest

import (
	"fmt"
	"testing"

	pkgutils "github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/utils"
	"github.com/spf13/viper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceCloudFolder(t *testing.T) {
	pkgutils.SkipAcc(t, "vmaas.datasource.cloud_folder")

	resource.ParallelTest(t, resource.TestCase{
		IsUnitTest: false,
		PreCheck:   func() { testAccPreCheck(t) },
		Providers:  testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCloudFolderConfig(),
				Check:  validateDataSourceID("data.hpegl_vmaas_cloud_folder.compute_folder"),
			},
		},
	})
}

func testAccDataSourceCloudFolderConfig() string {
	return providerStanza + fmt.Sprintf(`
	data "hpegl_vmaas_cloud_folder" "compute_folder" {
		cloud_id = %d
		name     = "%s"
	  }
`, viper.GetInt("vmaas.datasource.cloud_folder.cloud_id"),
		viper.GetString("vmaas.datasource.cloud_folder.name"))
}
