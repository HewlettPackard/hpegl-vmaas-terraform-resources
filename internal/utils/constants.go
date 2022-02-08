// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package utils

import "time"

const (
	// state constants
	StateUnknown      = "unknown"
	StateRunning      = "running"
	StateFailed       = "failed"
	StateProvisioning = "provisioning"
	StateStopped      = "stopped"
	StateStopping     = "stopping"
	StateSuspended    = "suspended"
	StateSuspending   = "suspending"
	StateResizing     = "resizing"
	StateRestarting   = "restarting"
	// data constants
	ErrInvalidType   = "invalid Type"
	ErrKeyNotDefined = "key is not defined"
	ErrSet           = "failed to set"
	NAN              = 0
	// retry constants
	defaultRetryDelay = time.Second * 5
	defaultRetryCount = 3
	defaultTimeout    = time.Hour * 24
	noRetryCount      = -1
	// power constants
	PowerOn  = "poweron"
	PowerOff = "poweroff"
	Restart  = "restart"
	Suspend  = "suspend"
	Deleting = "deleting"
	Deleted  = "deleted"
	Failed   = "failed"
	// Router Constants
	BgpEnabledPath      = "tier0_config.0.bgp.0.enable_bgp"
	BgpRestartTimerPath = "tier0_config.0.bgp.0.restart_time"
	BgpStaleTimerPath   = "tier0_config.0.bgp.0.stale_route_time"
	DefaultRestartTimer = 180
	DefaultStaleTimer   = 600
)
