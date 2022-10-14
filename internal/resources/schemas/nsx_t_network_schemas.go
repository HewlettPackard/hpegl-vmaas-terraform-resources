package schemas

import (
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/resources/validations"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	DSGroup         = "hpegl_vmaas_group"
	DSNetworkPool   = "hpegl_vmaas_network_pool"
	DSNetworkDomain = "hpegl_vmaas_network_domain"
	DSNetworkProxy  = "hpegl_vmaas_network_proxy"
	DSTransportZone = "hpegl_vmaas_transport_zone"
	DSRouter        = "hpegl_vmaas_router"
)

func StaticNetworkSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "static Network configuration",
		MaxItems:    1,
		ExactlyOneOf: []string{
			"static_network",
			"dhcp_network",
		},
		ConflictsWith: []string{
			"dhcp_network",
		},
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Name of the NSX-T Static Segment to be created.",
				},
				"description": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Description of the network to be created.",
				},
				"display_name": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Display name of the NSX-T network.",
				},
				"group_id": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Group ID of the Network. Please use " + DSGroup + " data source to retrieve ID or pass `shared`.",
				},
				"code": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Network Type code",
				},
				"type_id": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "Type ID for the NSX-T Network.",
				},
				"pool_id": {
					Type:        schema.TypeInt,
					Optional:    true,
					Description: "Pool ID can be obtained with " + DSNetworkPool + " data source.",
				},
				"external_id": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "External ID of the network",
				},
				"internal_id": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Internal ID of the network",
				},
				"unique_id": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Unique ID of the network",
				},
				"gateway": {
					Type:             schema.TypeString,
					Optional:         true,
					Description:      "Gateway IP address of the network",
					ValidateDiagFunc: validations.ValidateIPAddress,
				},
				"primary_dns": {
					Type:             schema.TypeString,
					Optional:         true,
					Description:      "Primary DNS IP Address",
					ValidateDiagFunc: validations.ValidateIPAddress,
				},
				"secondary_dns": {
					Type:             schema.TypeString,
					Optional:         true,
					Description:      "Secondary DNS IP Address",
					ValidateDiagFunc: validations.ValidateIPAddress,
				},
				"cidr": {
					Type:             schema.TypeString,
					Required:         true,
					Description:      "Gateway Classless Inter-Domain Routing (CIDR) of the network",
					ValidateDiagFunc: validations.ValidateCidr,
				},
				"active": {
					Type:        schema.TypeBool,
					Optional:    true,
					Default:     true,
					Description: "Activate (`true`) or disable (`false`) the network",
				},
				"scan_network": {
					Type:        schema.TypeBool,
					Optional:    true,
					Default:     false,
					Description: "Scan Network",
				},
				// "dhcp_server": {
				// 	Type:        schema.TypeBool,
				// 	Optional:    true,
				// 	Default:     false,
				// 	Description: "Enable DHCP Server.",
				// },
				"appliance_url_proxy_bypass": {
					Type:        schema.TypeBool,
					Optional:    true,
					Default:     true,
					Description: "Bypass Proxy for Appliance URL",
				},
				"no_proxy": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "List of IP addresses or name servers for which to exclude proxy traversal.",
				},
				"domain_id": {
					Type:        schema.TypeInt,
					Optional:    true,
					Description: "ID of the Network domain. Use " + DSNetworkDomain + " datasource to obtain the ID.",
				},
				"proxy_id": {
					Type:        schema.TypeInt,
					Optional:    true,
					Description: "Network Proxy ID. Use " + DSNetworkProxy + " data source to obtain the ID.",
				},
				"search_domains": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Search Domains",
				},
				"allow_static_override": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "If set to true, network will allow static override",
				},
				"resource_permissions": {
					Type:     schema.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"all": {
								Type:        schema.TypeBool,
								Optional:    true,
								Default:     true,
								Description: "Pass `true` to allow access to all groups.",
							},
							"sites": {
								Type:        schema.TypeList,
								Optional:    true,
								Description: "List of sites/groups",
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"id": {
											Type:        schema.TypeInt,
											Required:    true,
											Description: "ID of the site/group",
										},
										"default": {
											Type:        schema.TypeBool,
											Optional:    true,
											Default:     false,
											Description: "Group Default Selection",
										},
									},
								},
							},
						},
					},
				},
				"config": {
					Type:        schema.TypeList,
					Optional:    true,
					Description: "Network configuration",
					MaxItems:    1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"connected_gateway": {
								Type:     schema.TypeString,
								Required: true,
								Description: "Connected Gateway. Pass Provider ID of the Tier1 gateway. Use " + DSRouter +
									".provider_id  here.",
							},
							"vlan_id": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "VLAN IDs eg. `0,3-5`. Use this field for VLAN based segments.",
							},
						},
					},
				},
				"status": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Status of the network",
				},
				"scope_id": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Transport Zone ID. Use " + DSTransportZone + " Data source's `provider_id` here.",
				},
			},
		},
	}
}

func DhcpNetworkSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "DHCP Network configuration",
		MaxItems:    1,
		ExactlyOneOf: []string{
			"static_network",
			"dhcp_network",
		},
		ConflictsWith: []string{
			"static_network",
		},
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Name of the NSX-T DHCP Segment to be created.",
				},
				"description": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Description of the network to be created.",
				},
				"display_name": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Display name of the NSX-T network.",
				},
				"gateway": {
					Type:             schema.TypeString,
					Optional:         true,
					Description:      "Gateway IP address of the network",
					ValidateDiagFunc: validations.ValidateIPAddress,
				},
				"primary_dns": {
					Type:             schema.TypeString,
					Optional:         true,
					Description:      "Primary DNS IP Address",
					ValidateDiagFunc: validations.ValidateIPAddress,
				},
				"secondary_dns": {
					Type:             schema.TypeString,
					Optional:         true,
					Description:      "Secondary DNS IP Address",
					ValidateDiagFunc: validations.ValidateIPAddress,
				},
				"gateway_cidr": {
					Type:             schema.TypeString,
					Required:         true,
					Description:      "Gateway Classless Inter-Domain Routing (CIDR) of the network",
					ValidateDiagFunc: validations.ValidateCidr,
				},
				"active": {
					Type:        schema.TypeBool,
					Optional:    true,
					Default:     true,
					Description: "Activate (`true`) or disable (`false`) the network",
				},
				"scan_network": {
					Type:        schema.TypeBool,
					Optional:    true,
					Default:     false,
					Description: "Scan Network",
				},
				"dhcp_enabled": {
					Type:        schema.TypeBool,
					Optional:    true,
					Default:     true,
					Description: "Enable DHCP Server.",
				},
				"bypass_proxy_for_appliance_url": {
					Type:        schema.TypeBool,
					Optional:    true,
					Default:     true,
					Description: "Bypass Proxy for Appliance URL",
				},
				"no_proxy": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "List of IP addresses or name servers for which to exclude proxy traversal.",
				},
				"search_domains": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Search Domains",
				},
				"allow_ip_override": {
					Type:        schema.TypeBool,
					Optional:    true,
					Default:     false,
					Description: "If set to true, network will allow static override",
				},
				"group": {
					Type:        schema.TypeList,
					Optional:    true,
					Description: "Group ID",
					MaxItems:    1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"id": {
								Type:     schema.TypeString,
								Required: true,
								Description: "Group ID. Get the Group ID Use " + DSGroup +
									"Pass `shared` to use this object across all the Groups.",
							},
						},
					},
				},
				"domain": {
					Type:        schema.TypeList,
					Optional:    true,
					Description: "Domain ID",
					MaxItems:    1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"id": {
								Type:     schema.TypeInt,
								Optional: true,
								Description: "Domain ID. Get the Network Domain ID Use " + DSNetworkDomain +
									"to get Network Domain ID",
							},
						},
					},
				},
				"network_proxy": {
					Type:        schema.TypeList,
					Optional:    true,
					Description: "Network Proxy ID",
					MaxItems:    1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"id": {
								Type:     schema.TypeInt,
								Optional: true,
								Description: "Network Proxy ID. Get the Network proxy ID Use " + DSNetworkProxy +
									"to get Network Proxy ID",
							},
						},
					},
				},
				"resource_permissions": {
					Type:     schema.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"all": {
								Type:        schema.TypeBool,
								Optional:    true,
								Default:     true,
								Description: "Pass `true` to allow access to all groups.",
							},
							"sites": {
								Type:        schema.TypeList,
								Optional:    true,
								Description: "List of sites/groups",
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"id": {
											Type:        schema.TypeInt,
											Required:    true,
											Description: "ID of the site/group",
										},
										"default": {
											Type:        schema.TypeBool,
											Optional:    true,
											Default:     false,
											Description: "Group Default Selection",
										},
									},
								},
							},
						},
					},
				},
				"config": {
					Type:        schema.TypeList,
					Required:    true,
					Description: "Network configuration",
					MaxItems:    1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"connected_gateway": {
								Type:     schema.TypeString,
								Optional: true,
								Description: "Connected Gateway. Pass Provider ID of the Tier1 gateway. Use " + DSRouter +
									".provider_id  here.",
							},
							"vlan": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "VLAN IDs eg. `0,3-5`. Use this field for VLAN based segments.",
							},
							"dhcp_type": {
								Type:        schema.TypeString,
								Required:    true,
								Description: `DHCP Server type. Supported Values are "dhcpLocal", "dhcpRelay", "gatewayDhcp"`,
							},
							"dhcp_server": {
								Type:        schema.TypeString,
								Required:    true,
								Description: "DHCP Server Config ID",
							},
							"dhcp_server_address": {
								Type:     schema.TypeString,
								Optional: true,
								Description: "DHCP Server address and its CIDR. This address must not overlap the" +
									"ip-ranges of the subnet, or the gateway address of the subnet," +
									"or the DHCP static-binding addresses of this segment",
							},
							"dhcp_range": {
								Type:        schema.TypeString,
								Required:    true,
								Description: "DHCP server IP Address range",
							},
							"dhcp_lease_time": {
								Type:        schema.TypeString,
								Required:    true,
								Description: "DHCP Server default lease time",
							},
						},
					},
				},
				"transport_zone": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Transport Zone ID. Use " + DSTransportZone + " Data source's `provider_id` here.",
				},
			},
		},
	}
}
