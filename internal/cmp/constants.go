// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package cmp

import "time"

const (
	vmware        = "vmware"
	nsxtSegment   = "NSX-T Segment"
	errExactMatch = "error, could not find the %s with the specified name. Please verify the name and try again"
	successErr    = "got success = 'false while %s"
	// query params keys
	provisionTypeKey = "provisionType"
	codeKey          = "code"
	nameKey          = "name"
	maxKey           = "max"
	externalNameKey  = "externalName"
	// retry related constants
	maxTimeout = time.Hour * 2
	// router consts
	tier0GatewayType = "NSX-T Tier-0 Gateway"
	tier1GatewayType = "NSX-T Tier-1 Gateway"
)
