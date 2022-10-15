package schemas

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	DSRouter = "hpegl_vmaas_router"
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
				"code": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Network Type code",
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
				"status": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Status of the network",
				},
				"static_config": {
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
				"dhcp_enabled": {
					Type:        schema.TypeBool,
					Optional:    true,
					Default:     true,
					Description: "Enable DHCP Server.",
				},
				"dhcp_config": {
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
							"vlan_ids": {
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
			},
		},
	}
}
