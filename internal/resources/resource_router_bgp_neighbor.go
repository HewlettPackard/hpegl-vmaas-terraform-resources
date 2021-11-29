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

func RouterBgpNeighbor() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"router_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Parent router ID, router_id can be obtained by using router datasource/resource.",
				ForceNew:    true,
			},
			"ip_address": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validations.ValidateIPAddress,
				Description:      "IP Address.",
			},
			"remote_as": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Remote AS number.",
			},
			"keepalive": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     60,
				Description: "Keep Alive Time.",
			},
			"holddown": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     60,
				Description: "Hold Down Time.",
			},
			"router_filtering_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "IP Address Family. Supported Values are `IPV4`, `IPV6` and `L2VPN_EVPN`",
			},
			"router_filtering_in": {
				Type:     schema.TypeString,
				Optional: true,
				Description: `In Filter. Supported Values are "DNS_FORWARDER_IP",
				"EVPN_TEP_IP", "EXTERNAL_INTERFACE", "INTERNAL_TRANSIT_SUBNET",
				"IPSEC_LOCAL_IP", "LOOPBACK_INTERFACE", "prefixlist-out-default",
				"ROUTER_LINK", "SEGMENT", "SERVICE_INTERFACE"`,
			},
			"router_filtering_out": {
				Type:     schema.TypeString,
				Optional: true,
				Description: `Out Filter. Supported Values are "DNS_FORWARDER_IP",
				"EVPN_TEP_IP", "EXTERNAL_INTERFACE", "INTERNAL_TRANSIT_SUBNET",
				"IPSEC_LOCAL_IP", "LOOPBACK_INTERFACE", "prefixlist-out-default",
				"ROUTER_LINK", "SEGMENT", "SERVICE_INTERFACE"`,
			},
			"bfd_enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "BFD Enabled.",
				//				Default:     false,
			},
			"bfd_interval": {
				Type:     schema.TypeInt,
				Required: true,
				//				Default:     1000,
				Description: "BFD Interval(ms).",
			},
			"bfd_multiple": {
				Type:     schema.TypeInt,
				Optional: true,
				//				Default:     3,
				Description: "BFD Multiplier.",
			},
			"allow_as_in": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Allow-as-in.",
				Default:     false,
			},
			"hop_limit": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Max Hop Limit.",
				//				Default:     1,
			},
			"restart_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Graceful Restart. Supported Values are "HELPER_ONLY" "GR_AND_HELPER" "DISABLE".`,
			},
			"config": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Interface configurations",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source_addresses": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Source Addresses. This can be retrieved using Network Router Data Source",
							Elem: &schema.Schema{
								Type:             schema.TypeString,
								ValidateDiagFunc: validations.ValidateIPAddress,
							},
						},
					},
				},
			},
		},
		ReadContext:   routerBgpNeighborReadContext,
		CreateContext: routerBgpNeighborCreateContext,
		UpdateContext: routerBgpNeighborUpdateContext,
		DeleteContext: routerBgpNeighborDeleteContext,
		Description: `Router Bgp Neighbor resource facilitates creating,
		updating and deleting NSX-T Network Router BGP Neighbors.`,
	}
}

func routerBgpNeighborReadContext(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(rd)
	if err := c.CmpClient.RouterBgpNeighbor.Read(ctx, data, meta); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func routerBgpNeighborCreateContext(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(rd)
	if err := c.CmpClient.RouterBgpNeighbor.Create(ctx, data, meta); err != nil {
		return diag.FromErr(err)
	}

	return routerBgpNeighborReadContext(ctx, rd, meta)
}

func routerBgpNeighborUpdateContext(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(rd)
	if err := c.CmpClient.RouterBgpNeighbor.Update(ctx, data, meta); err != nil {
		return diag.FromErr(err)
	}

	return routerBgpNeighborReadContext(ctx, rd, meta)
}

func routerBgpNeighborDeleteContext(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(rd)
	if err := c.CmpClient.RouterBgpNeighbor.Delete(ctx, data, meta); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
