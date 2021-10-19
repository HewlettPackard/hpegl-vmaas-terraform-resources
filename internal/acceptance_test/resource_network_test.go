// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package acceptancetest

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"testing"
	"time"

	api_client "github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/client"
	pkgutils "github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(testVmaasNetworkDestroy("hpegl_vmaas_network.tf_network")),
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

func testVmaasNetworkDestroy(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("resource %s not found", name)
		}
		id, err := strconv.Atoi(rs.Primary.Attributes["id"])
		if err != nil {
			return fmt.Errorf("error while converting id into int, %w", err)
		}

		apiClient, cfg := getAPIClient()
		iClient := api_client.NetworksAPIService{
			Client: apiClient,
			Cfg:    cfg,
		}
		_, err = iClient.GetSpecificNetwork(context.Background(), id)

		statusCode := pkgutils.GetStatusCode(err)
		if statusCode != http.StatusNotFound {
			return fmt.Errorf("Expected %d statuscode, but got %d", 404, statusCode)
		}

		return nil
	}
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
