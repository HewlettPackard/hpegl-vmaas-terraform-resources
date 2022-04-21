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

func RouterRoute() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"router_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Parent router ID, router_id can be obtained by using router datasource/resource.",
				ForceNew:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the route.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description for the route.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Default:     true,
				Optional:    true,
				Description: "If `true` then route will be active/enabled.",
			},
			"default_route": {
				Type:        schema.TypeBool,
				Default:     false,
				Optional:    true,
				Description: "If `true` then the route will considered as the default route.",
			},
			"network": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validations.ValidateCidr,
				Description:      "Source Network CIDR Address",
			},
			"next_hop": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validations.ValidateIPAddress,
				Description:      "Next Hop/Destination IPv4 Address",
			},
			"mtu": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Network MTU",
			},
			"priority": {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          100,
				Description:      "Priority for the route",
				ValidateDiagFunc: validations.IntAtLeast(1),
			},
			"is_deprecated": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If parent router not found, then is_deprecated will be true",
			},
			"code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"route_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"external_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"provider_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		ReadContext:   routerRouteReadContext,
		CreateContext: routerRouteCreateContext,
		UpdateContext: routerRouteUpdateContext,
		DeleteContext: routerRouteDeleteContext,
		Description: `Router route resource facilitates creating,
		updating and deleting NSX-T Network Router routes.`,
	}
}

func routerRouteReadContext(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(rd)
	if err := c.CmpClient.RouterRoute.Read(ctx, data, meta); err != nil {
		return diag.FromErr(err)
	}
	isDeprecated := data.GetBool("is_deprecated")
	if isDeprecated {
		return diag.Diagnostics{
			{
				Severity: diag.Warning,
				Summary:  "Parent router is deleted. This resource is deprecated!!!",
			},
		}
	}

	return nil
}

func routerRouteCreateContext(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(rd)
	if err := c.CmpClient.RouterRoute.Create(ctx, data, meta); err != nil {
		return diag.FromErr(err)
	}

	return routerRouteReadContext(ctx, rd, meta)
}

func routerRouteUpdateContext(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(rd)
	if err := c.CmpClient.RouterRoute.Update(ctx, data, meta); err != nil {
		return diag.FromErr(err)
	}

	return routerRouteReadContext(ctx, rd, meta)
}

func routerRouteDeleteContext(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(rd)
	if err := c.CmpClient.RouterRoute.Delete(ctx, data, meta); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
