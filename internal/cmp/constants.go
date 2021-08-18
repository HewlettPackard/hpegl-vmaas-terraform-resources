// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package cmp

import "time"

const (
	vmware           = "vmware"
	errExactMatch    = "error, could not find the %s with the specified name. Please verify the name and try again"
	provisionTypeKey = "provisionType"
	codeKey          = "code"
	nameKey          = "name"
	maxKey           = "max"
	// retry related constants
	maxTimeout = time.Hour * 2
)
