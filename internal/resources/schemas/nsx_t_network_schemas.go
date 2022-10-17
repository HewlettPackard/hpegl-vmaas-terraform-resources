package schemas

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	DSNetworkPool = "hpegl_vmaas_network_pool"
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
				"pool_id": {
					Type:        schema.TypeInt,
					Optional:    true,
					Description: "Pool ID can be obtained with " + DSNetworkPool + " data source.",
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
	}
}
