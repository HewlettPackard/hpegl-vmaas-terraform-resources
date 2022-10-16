// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"

	diffvalidation "github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/resources/diffValidation"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/resources/schemas"
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
			"scope_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Transport Zone ID. Use " + DSTransportZone + " Data source's `provider_id` here.",
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
			"static_network": schemas.StaticNetworkSchema(),
			"dhcp_network":   schemas.DhcpNetworkSchema(),
		},
		SchemaVersion: 0,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		ReadContext:   resNetworkReadContext,
		CreateContext: resNetworkCreateContext,
		UpdateContext: resNetworkUpdateContext,
		DeleteContext: resNetworkDeleteContext,
		CustomizeDiff: networkCustomDiff,
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

func networkCustomDiff(ctx context.Context, diff *schema.ResourceDiff, meta interface{}) error {
	return diffvalidation.NewNetworkValidate(diff).DiffValidate()
}
