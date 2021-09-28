package schemas

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func RouterTier0ConfigSchema() *schema.Schema {
	return &schema.Schema{
		Type:          schema.TypeSet,
		Optional:      true,
		Description:   "Tier0 Gateway configuration",
		MaxItems:      1,
		ConflictsWith: []string{"tier1_config"},
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"ha_mode": {
					Type:         schema.TypeString,
					Optional:     true,
					ExactlyOneOf: []string{"ACTIVE_ACTIVE"},
					Description:  "Available values are 'ACTIVE_ACTIVE'",
				},
				"fail_over": {
					Type:         schema.TypeString,
					Required:     true,
					ExactlyOneOf: []string{"NON_PREEMPTIVE"},
					Description:  "Available values are 'NON_PREEMPTIVE'",
				},
				"enable_bgp": {
					Type:     schema.TypeBool,
					Optional: true,
				},
				"route_redistribution_tier0": {
					Type:     schema.TypeSet,
					MaxItems: 1,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"TIER0_STATIC": {
								Type:     schema.TypeBool,
								Optional: true,
							},
							"TIER0_NAT": {
								Type:     schema.TypeBool,
								Optional: true,
							},
							"TIER0_IPSEC_LOCAL_IP": {
								Type:     schema.TypeBool,
								Optional: true,
							},
							"TIER0_DNS_FORWARDER_IP": {
								Type:     schema.TypeBool,
								Optional: true,
							},
							"TIER0_SERVICE_INTERFACE": {
								Type:     schema.TypeBool,
								Optional: true,
							},
							"TIER0_EXTERNAL_INTERFACE": {
								Type:     schema.TypeBool,
								Optional: true,
							},
							"TIER0_LOOPBACK_INTERFACE": {
								Type:     schema.TypeBool,
								Optional: true,
							},
							"TIER0_SEGMENT": {
								Type:     schema.TypeBool,
								Optional: true,
							},
						},
					},
				},
				"route_redistribution_tier1": {
					Type:     schema.TypeSet,
					MaxItems: 1,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"TIER1_DNS_FORWARDER_IP": {
								Type:     schema.TypeBool,
								Optional: true,
							},
							"TIER1_STATIC": {
								Type:     schema.TypeBool,
								Optional: true,
							},
							"TIER1_LB_VIP": {
								Type:     schema.TypeBool,
								Optional: true,
							},
							"TIER1_NAT": {
								Type:     schema.TypeBool,
								Optional: true,
							},
							"TIER1_LB_SNAT": {
								Type:     schema.TypeBool,
								Optional: true,
							},
							"TIER1_IPSEC_LOCAL_ENDPOINT": {
								Type:     schema.TypeBool,
								Optional: true,
							},
							"TIER1_SERVICE_INTERFACE": {
								Type:     schema.TypeBool,
								Optional: true,
							},
							"TIER1_SEGMENT": {
								Type:     schema.TypeBool,
								Optional: true,
							},
						},
					},
				},
				"bgp": {
					Type:     schema.TypeSet,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"LOCAL_AS_NUM": {
								Type:     schema.TypeInt,
								Optional: true,
							},
							"ECMP": {
								Type:     schema.TypeBool,
								Optional: true,
							},
							"MULTIPATH_RELAX": {
								Type:     schema.TypeBool,
								Optional: true,
							},
							"INTER_SR_IBGP": {
								Type:     schema.TypeBool,
								Optional: true,
							},
							"RESTART_MODE": {
								Type:     schema.TypeString,
								Optional: true,
								ExactlyOneOf: []string{
									"HELPER_ONLY",
									"GRACEFUL_RESTART_AND_HELPER",
									"DISABLE",
								},
							},
							"RESTART_TIME": {
								Type:     schema.TypeInt,
								Optional: true,
							},
							"STALE_ROUTE_TIME": {
								Type:     schema.TypeInt,
								Optional: true,
							},
						},
					},
				},
			},
		},
	}
}

func RouterTier1ConfigSchema() *schema.Schema {
	return &schema.Schema{
		Type:          schema.TypeSet,
		Optional:      true,
		Description:   "Tier1 Gateway configuration",
		ConflictsWith: []string{"tier0_config"},
		MaxItems:      1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"tier0_gateway": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"edge_cluster": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"route_advertisement": {
					Type:     schema.TypeSet,
					Optional: true,
					MaxItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"TIER1_CONNECTED": {
								Type:     schema.TypeBool,
								Optional: true,
							},
							"TIER1_NAT": {
								Type:     schema.TypeBool,
								Optional: true,
							},
							"TIER1_STATIC_ROUTES": {
								Type:     schema.TypeBool,
								Optional: true,
							},
							"TIER1_LB_VIP": {
								Type:     schema.TypeBool,
								Optional: true,
							},
							"TIER1_LB_SNAT": {
								Type:     schema.TypeBool,
								Optional: true,
							},
							"TIER1_DNS_FORWARDER_IP": {
								Type:     schema.TypeBool,
								Optional: true,
							},
							"TIER1_IPSEC_LOCAL_ENDPOINT": {
								Type:     schema.TypeBool,
								Optional: true,
							},
						},
					},
				},
			},
		},
	}
}
