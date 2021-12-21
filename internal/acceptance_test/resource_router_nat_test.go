// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package acceptancetest

import (
	"fmt"
	"net/http"
	"testing"

	api_client "github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/client"
	pkgutils "github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/utils"
	"github.com/hewlettpackard/hpegl-provider-lib/pkg/atf"
)

func TestVmaasRouterNatPlan(t *testing.T) {
	acc := &atf.Acc{
		PreCheck:     testAccPreCheck,
		Providers:    testAccProviders,
		ResourceName: "hpegl_vmaas_router_nat_rule",
	}
	acc.RunResourcePlanTest(t)
}

func TestAccResourceRouterNatCreate(t *testing.T) {
	acc := &atf.Acc{
		ResourceName: "hpegl_vmaas_router_nat_rule",
		PreCheck:     testAccPreCheck,
		Providers:    testAccProviders,
		GetAPI:       getResourceRouterNatCreate,
		ValidateResourceDestroy: func(attr map[string]string) error {
			_, err := getResourceRouterNatCreate(attr)
			statusCode := pkgutils.GetStatusCode(err)
			if statusCode != http.StatusNotFound {
				return fmt.Errorf("expected %d statuscode, but got %d", 404, statusCode)
			}

			return nil
		},
	}

	acc.RunResourceTests(t)
}

func getResourceRouterNatCreate(attr map[string]string) (interface{}, error) {
	cl, cfg := getAPIClient()
	iClient := api_client.RouterAPIService{
		Client: cl,
		Cfg:    cfg,
	}
	id := toInt(attr["id"])
	routerID := toInt(attr["router_id"])

	return iClient.GetSpecificRouterNat(getAccContext(), routerID, id)
}
