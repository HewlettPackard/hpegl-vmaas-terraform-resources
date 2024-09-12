// (C) Copyright 2024 Hewlett Packard Enterprise Development LP

package acceptancetest

import (
	"testing"

	api_client "github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/client"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/atf"
)

func TestAccDataSourceInstanceStorageType(t *testing.T) {
	acc := &atf.Acc{
		PreCheck:     testAccPreCheck,
		Providers:    testAccProviders,
		ResourceName: "hpegl_vmaas_instance_disk_type",
		GetAPI: func(attr map[string]string) (interface{}, error) {
			cl, cfg := getAPIClient()
			iClient := api_client.InstancesAPIService{
				Client: cl,
				Cfg:    cfg,
			}

			return iClient.GetStorageVolTypeID(getAccContext(), attr["cloud_id"], attr["layout_id"])
		},
	}

	acc.RunDataSourceTests(t)
}
