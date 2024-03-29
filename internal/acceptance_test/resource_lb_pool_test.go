// (C) Copyright 2022 Hewlett Packard Enterprise Development LP

package acceptancetest

import (
	"context"
	"testing"

	api_client "github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/client"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/atf"
)

func TestVmaasLBPoolPlan(t *testing.T) {
	acc := &atf.Acc{
		PreCheck:     testAccPreCheck,
		Providers:    testAccProviders,
		ResourceName: "hpegl_vmaas_load_balancer_pool",
	}
	acc.RunResourcePlanTest(t)
}

func TestAccResourceLBPoolCreate(t *testing.T) {
	acc := &atf.Acc{
		ResourceName: "hpegl_vmaas_load_balancer_pool",
		PreCheck:     testAccPreCheck,
		Providers:    testAccProviders,
		GetAPI: func(attr map[string]string) (interface{}, error) {
			cl, cfg := getAPIClient()
			iClient := api_client.LoadBalancerAPIService{
				Client: cl,
				Cfg:    cfg,
			}
			id := toInt(attr["id"])
			lbID := toInt(attr["lb_id"])

			return iClient.GetSpecificLBPool(context.Background(), lbID, id)
		},
	}

	acc.RunResourceTests(t)
}

func TestAccResourceLBPoolCreate_poolErr(t *testing.T) {
	acc := &atf.Acc{
		ResourceName: "hpegl_vmaas_load_balancer_pool",
		PreCheck:     testAccPreCheck,
		Providers:    testAccProviders,
		Version:      "pool_err",
		GetAPI: func(attr map[string]string) (interface{}, error) {
			cl, cfg := getAPIClient()
			iClient := api_client.LoadBalancerAPIService{
				Client: cl,
				Cfg:    cfg,
			}
			id := toInt(attr["id"])
			lbID := toInt(attr["lb_id"])

			return iClient.GetSpecificLBPool(context.Background(), lbID, id)
		},
	}

	acc.RunResourceTests(t)
}
