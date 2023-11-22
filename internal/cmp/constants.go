// (C) Copyright 2021-2022 Hewlett Packard Enterprise Development LP

package cmp

import "time"

const (
	vmware        = "vmware"
	nsxtSegment   = "NSX Segment"
	nsxt          = "NSX"
	errExactMatch = "error, could not find the %s with the specified name. Please verify the name and try again"
	successErr    = "got success = 'false while %s"
	// query params keys
	provisionTypeKey = "provisionType"
	codeKey          = "code"
	nameKey          = "name"
	maxKey           = "max"
	externalNameKey  = "externalName"
	filterTypeKey    = "filterType"
	// retry related constants
	maxTimeout = time.Hour * 2
	// router consts
	tier0GatewayType             = "NSX Tier-0 Gateway"
	tier1GatewayType             = "NSX Tier-1 Gateway"
	routerFirewallExternalPolicy = "GatewayPolicy"
	syncedTypeValue              = "Synced"

	// load balancer consts
	TCP      = "tcp"
	UDP      = "udp"
	HTTP     = "http"
	COOKIE   = "COOKIE"
	SOURCEIP = "SOURCE_IP"
)
