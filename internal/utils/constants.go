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
	// data constants
	ErrInvalidType   = "invalid Type"
	ErrKeyNotDefined = "key is not defined"
	ErrSet           = "failed to set"
	NAN              = 0
	// retry constants
	defaultRetryTimeout = time.Second * 5
	defaultRetryCount   = 3
	defaultTimeout      = time.Duration(0)
	// power constants
	PowerOn         = "poweron"
	PowerOff        = "poweroff"
	Restart         = "restart"
	Suspend         = "suspend"
	Deleting        = "deleting"
	Deleted         = "deleted"
	Failed          = "failed"
	StateRestarting = "restarting"
)
