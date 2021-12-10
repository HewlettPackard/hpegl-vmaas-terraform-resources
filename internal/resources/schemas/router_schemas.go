package schemas

import (
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/resources/validations"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
					Optional: true,
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
				"edge_cluster": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Edge Cluster. Use EdgeCluster's provided_id here using EdgeCluster Data Source.",
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
					Description: "Edge Cluster. Use EdgeCluster's provided_id here using EdgeCluster Data Source.",
				},
				"fail_over": {
					Type:     schema.TypeString,
					Optional: true,
					ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
						"NON_PREEMPTIVE", "PREEMPTIVE",
					}, false)),
					RequiredWith: []string{"tier1_config.edge_cluster"},
					Description:  "Failover. Available values are 'PREEMPTIVE' or 'NON_PREEMPTIVE'",
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
