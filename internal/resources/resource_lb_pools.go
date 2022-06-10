// (C) Copyright 2022 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"

	diffvalidation "github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/resources/diffValidation"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/resources/schemas"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/utils"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResLoadBalancerPools() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Network loadbalancer pool name",
			},
			"description": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Creating the Network loadbalancer pool.",
				ForceNew:    true,
			},
			"minActive": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "minimum active members for the Network loadbalancer pool",
				ForceNew:    true,
			},
			"vipBalance": {
				Type: schema.TypeString,
				            ValidateDiagFunc: validations.StringInSlice([]string{
					            "ROUND_ROBIN", "WEIGHTED_ROUND_ROBIN", " LEAST_CONNECTION",
								 "WEIGHTED_LEAST_CONNECTION","IP_HASH",
				            }, false),
				            Required:    true,
				            Description: "Network Loadbalancer Supported values are `ROUND_ROBIN`, 
				                `WEIGHTED_ROUND_ROBIN`, `LEAST_CONNECTION`, `WEIGHTED_LEAST_CONNECTION`,`IP_HASH`",
						},
			},
			"config": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "pool Configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"snatTranslationType": {
							Type: schema.TypeString,
				            ValidateDiagFunc: validations.StringInSlice([]string{
					            "LBSnatAutoMap", "LBSnatDisabled", "LBSnatIpPool"
				            }, false),
				            Required:    true,
				            Description: "Network Loadbalancer Supported values are `LBSnatAutoMap`, 
				                `LBSnatDisabled`, `LBSnatIpPool`",
						},
						},
						"passiveMonitorPath": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "passiveMonitorPath for Network loadbalancer pool",
						},
						"activeMonitorPaths": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "activeMonitorPaths for Network loadbalancer pool",
						},
						"tcpMultiplexing": {
							Type:        schema.TypeBool,
							Required:    true,
							Default:     true,
							Description: "tcpMultiplexing for Network loadbalancer pool",
						},
						"tcpMultiplexingNumber": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "tcpMultiplexingNumber for Network loadbalancer pool",
						},
						"snatIpAddress": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "snatIpAddress for Network loadbalancer pool",
						},
						"memberGroup": {
							Type:        schema.TypeList,
				            Required:    true,
				            Description: "memberGroup Configuration",
				            Elem: &schema.Resource{
					            Schema: map[string]*schema.Schema{
						            "name": {
							            Type:        schema.TypeString,
										Required:    true,
										Description: "name of the member group",
						            },
						            "path": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "path of the member group",
									},
									"ipRevisionFilter": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "ipRevisionFilter of the member group",
									},
									"port": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "port of the member group",
									},
								},
							},			
						
						},
					},
				},
			},
		},
		ReadContext:   loadbalancerPoolReadContext,
		CreateContext: loadbalancerPoolCreateContext,
		DeleteContext: loadbalancerPoolDeleteContext,
		Description: `loadbalancer Pool resource facilitates creating,
		and deleting NSX-T  Network Load Balancers.`,
	}
}

func loadbalancerPoolReadContext(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(rd)
	if err := c.CmpClient.ResLoadBalancerPools.Read(ctx, data, meta); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func loadbalancerPoolCreateContext(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(rd)
	if err := c.CmpClient.ResLoadBalancerPools.Create(ctx, data, meta); err != nil {
		return diag.FromErr(err)
	}

	return loadbalancerPoolReadContext(ctx, rd, meta)
}

func loadbalancerPoolDeleteContext(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(rd)
	if err := c.CmpClient.ResLoadBalancerPools.Delete(ctx, data, meta); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
