// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package acceptancetest

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	pkgutils "github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/spf13/viper"
)

func TestVmaasRouterTier0Plan(t *testing.T) {
	pkgutils.SkipAcc(t, "vmaas.resource.router")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:             testAccResourceTier0Router(),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccResourceTier0RouterCreate(t *testing.T) {
	pkgutils.SkipAcc(t, "vmaas.resource.router")
	if testing.Short() {
		t.Skip("Skipping router resource creation in short mode")
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(checkResourceDestroy("hpegl_vmaas_router.tf_tier0")),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceTier0Router(),
				Check: resource.ComposeTestCheckFunc(
					validateResource(
						"hpegl_vmaas_router.tf_tier0",
					),
				),
			},
		},
	})
}

func testAccResourceTier0Router() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	return providerStanza + fmt.Sprintf(`
	resource "hpegl_vmaas_router" "tf_tier0" {
		name     = "%s_%d"
		enable   = %t
		group_id = "%s"
		tier0_config {
			bgp {
			ecmp             = %t
			enable_bgp       = %t
			inter_sr_ibgp    = %t
			local_as_num     = %d
			multipath_relax  = %t
			restart_mode     = "%s"
			restart_time     = %d
			stale_route_time = %d
			}
			route_redistribution_tier0 {
			tier0_dns_forwarder_ip   = %t
			tier0_external_interface = %t
			tier0_ipsec_local_ip     = %t
			tier0_loopback_interface = %t
			tier0_nat                = %t
			tier0_segment            = %t
			tier0_service_interface  = %t
			tier0_static             = %t
			}
			route_redistribution_tier1 {
			tier1_dns_forwarder_ip     = %t
			tier1_service_interface    = %t
			tier1_ipsec_local_endpoint = %t
			tier1_lb_snat              = %t
			tier1_lb_vip               = %t
			tier1_nat                  = %t
			tier1_segment              = %t
			tier1_static               = %t
			}
			fail_over = "%s"
			ha_mode   = "%s"
		}
	}
`,
		viper.GetString("vmaas.resource.router.tier0.name"),
		r.Int63n(999999),
		viper.GetBool("vmaas.resource.router.tier0.enable"),
		viper.GetString("vmaas.resource.router.tier0.group_id"),
		viper.GetBool("vmaas.resource.router.tier0.tier0_config.bgp.ecmp"),
		viper.GetBool("vmaas.resource.router.tier0.tier0_config.bgp.enable_bgp"),
		viper.GetBool("vmaas.resource.router.tier0.tier0_config.bgp.inter_sr_ibgp"),
		viper.GetInt("vmaas.resource.router.tier0.tier0_config.bgp.local_as_num"),
		viper.GetBool("vmaas.resource.router.tier0.tier0_config.bgp.multipath_relax"),
		viper.GetString("vmaas.resource.router.tier0.tier0_config.bgp.restart_mode"),
		viper.GetInt("vmaas.resource.router.tier0.tier0_config.bgp.restart_time"),
		viper.GetInt("vmaas.resource.router.tier0.tier0_config.bgp.stale_route_time"),
		viper.GetBool("vmaas.resource.router.tier0.tier0_config.route_redistribution_tier0.tier0_dns_forwarder_ip"),
		viper.GetBool("vmaas.resource.router.tier0.tier0_config.route_redistribution_tier0.tier0_external_interface"),
		viper.GetBool("vmaas.resource.router.tier0.tier0_config.route_redistribution_tier0.tier0_ipsec_local_ip"),
		viper.GetBool("vmaas.resource.router.tier0.tier0_config.route_redistribution_tier0.tier0_loopback_interface"),
		viper.GetBool("vmaas.resource.router.tier0.tier0_config.route_redistribution_tier0.tier0_nat"),
		viper.GetBool("vmaas.resource.router.tier0.tier0_config.route_redistribution_tier0.tier0_segment"),
		viper.GetBool("vmaas.resource.router.tier0.tier0_config.route_redistribution_tier0.tier0_service_interface"),
		viper.GetBool("vmaas.resource.router.tier0.tier0_config.route_redistribution_tier0.tier0_static"),
		viper.GetBool("vmaas.resource.router.tier0.tier0_config.route_redistribution_tier1.tier1_dns_forwarder_ip"),
		viper.GetBool("vmaas.resource.router.tier0.tier0_config.route_redistribution_tier1.tier1_service_interface"),
		viper.GetBool("vmaas.resource.router.tier0.tier0_config.route_redistribution_tier1.tier1_ipsec_local_endpoint"),
		viper.GetBool("vmaas.resource.router.tier0.tier0_config.route_redistribution_tier1.tier1_lb_snat"),
		viper.GetBool("vmaas.resource.router.tier0.tier0_config.route_redistribution_tier1.tier1_lb_vip"),
		viper.GetBool("vmaas.resource.router.tier0.tier0_config.route_redistribution_tier1.tier1_nat"),
		viper.GetBool("vmaas.resource.router.tier0.tier0_config.route_redistribution_tier1.tier1_segment"),
		viper.GetBool("vmaas.resource.router.tier0.tier0_config.route_redistribution_tier1.tier1_static"),
		viper.GetString("vmaas.resource.router.tier0.tier0_config.fail_over"),
		viper.GetString("vmaas.resource.router.tier0.tier0_config.ha_mode"))
}
