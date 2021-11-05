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

func TestVmaasNetworkPlan(t *testing.T) {
	pkgutils.SkipAcc(t, "vmaas.resource.network")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:             testAccResourceNetwork(),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccResourceNetworkCreate(t *testing.T) {
	pkgutils.SkipAcc(t, "vmaas.resource.network")
	if testing.Short() {
		t.Skip("Skipping network resource creation in short mode")
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			checkResourceDestroy(
				"hpegl_vmaas_network.tf_network",
				func(cl *api_client.APIClient, cfg api_client.Configuration, id int, attr map[string]string,
				) (interface{}, error) {
					iClient := api_client.NetworksAPIService{
						Client: cl,
						Cfg:    cfg,
					}

					return iClient.GetSpecificNetwork(context.Background(), id)
				},
			),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNetwork(),
				Check: resource.ComposeTestCheckFunc(
					validateResource(
						"hpegl_vmaas_network.tf_network",
					),
				),
			},
		},
	})
}

func testAccResourceNetwork() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	return providerStanza + fmt.Sprintf(`
	resource "hpegl_vmaas_network" "tf_network" {
			name                = "%s_%d"
			group_id   			= "%s"
			active   			= %t
			dhcp_server   		= %t
			description   		= "%s"
			cidr   			    = "%s"
			gateway   			= "%s"
			pool_id 			= %d
	}
`,
		viper.GetString("vmaas.resource.network.name"),
		r.Int63n(999999),
		viper.GetString("vmaas.resource.network.group_id"),
		viper.GetBool("vmaas.resource.network.active"),
		viper.GetBool("vmaas.resource.network.dhcp_server"),
		viper.GetString("vmaas.resource.network.description"),
		viper.GetString("vmaas.resource.network.cidr"),
		viper.GetString("vmaas.resource.network.gateway"),
		viper.GetInt("vmaas.resource.network.pool_id"))
}
