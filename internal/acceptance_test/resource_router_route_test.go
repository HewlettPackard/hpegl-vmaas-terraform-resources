// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package acceptancetest

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"

	api_client "github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/client"
	pkgutils "github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/spf13/viper"
)

func TestVmaasRouterRoutePlan(t *testing.T) {
	pkgutils.SkipAcc(t, "vmaas.resource.router_route")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:             testAccResourceRouterRoute(),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccResourceRouterRouteCreate(t *testing.T) {
	pkgutils.SkipAcc(t, "vmaas.resource.router_route")
	if testing.Short() {
		t.Skip("Skipping router resource creation in short mode")
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			checkResourceDestroy(
				"hpegl_vmaas_router_route.tf_route",
				func(cl *api_client.APIClient, cfg api_client.Configuration, id int, attr map[string]string,
				) (interface{}, error) {
					iClient := api_client.RouterAPIService{
						Client: cl,
						Cfg:    cfg,
					}
					routerStr := attr["router_id"]
					routerID, err := strconv.Atoi(routerStr)
					if err != nil {
						return nil, err
					}

					return iClient.GetSpecificRouterRoute(context.Background(), routerID, id)
				},
			),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceRouterRoute(),
				Check: resource.ComposeTestCheckFunc(
					validateResource(
						"hpegl_vmaas_router_route.tf_route",
					),
				),
			},
		},
	})
}

func testAccResourceRouterRoute() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	return providerStanza + fmt.Sprintf(`
	resource "hpegl_vmaas_router_route" "tf_route" {
		name          = "%s_%d"
		router_id     = %d
		description   = "%s"
		enabled       = %t
		default_route = %t
		network       = "%s"
		next_hop      = "%s"
		mtu           = "%s"
		priority      = %d
	  }
`,
		viper.GetString("vmaas.resource.router_route.name"),
		r.Int63n(999999),
		viper.GetInt("vmaas.resource.router_route.router_id"),
		viper.GetString("vmaas.resource.router_route.description"),
		viper.GetBool("vmaas.resource.router_route.enabled"),
		viper.GetBool("vmaas.resource.router_route.default_route"),
		viper.GetString("vmaas.resource.router_route.network"),
		viper.GetString("vmaas.resource.router_route.next_hop"),
		viper.GetString("vmaas.resource.router_route.mtu"),
		viper.GetInt("vmaas.resource.router_route.priority"),
	)
}
