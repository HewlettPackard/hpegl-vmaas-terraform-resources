// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"

	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/resources/schemas"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/utils"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Router() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Network router name",
			},
			"type_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Network router type ID.",
				ForceNew:    true,
			},
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Group ID. Available values are either 'Shared' or ID fetched from " + DSGroup,
				ForceNew:    true,
			},
			"enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Can be used to enable/disable the network router",
				Default:     true,
			},
			"network_server_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "NSX-T Integration ID",
				ForceNew:    true,
			},
			"tier0_config": schemas.RouterTier0ConfigSchema(),
			"tier1_config": schemas.RouterTier1ConfigSchema(),
		},
		ReadContext:   routerReadContext,
		CreateContext: routerCreateContext,
		UpdateContext: routerUpdateContext,
		DeleteContext: routerDeleteContext,
	}
}

func routerReadContext(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(rd)
	if err := c.CmpClient.Router.Read(ctx, data, meta); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func routerCreateContext(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(rd)
	if err := c.CmpClient.Router.Create(ctx, data, meta); err != nil {
		return diag.FromErr(err)
	}

	return routerReadContext(ctx, rd, meta)
}

func routerUpdateContext(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(rd)
	if err := c.CmpClient.Router.Update(ctx, data, meta); err != nil {
		return diag.FromErr(err)
	}

	return routerReadContext(ctx, rd, meta)
}

func routerDeleteContext(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(rd)
	if err := c.CmpClient.Router.Delete(ctx, data, meta); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
