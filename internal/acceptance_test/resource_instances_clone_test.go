// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package acceptancetest

import (
	"testing"

	api_client "github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/client"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/atf"
)

func TestVmaasInstanceClonePlan(t *testing.T) {
	acc := &atf.Acc{
		PreCheck:     testAccPreCheck,
		Providers:    testAccProviders,
		ResourceName: "hpegl_vmaas_instance_clone",
	}
	acc.RunResourcePlanTest(t)
}

func TestAccResourceInstanceCloneCreate(t *testing.T) {
	acc := &atf.Acc{
		ResourceName: "hpegl_vmaas_instance_clone",
		PreCheck:     testAccPreCheck,
		Providers:    testAccProviders,
		GetAPI: func(attr map[string]string) (interface{}, error) {
			cl, cfg := getAPIClient()
			iClient := api_client.InstancesAPIService{
				Client: cl,
				Cfg:    cfg,
			}
			id := toInt(attr["id"])

			return iClient.GetASpecificInstance(getAccContext(), id)
		},
	}

	acc.RunResourceTests(t)
}
