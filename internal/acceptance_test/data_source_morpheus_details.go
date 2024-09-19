// (C) Copyright 2024 Hewlett Packard Enterprise Development LP

package acceptancetest

import (
	"testing"

	api_client "github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/client"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/atf"
)

func TestAccDataSourceMorpheusDetails(t *testing.T) {
	acc := &atf.Acc{
		PreCheck:     testAccPreCheck,
		Providers:    testAccProviders,
		ResourceName: "hpegl_vmaas_morpheus_details",
		GetAPI: func(attr map[string]string) (interface{}, error) {
			cl, cfg := getBrokerAPIClient()
			iClient := api_client.BrokerAPIService{
				Client: cl,
				Cfg:    cfg,
			}

			return iClient.GetMorpheusDetails(getAccContext())
		},
	}

	acc.RunDataSourceTests(t)
}
