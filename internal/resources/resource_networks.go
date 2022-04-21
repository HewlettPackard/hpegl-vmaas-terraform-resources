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
