//  (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package diffvalidation

import (
	"fmt"

	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Router struct {
	diff *schema.ResourceDiff
}

func NewRouterValidate(diff *schema.ResourceDiff) *Router {
	return &Router{
		diff: diff,
	}
}

func (r *Router) DiffValidate() error {

	err := r.validateTier0Config()
	if err != nil {
		return err
	}

	return nil
}

func (r *Router) validateTier0Config() error {

	// BGP inter SR routing only applicable for Tier0 in active-active HA-mode
	haModePath := "tier0_config.0.ha_mode"
	interSRiBGPPath := "tier0_config.0.bgp.0.inter_sr_ibgp"
	if r.diff.HasChange(haModePath) || r.diff.HasChange(interSRiBGPPath) {
		action := r.diff.Get(haModePath)
		if action == "ACTIVE_STANDBY" {
			if (r.diff.Get(interSRiBGPPath)).(bool) {
				return fmt.Errorf("BGP inter SR routing only applicable for Tier0 in active-active HA-mode")
			}
		}
	}
	if r.diff.HasChange(haModePath) {
		action := r.diff.Get(haModePath)
		if action == "ACTIVE_STANDBY" {
			if r.diff.Get("tier0_config.0.fail_over") == "" {
				return fmt.Errorf("failover mode is required when HA mode is set to ACTIVE_STANDBY")
			}
		}
	}

	//BGP graceful restart timers cannot be updated when BGP config is enabled
	bgpEnabledPath := "tier0_config.0.bgp.0.enable_bgp"
	if r.diff.HasChange(bgpEnabledPath) || (r.diff.HasChange("tier0_config.0.bgp.0.restart_time") || r.diff.HasChange("tier0_config.0.bgp.0.stale_route_time")) {
		action := r.diff.Get(bgpEnabledPath)
		if action.(bool) {
			oldRestartTimer, newRestartTimer := r.diff.GetChange("tier0_config.0.bgp.0.restart_time")
			oldStaleTimer, newStaleTimer := r.diff.GetChange("tier0_config.0.bgp.0.stale_route_time")
			// During the creation of the Router
			if (r.diff.HasChange("tier0_config.0.bgp.0.restart_time") && utils.IsEmpty(oldRestartTimer) && newRestartTimer.(int) != utils.DefaultRestartTimer) || (r.diff.HasChange("tier0_config.0.bgp.0.stale_route_time") && utils.IsEmpty(oldStaleTimer) && newStaleTimer.(int) != utils.DefaultStaleTimer) {
				return fmt.Errorf("BGP graceful restart timers cannot be updated when BGP config is enabled")
				// While updating the Router
			} else if (r.diff.HasChange("tier0_config.0.bgp.0.restart_time") && !utils.IsEmpty(oldRestartTimer)) || (r.diff.HasChange("tier0_config.0.bgp.0.stale_route_time") && !utils.IsEmpty(oldStaleTimer)) {
				return fmt.Errorf("BGP graceful restart timers cannot be updated when BGP config is enabled")
				// While updating the Router
			}
		}
	}

	return nil
}
