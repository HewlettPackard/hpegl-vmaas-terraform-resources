package schemas

import (
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/resources/validations"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func RouterTier0ConfigSchema() *schema.Schema {
	return &schema.Schema{
		Type:          schema.TypeList,
		Optional:      true,
		Description:   "Tier0 Gateway configuration",
		MaxItems:      1,
		ExactlyOneOf:  []string{"tier1_config", "tier0_config"},
		ConflictsWith: []string{"tier1_config"},
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"bgp": {
					Type:     schema.TypeList,
					Required: true,
					MaxItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"local_as_num": {
								Type:             schema.TypeInt,
								Required:         true,
								ValidateDiagFunc: validations.IntAtLeast(1),
								Description:      "Local AS Number",
							},
							"ecmp": {
								Type:        schema.TypeBool,
								Optional:    true,
								Default:     true,
								Description: "ECMP",
							},
							"multipath_relax": {
								Type:        schema.TypeBool,
								Optional:    true,
								Default:     true,
								Description: "Multipath Relax",
							},
							"inter_sr_ibgp": {
								Type:        schema.TypeBool,
								Optional:    true,
								Default:     true,
								Description: "Inter SR iBGP",
							},
							"restart_mode": {
								Type:         schema.TypeString,
								Required:     true,
								InputDefault: "HELPER_ONLY",
								ValidateDiagFunc: validations.StringInSlice([]string{
									"HELPER_ONLY",
									"GRACEFUL_RESTART_AND_HELPER",
									"DISABLE",
								}, false),
								Description: "Graceful Restart",
							},
							"restart_time": {
								Type:             schema.TypeInt,
								Required:         true,
								ValidateDiagFunc: validations.IntBetween(1, 3600),
								Description:      "Graceful Restart Timer",
							},
							"stale_route_time": {
								Required:         true,
								Type:             schema.TypeInt,
								ValidateDiagFunc: validations.IntBetween(1, 3600),
								Description:      "Graceful Restart Stale Timer",
							},
							"enable_bgp": {
								Type:     schema.TypeBool,
								Optional: true,
								Default:  true,
							},
						},
					},
				},
				"fail_over": {
					Type:     schema.TypeString,
					Required: true,
					ValidateDiagFunc: validations.StringInSlice([]string{
						"NON_PREEMPTIVE", "PREEMPTIVE",
					}, false),
					Description: "Failover. Available values are 'PREEMPTIVE' or 'NON_PREEMPTIVE'",
				},
				"ha_mode": {
					Type:     schema.TypeString,
					Required: true,
					ValidateDiagFunc: validations.StringInSlice([]string{
						"ACTIVE_ACTIVE", "ACTIVE_STANDBY",
					}, false),
					Description: "HA Mode. Available values are 'ACTIVE_ACTIVE' or 'ACTIVE_STANDBY'",
					ForceNew:    true,
				},
				"route_redistribution_tier0": {
					Type:     schema.TypeList,
					MaxItems: 1,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"tier0_static": {
								Type:        schema.TypeBool,
								Optional:    true,
								Default:     false,
								Description: "Static Routes",
							},
							"tier0_nat": {
								Type:        schema.TypeBool,
								Optional:    true,
								Default:     false,
								Description: "NAT IP",
							},
							"tier0_ipsec_local_ip": {
								Type:        schema.TypeBool,
								Optional:    true,
								Default:     false,
								Description: "IP Sec Local IP",
							},
							"tier0_dns_forwarder_ip": {
								Type:        schema.TypeBool,
								Optional:    true,
								Default:     false,
								Description: "DNS Forwarder IP",
							},
							"tier0_service_interface": {
								Type:        schema.TypeBool,
								Optional:    true,
								Default:     false,
								Description: "Service Interface Subnet",
							},
							"tier0_external_interface": {
								Type:        schema.TypeBool,
								Optional:    true,
								Default:     false,
								Description: "External Interface Subnet",
							},
							"tier0_loopback_interface": {
								Type:        schema.TypeBool,
								Optional:    true,
								Default:     false,
								Description: "Loopback Interface Subnet",
							},
							"tier0_segment": {
								Type:        schema.TypeBool,
								Optional:    true,
								Default:     false,
								Description: "Connected Segment",
							},
						},
					},
				},
				"route_redistribution_tier1": {
					Type:     schema.TypeList,
					MaxItems: 1,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"tier1_dns_forwarder_ip": {
								Type:        schema.TypeBool,
								Optional:    true,
								Default:     false,
								Description: "DNS Forwarder IP",
							},
							"tier1_static": {
								Type:        schema.TypeBool,
								Optional:    true,
								Default:     false,
								Description: "Static Routes",
							},
							"tier1_lb_vip": {
								Type:        schema.TypeBool,
								Optional:    true,
								Default:     false,
								Description: "LB VIP",
							},
							"tier1_nat": {
								Type:        schema.TypeBool,
								Optional:    true,
								Description: "NAT IP",
							},
							"tier1_lb_snat": {
								Type:        schema.TypeBool,
								Optional:    true,
								Default:     false,
								Description: "LB SNAT IP",
							},
							"tier1_ipsec_local_endpoint": {
								Type:        schema.TypeBool,
								Optional:    true,
								Default:     false,
								Description: "IPSec Local Endpoint",
							},
							"tier1_service_interface": {
								Type:        schema.TypeBool,
								Optional:    true,
								Default:     false,
								Description: "Service Interface Subnet",
							},
							"tier1_segment": {
								Type:        schema.TypeBool,
								Optional:    true,
								Default:     false,
								Description: "Connected Segment",
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
		Type:          schema.TypeList,
		Optional:      true,
		Description:   "Tier1 Gateway configuration",
		ConflictsWith: []string{"tier0_config"},
		MaxItems:      1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"tier0_gateway": {
					Type:     schema.TypeString,
					Optional: true,
					Description: "Provider ID of the Tier0 Gateway. Use Tier0 Router's " +
						" .provider_id  here.",
				},
				"edge_cluster": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Edge Cluster",
				},
				"fail_over": {
					Type:     schema.TypeString,
					Required: true,
					ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
						"NON_PREEMPTIVE", "PREEMPTIVE",
					}, false)),
					Description: "Failover. Available values are 'PREEMPTIVE' or 'NON_PREEMPTIVE'",
				},
				"route_advertisement": {
					Type:     schema.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"tier1_connected": {
								Type:        schema.TypeBool,
								Optional:    true,
								Default:     false,
								Description: "Connected Routes",
							},
							"tier1_nat": {
								Type:        schema.TypeBool,
								Optional:    true,
								Default:     false,
								Description: "NAT IPs",
							},
							"tier1_static_routes": {
								Type:        schema.TypeBool,
								Optional:    true,
								Default:     false,
								Description: "Static Routes",
							},
							"tier1_lb_vip": {
								Type:        schema.TypeBool,
								Optional:    true,
								Default:     false,
								Description: "LB VIP Routes",
							},
							"tier1_lb_snat": {
								Type:        schema.TypeBool,
								Optional:    true,
								Default:     false,
								Description: "LB SNAT IP Routes",
							},
							"tier1_dns_forwarder_ip": {
								Type:        schema.TypeBool,
								Optional:    true,
								Default:     false,
								Description: "DNS Forwarder IP Routes",
							},
							"tier1_ipsec_local_endpoint": {
								Type:        schema.TypeBool,
								Optional:    true,
								Default:     true,
								Description: "IPSec Local Endpoint",
							},
						},
					},
				},
			},
		},
	}
}

func RouterNatRuleSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		// ValidateDiagFunc: validations.ValidateUniqueNameInList,
		Description: `NAT Rules for the specific router configuration. Please note that changing
		order of nat_rule list will result into unwanted behaviour.`,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "ID of the NAT rule.",
				},
				"name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Name of the NAT rule.",
				},
				"description": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Description for the NAT rule.",
				},
				"enabled": {
					Type:        schema.TypeBool,
					Default:     false,
					Optional:    true,
					Description: "If true then NAT rule will be active/enabled.",
				},
				"config": {
					Type:        schema.TypeList,
					MaxItems:    1,
					Required:    true,
					Description: "NAT configurations",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"action": {
								Type: schema.TypeString,
								ValidateDiagFunc: validations.StringInSlice([]string{
									"DNAT", "SNAT",
								}, false),
								Required:    true,
								Description: "Supported values are DNAT and SNAT",
							},
							"service": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "Type of the service",
							},
							"firewall": {
								Type:     schema.TypeString,
								Optional: true,
								Default:  "MATCH_INTERNAL_ADDRESS",
								ValidateDiagFunc: validations.StringInSlice([]string{
									"MATCH_EXTERNAL_ADDRESS", "MATCH_INTERNAL_ADDRESS", "BYPASS",
								}, false),
								// "MATCH_INTERNAL_ADDRESS",
							},
							// This field will added on later versions
							// "scope": {
							// 	Type:        schema.TypeString,
							// 	Optional:    true,
							// 	Description: "Scope to particular router interface",
							// },
							"logging": {
								Type:     schema.TypeBool,
								Optional: true,
							},
						},
					},
				},
				"source_network": {
					Type:             schema.TypeString,
					Optional:         true,
					ValidateDiagFunc: validations.ValidateCidr,
					Description:      "Source Network CIDR Address",
				},
				"destination_network": {
					Type:             schema.TypeString,
					Optional:         true,
					ValidateDiagFunc: validations.ValidateCidr,
					Description:      "Destination Network CIDR Address",
				},
				"translated_network": {
					Type:             schema.TypeString,
					Required:         true,
					ValidateDiagFunc: validations.ValidateCidr,
					Description:      "Translated Network CIDR Address",
				},
				"translated_ports": {
					Type:        schema.TypeInt,
					Optional:    true,
					Description: "Translated Network Port",
				},
				"priority": {
					Type:             schema.TypeInt,
					Optional:         true,
					Default:          100,
					Description:      "Priority for the rule",
					ValidateDiagFunc: validations.IntAtLeast(1),
				},
			},
		},
	}
}
