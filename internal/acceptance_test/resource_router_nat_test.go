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

func TestVmaasRouterNatPlan(t *testing.T) {
	pkgutils.SkipAcc(t, "vmaas.resource.router_nat")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:             testAccResourceRouterNat(),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccResourceRouterNatCreate(t *testing.T) {
	pkgutils.SkipAcc(t, "vmaas.resource.router_nat")
	if testing.Short() {
		t.Skip("Skipping router resource creation in short mode")
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			checkResourceDestroy(
				"hpegl_vmaas_router_nat_rule.tf_nat",
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
					return iClient.GetSpecificRouterNat(context.Background(), routerID, id)
				},
			),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceRouterNat(),
				Check: resource.ComposeTestCheckFunc(
					validateResource(
						"hpegl_vmaas_router_nat_rule.tf_nat",
					),
				),
			},
		},
	})
}

func testAccResourceRouterNat() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	return providerStanza + fmt.Sprintf(`
	resource "hpegl_vmaas_router_nat_rule" "tf_nat" {
		name        = "%s_%d"
		router_id   = %d
		enabled     = %t
		description = "%s"
		config {
		  action   = "%s"
		  logging  = %t
		  firewall = "%s"
		}
		source_network      = "%s"
		translated_network  = "%s"
		destination_network = "%s"
		translated_ports    = %d
		priority            = %d
	  }
`,
		viper.GetString("vmaas.resource.router.tier0.name"),
		r.Int63n(999999),
		viper.GetInt("vmaas.resource.router_nat.router_id"),
		viper.GetBool("vmaas.resource.router_nat.enabled"),
		viper.GetString("vmaas.resource.router_nat.description"),
		viper.GetString("vmaas.resource.router_nat.config.action"),
		viper.GetBool("vmaas.resource.router_nat.config.logging"),
		viper.GetString("vmaas.resource.router_nat.config.firewall"),
		viper.GetString("vmaas.resource.router_nat.source_network"),
		viper.GetString("vmaas.resource.router_nat.translated_network"),
		viper.GetString("vmaas.resource.router_nat.destination_network"),
		viper.GetInt("vmaas.resource.router_nat.translated_ports"),
		viper.GetInt("vmaas.resource.router_nat.priority"),
	)
}
