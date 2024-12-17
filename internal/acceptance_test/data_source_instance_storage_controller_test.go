// (C) Copyright 2024 Hewlett Packard Enterprise Development LP

package acceptancetest

import (
	"testing"

	api_client "github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/client"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/atf"
)

func TestAccDataSourceInstanceStorageController(t *testing.T) {
	acc := &atf.Acc{
		PreCheck:     testAccPreCheck,
		Providers:    testAccProviders,
		ResourceName: "hpegl_vmaas_instance_storage_controller",
		GetAPI: func(attr map[string]string) (interface{}, error) {
			cl, cfg := getAPIClient()
			iClient := api_client.InstancesAPIService{
				Client: cl,
				Cfg:    cfg,
			}

			return iClient.GetStorageControllerMount(getAccContext(), attr["layout_id"], attr["controller_name"], toInt(attr["bus_number"]), toInt(attr["interface_number"]))
		},
	}

	acc.RunDataSourceTests(t)
}
