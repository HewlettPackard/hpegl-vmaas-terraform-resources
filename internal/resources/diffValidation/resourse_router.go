//  (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package diffvalidation

import (
	"fmt"

	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	haModePath          = "tier0_config.0.ha_mode"
	interSRiBGPPath     = "tier0_config.0.bgp.0.inter_sr_ibgp"
	failOverPath        = "tier0_config.0.fail_over"
	bgpEnabledPath      = "tier0_config.0.bgp.0.enable_bgp"
	bgpRestartTimerPath = "tier0_config.0.bgp.0.restart_time"
	bgpStaleTimerPath   = "tier0_config.0.bgp.0.stale_route_time"
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

func (r *Router) validateHAModeIsACTIVE_STANDBY() error {
	// BGP inter SR routing only applicable for Tier0 in active-active HA-mode

	if r.diff.HasChange(haModePath) || r.diff.HasChange(interSRiBGPPath) {
		action := r.diff.Get(haModePath)
		if action == "ACTIVE_STANDBY" {
			if (r.diff.Get(interSRiBGPPath)).(bool) {
				return fmt.Errorf("BGP inter SR routing only applicable for Tier0 in active-active HA-mode")
			}
			if r.diff.HasChange(haModePath) && r.diff.Get(failOverPath) == "" {
				return fmt.Errorf("failover mode is required when HA mode is set to ACTIVE_STANDBY")
			}
		}
	}

	return nil
}

func (r *Router) validateBGPTimers() error {
	// BGP graceful restart timers cannot be updated when BGP config is enabled
	if r.diff.HasChange(bgpRestartTimerPath) || r.diff.HasChange(bgpStaleTimerPath) {
		action := r.diff.Get(bgpEnabledPath)
		if action.(bool) {
			currRestartTimer, newRestartTimer := r.diff.GetChange(bgpRestartTimerPath)
			currStaleTimer, newStaleTimer := r.diff.GetChange(bgpStaleTimerPath)

			// During the creation of the Router
			validateTimerCreation := (r.diff.HasChange(bgpRestartTimerPath) && utils.IsEmpty(currRestartTimer) && newRestartTimer.(int) != DefaultRestartTimer) ||
				(r.diff.HasChange(bgpStaleTimerPath) && utils.IsEmpty(currStaleTimer) && newStaleTimer.(int) != DefaultStaleTimer)
			// While updating the Router
			validateTimerUpdation := (r.diff.HasChange(bgpRestartTimerPath) && !utils.IsEmpty(currRestartTimer)) || (r.diff.HasChange(bgpStaleTimerPath) && !utils.IsEmpty(currStaleTimer))

			if validateTimerCreation || validateTimerUpdation {
				return fmt.Errorf("BGP graceful restart timers cannot be updated when BGP config is enabled")
			}
		}
	}

	return nil
}

func (r *Router) validateTier0Config() error {
	err := r.validateHAModeIsACTIVE_STANDBY()
	if err != nil {
		return err
	}

	err = r.validateBGPTimers()
	if err != nil {
		return err
	}

	return nil
}
