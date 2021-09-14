// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"

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
			"code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "network code",
			},
			"cloud_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Cloud ID or the zone ID",
			},
			"type_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Type id for the NSX-T. This value will be constant always",
			},
			"external_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "External ID ",
			},
			"internal_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Internal ID",
			},
			"unique_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Unique ID",
			},
			"gateway": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Gateway address for the network",
			},
			"netmask": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Netmask address for the network",
			},
			"primary_dns": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Primary DNS",
			},
			"secondary_dns": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Secondary DNS",
			},
			"cidr": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "CIDR of the network",
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
				Type:        schema.TypeString,
				Optional:    true,
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
				Required:    true,
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
							Elem: &schema.Schema{
								Type: schema.TypeInt,
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

	return nil
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
