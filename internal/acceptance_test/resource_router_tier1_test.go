// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package acceptancetest

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	api_client "github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/client"
	pkgutils "github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/spf13/viper"
)

func TestVmaasRouterTier1Plan(t *testing.T) {
	pkgutils.SkipAcc(t, "vmaas.resource.router")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:             testAccResourceTier1Router(),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccResourceTier1RouterCreate(t *testing.T) {
	pkgutils.SkipAcc(t, "vmaas.resource.router")
	if testing.Short() {
		t.Skip("Skipping router resource creation in short mode")
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			checkResourceDestroy(
				"hpegl_vmaas_router.tf_tier1",
				func(cl *api_client.APIClient, cfg api_client.Configuration, id int, attr map[string]string,
				) (interface{}, error) {
					iClient := api_client.RouterAPIService{
						Client: cl,
						Cfg:    cfg,
					}
					return iClient.GetSpecificRouter(context.Background(), id)
				},
			),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceTier1Router(),
				Check: resource.ComposeTestCheckFunc(
					validateResource(
						"hpegl_vmaas_router.tf_tier1",
					),
				),
			},
		},
	})
}

func testAccResourceTier1Router() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	return providerStanza + fmt.Sprintf(`
	resource "hpegl_vmaas_router" "tf_tier1" {
		name     = "%s_%d"
		enable   = %t
		group_id = "%s"
		tier1_config {
			fail_over = "%s"
			route_advertisement {
				tier1_connected = %t
				tier1_static_routes = %t
				tier1_dns_forwarder_ip = %t
				tier1_lb_vip = %t
				tier1_nat = %t
				tier1_lb_snat = %t
				tier1_ipsec_local_endpoint = %t
			}
		}
	}
`,
		viper.GetString("vmaas.resource.router.tier1.name"),
		r.Int63n(999999),
		viper.GetBool("vmaas.resource.router.tier1.enable"),
		viper.GetString("vmaas.resource.router.tier1.group_id"),
		viper.GetString("vmaas.resource.router.tier1.tier1_config.fail_over"),
		viper.GetBool("vmaas.resource.router.tier1.tier1_config.route_advertisement.tier1_connected"),
		viper.GetBool("vmaas.resource.router.tier1.tier1_config.route_advertisement.tier1_static_routes"),
		viper.GetBool("vmaas.resource.router.tier1.tier1_config.route_advertisement.tier1_dns_forwarder_ip"),
		viper.GetBool("vmaas.resource.router.tier1.tier1_config.route_advertisement.tier1_lb_vip"),
		viper.GetBool("vmaas.resource.router.tier1.tier1_config.route_advertisement.tier1_nat"),
		viper.GetBool("vmaas.resource.router.tier1.tier1_config.route_advertisement.tier1_lb_snat"),
		viper.GetBool("vmaas.resource.router.tier1.tier1_config.route_advertisement.tier1_ipsec_local_endpoint"))
}
