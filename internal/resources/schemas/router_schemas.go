package schemas

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func RouterTier0ConfigSchema() *schema.Schema {
	return &schema.Schema{
		Type:          schema.TypeSet,
		Optional:      true,
		Description:   "Tier0 Gateway configuration",
		MaxItems:      1,
		ExactlyOneOf:  []string{"tier1_config", "tier0_config"},
		ConflictsWith: []string{"tier1_config"},
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"ha_mode": {
					Type:     schema.TypeString,
					Required: true,
					ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
						"ACTIVE_ACTIVE",
					}, false)),
					Description: "Available values are 'ACTIVE_ACTIVE'",
				},
				"fail_over": {
					Type:     schema.TypeString,
					Required: true,
					ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
						"NON_PREEMPTIVE", "PREEMPTIVE",
					}, false)),
					Description: "Available values are 'NON_PREEMPTIVE'",
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
							"tier0_static": {
								Type:     schema.TypeBool,
								Optional: true,
							},
							"tier0_nat": {
								Type:     schema.TypeBool,
								Optional: true,
							},
							"tier0_ipsec_local_ip": {
								Type:     schema.TypeBool,
								Optional: true,
							},
							"tier0_dns_forwarder_ip": {
								Type:     schema.TypeBool,
								Optional: true,
							},
							"tier0_service_interface": {
								Type:     schema.TypeBool,
								Optional: true,
							},
							"tier0_external_interface": {
								Type:     schema.TypeBool,
								Optional: true,
							},
							"tier0_loopback_interface": {
								Type:     schema.TypeBool,
								Optional: true,
							},
							"tier0_segment": {
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
							"tier1_dns_forwarder_ip": {
								Type:     schema.TypeBool,
								Optional: true,
							},
							"tier1_static": {
								Type:     schema.TypeBool,
								Optional: true,
							},
							"tier1_lb_vip": {
								Type:     schema.TypeBool,
								Optional: true,
							},
							"tier1_nat": {
								Type:     schema.TypeBool,
								Optional: true,
							},
							"tier1_lb_snat": {
								Type:     schema.TypeBool,
								Optional: true,
							},
							"tier1_ipsec_local_endpoint": {
								Type:     schema.TypeBool,
								Optional: true,
							},
							"tier1_service_interface": {
								Type:     schema.TypeBool,
								Optional: true,
							},
							"tier1_segment": {
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
							"local_as_num": {
								Type:     schema.TypeInt,
								Optional: true,
							},
							"ecmp": {
								Type:     schema.TypeBool,
								Optional: true,
							},
							"multipath_relax": {
								Type:     schema.TypeBool,
								Optional: true,
							},
							"inter_sr_ibgp": {
								Type:     schema.TypeBool,
								Optional: true,
							},
							"restart_mode": {
								Type:     schema.TypeString,
								Optional: true,
								ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
									"HELPER_ONLY",
									"GRACEFUL_RESTART_AND_HELPER",
									"DISABLE",
								}, false)),
							},
							"restart_time": {
								Type:     schema.TypeInt,
								Optional: true,
							},
							"stale_route_time": {
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
							"tier1_connected": {
								Type:     schema.TypeBool,
								Optional: true,
							},
							"tier1_nat": {
								Type:     schema.TypeBool,
								Optional: true,
							},
							"tier1_static_routes": {
								Type:     schema.TypeBool,
								Optional: true,
							},
							"tier1_lb_vip": {
								Type:     schema.TypeBool,
								Optional: true,
							},
							"tier1_lb_snat": {
								Type:     schema.TypeBool,
								Optional: true,
							},
							"tier1_dns_forwarder_ip": {
								Type:     schema.TypeBool,
								Optional: true,
							},
							"tier1_ipsec_local_endpoint": {
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
