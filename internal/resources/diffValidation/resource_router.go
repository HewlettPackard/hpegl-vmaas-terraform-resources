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

func (r *Router) validateHAModeIsActiveStandby() error {
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
	// During tier0 router creation with bgp enabled, the timers should have default values
	if r.diff.HasChange(bgpRestartTimerPath) || r.diff.HasChange(bgpStaleTimerPath) {
		isBgpEnabled := r.diff.Get(bgpEnabledPath)
		if isBgpEnabled.(bool) {
			var checkRestartTimerDuringCreation, checkRestartTimerDuringUpdation bool
			var checkStaleTimerDuringCreation, checkStaleTimerDuringUpdation bool

			currRestartTimer, newRestartTimer := r.diff.GetChange(bgpRestartTimerPath)
			currStaleTimer, newStaleTimer := r.diff.GetChange(bgpStaleTimerPath)

			// Restart timer
			if r.diff.HasChange(bgpRestartTimerPath) {
				if utils.IsEmpty(currRestartTimer) {
					// Check: During tier0 router creation with bgp enabled,
					// the restart timer should not be different from DefaultRestartTimer
					checkRestartTimerDuringCreation = (newRestartTimer.(int) != DefaultRestartTimer)
				} else {
					// With bgp enabled, restart timer should not be updated
					checkRestartTimerDuringUpdation = true
				}
			}
			// Stale timer
			if r.diff.HasChange(bgpStaleTimerPath) {
				if utils.IsEmpty(currStaleTimer) {
					// Check: During tier0 router creation with bgp enabled,
					// the stale timer should not be different from DefaultStaleTimer
					checkStaleTimerDuringCreation = newStaleTimer.(int) != DefaultStaleTimer
				} else {
					// With bgp enabled, stale timer should not be updated
					checkStaleTimerDuringUpdation = true
				}
			}

			if checkRestartTimerDuringCreation || checkRestartTimerDuringUpdation ||
				checkStaleTimerDuringCreation || checkStaleTimerDuringUpdation {
				return fmt.Errorf("BGP graceful restart timers cannot be updated when BGP config is enabled")
			}
		}
	}

	return nil
}

func (r *Router) validateTier0Config() error {
	err := r.validateHAModeIsActiveStandby()
	if err != nil {
		return err
	}

	err = r.validateBGPTimers()
	if err != nil {
		return err
	}

	return nil
}
