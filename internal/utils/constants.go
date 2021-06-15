// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package utils

import "time"

const (
	// state constants
	StateUnknown      = "unknown"
	StateRunning      = "running"
	StateFailed       = "failed"
	StateProvisioning = "provisioning"
	// data constants
	ErrInvalidType   = "invalid Type"
	ErrKeyNotDefined = "key is not defined"
	ErrSet           = "failed to set"
	NAN              = 0
	// retry constants
	defaultTimeout    = time.Second * 5
	defaultRetryCount = 3
	PowerOn           = "poweron"
	PowerOff          = "poweroff"
	Restart           = "restart"
	Suspend           = "suspend"
)
