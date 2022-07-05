// (C) Copyright 2022 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"

	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/utils"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func LoadBalancer() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Network loadbalancer name",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of Network loadbalancer",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "creating Network loadbalancer",
			},
			"network_server_id": {
				Type:        schema.TypeInt,
				Description: "NSX-T Integration ID",
				Computed:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Network Loadbalancer configuration enabled",
				Default:     true,
			},
			"visibility": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Network Loadbalancer is public/private visibility mode",
			},
			"resource_permission": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "permission access for Loadbalancer",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"all": {
							Type:        schema.TypeBool,
							Default:     true,
							Optional:    true,
							Description: "If `true` then resource_permission rule will be active/enabled.",
						},
					},
				},
			},
			"config": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Network Load Balancer Configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"admin_state": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "If `true` then admin State rule will be active/enabled.",
						},
						"size": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "SMALL",
							Description: `In Filter. Supported Values are "SMALL", "MEDIUM", "LARGE"`,
						},
						"loglevel": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "INFO",
							Description: `In Filter. Supported Values are "DEBUG", "INFO", "WARNING", "ERROR", "CRITICAL", "ALERT", "EMERGENCY"`,
						},
						"tier1": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     true,
							Description: "Network Loadbalancer NSX-T tier1 gateway",
						},
					},
				},
			},
		},
		SchemaVersion: 0,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		ReadContext:   LoadBalancerReadContext,
		UpdateContext: LoadbalancerUpdateContext,
		CreateContext: LoadbalancerCreateContext,
		DeleteContext: LoadbalancerDeleteContext,
		Description: `loadbalancer resource facilitates creating, updating
		and deleting NSX-T  Network Load Balancers.`,
	}
}

func LoadbalancerReadContext(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(rd)
	if err := c.CmpClient.LoadBalancer.Read(ctx, data, meta); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func LoadbalancerCreateContext(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(rd)
	if err := c.CmpClient.LoadBalancer.Create(ctx, data, meta); err != nil {
		return diag.FromErr(err)
	}

	return LoadbalancerReadContext(ctx, rd, meta)
}

func LoadbalancerUpdateContext(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(rd)
	if err := c.CmpClient.LoadBalancer.Update(ctx, data, meta); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func LoadbalancerDeleteContext(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(rd)
	if err := c.CmpClient.LoadBalancer.Delete(ctx, data, meta); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
