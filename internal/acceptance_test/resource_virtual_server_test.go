// (C) Copyright 2022 Hewlett Packard Enterprise Development LP

package acceptancetest

import (
	"context"
	"testing"

	api_client "github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/client"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/atf"
)

func TestVmaasLBVirtualServerPlan(t *testing.T) {
	acc := &atf.Acc{
		PreCheck:     testAccPreCheck,
		Providers:    testAccProviders,
		ResourceName: "hpegl_vmaas_load_balancer_virtual_server",
	}
	acc.RunResourcePlanTest(t)
}

func TestAccResourceLBVirtualServerCreate(t *testing.T) {
	acc := &atf.Acc{
		ResourceName: "hpegl_vmaas_load_balancer_virtual_server",
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

			return iClient.GetSpecificLBVirtualServer(context.Background(), lbID, id)
		},
	}

	acc.RunResourceTests(t)
}

func TestAccResourceLBVirtualServerCreate_virtualserverErr(t *testing.T) {
	acc := &atf.Acc{
		ResourceName: "hpegl_vmaas_load_balancer_virtual_server",
		PreCheck:     testAccPreCheck,
		Providers:    testAccProviders,
		Version:      "virtualserver_err",
		GetAPI: func(attr map[string]string) (interface{}, error) {
			cl, cfg := getAPIClient()
			iClient := api_client.LoadBalancerAPIService{
				Client: cl,
				Cfg:    cfg,
			}
			id := toInt(attr["id"])
			lbID := toInt(attr["lb_id"])

			return iClient.GetSpecificLBVirtualServer(context.Background(), lbID, id)
		},
	}

	acc.RunResourceTests(t)
}
