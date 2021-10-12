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
				Description: "Name of the network to be created",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the network to be created",
			},
			"display_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Display name of the network",
			},
			"group_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "ID of the group in which network associated. Please use " + DSGroup + " data source to retrieve ID",
			},
			"code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "network code",
			},
			"cloud_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Cloud ID or the zone ID",
				ForceNew:    true,
			},
			"type_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Type id for the NSX-T. Use " + DSNetworkType + " Data source for retrieving type ID",
			},
			"parent_network_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Description: "Parent network ID can be obtained with " + DSNetwork + " data source/resource. This field" +
					"should be set only when creating 'Custom Network'.",
			},
			"pool_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Pool ID can be obtained with " + DSNetworkPool + " data source. pool_id will not support with NSX-T segment",
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
				Description:      "Gateway address for the network",
				ValidateDiagFunc: validations.ValidateIPAddress,
			},
			"netmask": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Netmask address for the network",
				ValidateDiagFunc: validations.ValidateIPAddress,
			},
			"primary_dns": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Primary DNS",
				ValidateDiagFunc: validations.ValidateIPAddress,
			},
			"secondary_dns": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Secondary DNS",
				ValidateDiagFunc: validations.ValidateIPAddress,
			},
			"cidr": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "CIDR of the network",
				ValidateDiagFunc: validations.ValidateCidr,
			},
			"active": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Denotes the network is active or not",
			},
			"scan_network": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Dentes whether scan network",
			},
			"dhcp_server": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "DHCP server address",
			},
			"appliance_url_proxy_bypass": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Addresses of appliances to proxy bypass",
			},
			"no_proxy": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "No proxy IPs/Adrresses",
			},
			"config": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Network configuration",
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"connected_gateway": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "ID for the gateway connection",
						},
						"vlan_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "VLAN ID",
						},
					},
				},
			},
			"resource_permission": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"all": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Dentes whether provide all permissions",
						},
						"sites": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of site id",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "id for the site",
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
				Optional:    true,
				Description: "Scope ID",
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
