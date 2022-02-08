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
	if r.diff.HasChange(utils.BgpEnabledPath) || (r.diff.HasChange(utils.BgpRestartTimerPath) || r.diff.HasChange(utils.BgpStaleTimerPath)) {
		action := r.diff.Get(utils.BgpEnabledPath)
		if action.(bool) {
			currRestartTimer, newRestartTimer := r.diff.GetChange(utils.BgpRestartTimerPath)
			currStaleTimer, newStaleTimer := r.diff.GetChange(utils.BgpStaleTimerPath)
			// During the creation of the Router
			if (r.diff.HasChange(utils.BgpRestartTimerPath) && utils.IsEmpty(currRestartTimer) && newRestartTimer.(int) != utils.DefaultRestartTimer) || (r.diff.HasChange(utils.BgpStaleTimerPath) && utils.IsEmpty(currStaleTimer) && newStaleTimer.(int) != utils.DefaultStaleTimer) {
				return fmt.Errorf("BGP graceful restart timers cannot be updated when BGP config is enabled")
				// While updating the Router
			} else if (r.diff.HasChange(utils.BgpRestartTimerPath) && !utils.IsEmpty(currRestartTimer)) || (r.diff.HasChange(utils.BgpStaleTimerPath) && !utils.IsEmpty(currStaleTimer)) {
				return fmt.Errorf("BGP graceful restart timers cannot be updated when BGP config is enabled")
				// While updating the Router
			}
		}
	}

	return nil
}
