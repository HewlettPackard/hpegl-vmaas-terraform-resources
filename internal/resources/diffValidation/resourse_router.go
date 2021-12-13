//  (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package diffvalidation

import (
	"fmt"

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
	if r.diff.HasChange(haModePath) {
		action := r.diff.Get(haModePath)
		if action == "ACTIVE_STANDBY" {
			if (r.diff.Get("tier0_config.0.bgp.0.inter_sr_ibgp")).(bool) {
				return fmt.Errorf("BGP inter SR routing only applicable for Tier0 in active-active HA-mode")
			}
			if r.diff.Get("tier0_config.0.fail_over") == "" {
				return fmt.Errorf("failover mode is required when HA mode is set to ACTIVE_STANDBY")
			}

		}
	}

	//BGP graceful restart timers cannot be updated when BGP config is enabled
	bgpEnabledPath := "tier0_config.0.bgp.0.enable_bgp"
	if r.diff.HasChange(bgpEnabledPath) {
		action := r.diff.Get(bgpEnabledPath)
		if action.(bool) {
			if r.diff.HasChange("tier0_config.0.bgp.0.restart_time") || r.diff.HasChange("tier0_config.0.bgp.0.stale_route_time") {
				return fmt.Errorf("BGP graceful restart timers cannot be updated when BGP config is enabled")
			}
		}
	}

	return nil
}