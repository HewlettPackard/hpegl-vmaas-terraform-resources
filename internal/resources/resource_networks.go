// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"

	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/resources/validations"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/utils"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Network() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the NSX-T Segment to be created.",
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
			"type_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Type ID for the NSX-T Network.",
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
			"type": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Type ID",
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Required: true,
							Description: "Network Type ID. Get the Network Type ID Use " + DSNetwork +
								"to Gets All Network Types API.",
						},
					},
				},
			},
			"network_server": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Network Server ID",
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Get the ID Use " + DSNetwork + "where `serviceType` is set to `networkServer`",
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
							Required: true,
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
							Required: true,
							Description: "Network Proxy ID. Get the Network proxy ID Use " + DSNetworkProxy +
								"to get Network Proxy ID",
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
							Description: "DHCP Server type.",
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
			"transport_zone": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Transport Zone ID. Use " + DSTransportZone + " Data source's `provider_id` here.",
			},
		},
		SchemaVersion: 0,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		ReadContext:   resNetworkReadContext,
		CreateContext: resNetworkCreateContext,
		UpdateContext: resNetworkUpdateContext,
		DeleteContext: resNetworkDeleteContext,
		Description: `Network resource facilitates creating,
		updating and deleting NSX-T Networks.`,
	}
}

func resNetworkReadContext(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(rd)
	err = c.CmpClient.ResNetwork.Read(ctx, data, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resNetworkCreateContext(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(rd)
	err = c.CmpClient.ResNetwork.Create(ctx, data, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	return resNetworkReadContext(ctx, rd, meta)
}

func resNetworkDeleteContext(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(rd)
	err = c.CmpClient.ResNetwork.Delete(ctx, data, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resNetworkUpdateContext(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(rd)
	err = c.CmpClient.ResNetwork.Update(ctx, data, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
