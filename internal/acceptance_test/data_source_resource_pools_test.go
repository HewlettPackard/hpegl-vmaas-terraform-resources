// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package acceptancetest

import (
	"testing"

	api_client "github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/client"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/atf"
)

func TestAccDataSourceResourcePool(t *testing.T) {
	acc := &atf.Acc{
		PreCheck:     testAccPreCheck,
		Providers:    testAccProviders,
		ResourceName: "hpegl_vmaas_resource_pool",
		GetAPI: func(attr map[string]string) (interface{}, error) {
			cl, cfg := getAPIClient()
			iClient := api_client.CloudsAPIService{
				Client: cl,
				Cfg:    cfg,
			}
			id := toInt(attr["id"])
			cloudID := toInt(attr["cloud_id"])

			return iClient.GetSpecificCloudResourcePool(getAccContext(), cloudID, id)
		},
	}

	acc.RunDataSourceTests(t)
}
