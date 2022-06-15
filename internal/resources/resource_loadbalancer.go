// (C) Copyright 2022 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"

	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/resources/validations"
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
				Description: "Type of Network loadbalancer",
				Computed:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "creating Network loadbalancer",
			},
			"network_server_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "NSX-T Integration ID",
				ForceNew:    true,
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
				Required:    true,
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
				Required:    true,
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
							Type:             schema.TypeString,
							ValidateDiagFunc: validations.StringInSlice([]string{"SMALL", "MEDIUM", "LARGE"}, false),
							Optional:         true,
							Default:          true,
							Description:      "Network Loadbalancer Supported values are `SMALL`, `MEDIUM`, `LARGE`",
						},
						"loglevel": {
							Type:             schema.TypeString,
							ValidateDiagFunc: validations.StringInSlice([]string{"DEBUG", "INFO", "WARNING", "ERROR", "CRITICAL", "ALERT", "EMERGENCY"}, false),
							Optional:         true,
							Default:          true,
							Description:      "Network Loadbalancer Supported values are `DEBUG`, `INFO`, `WARNING`, `ERROR`, `CRITICAL`, `ALERT`, `EMERGENCY`",
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
		ReadContext:   loadbalancerReadContext,
		UpdateContext: loadbalancerReadContext,
		CreateContext: loadbalancerCreateContext,
		DeleteContext: loadbalancerDeleteContext,
		Description: `loadbalancer resource facilitates creating,
		and deleting NSX-T  Network Load Balancers.`,
	}
}

func loadbalancerReadContext(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

func loadbalancerCreateContext(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(rd)
	if err := c.CmpClient.LoadBalancer.Create(ctx, data, meta); err != nil {
		return diag.FromErr(err)
	}

	return loadbalancerReadContext(ctx, rd, meta)
}

func loadbalancerDeleteContext(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
