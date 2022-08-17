// (C) Copyright 2022 Hewlett Packard Enterprise Development LP

package acceptancetest

import (
	"context"
	"testing"

	api_client "github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/client"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/atf"
)

func TestVmaasLoadBalancerPlan(t *testing.T) {
	acc := &atf.Acc{
		PreCheck:     testAccPreCheck,
		Providers:    testAccProviders,
		ResourceName: "hpegl_vmaas_load_balancer",
	}
	acc.RunResourcePlanTest(t)
}

func TestAccResourceLoadBalancerCreate(t *testing.T) {
	acc := &atf.Acc{
		ResourceName: "hpegl_vmaas_load_balancer",
		PreCheck:     testAccPreCheck,
		Providers:    testAccProviders,
		GetAPI: func(attr map[string]string) (interface{}, error) {
			cl, cfg := getAPIClient()
			iClient := api_client.LoadBalancerAPIService{
				Client: cl,
				Cfg:    cfg,
			}
			id := toInt(attr["id"])

			return iClient.GetSpecificLoadBalancers(context.Background(), id)
		},
	}

	acc.RunResourceTests(t)
}
