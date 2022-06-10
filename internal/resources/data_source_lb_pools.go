// (C) Copyright 2022 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"

	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/utils"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func LBPoolData() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: f(generalNamedesc, "ResLoadBalancerPool", "ResLoadBalancerPool"),
			},
			"minActive": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "minimum active members for the Network loadbalancer pool",
			},
			"vipBalance": {
				Type: schema.TypeString,
				            ValidateDiagFunc: validations.StringInSlice([]string{
					            "ROUND_ROBIN", "WEIGHTED_ROUND_ROBIN", " LEAST_CONNECTION",
								 "WEIGHTED_LEAST_CONNECTION","IP_HASH",
				            }, false),
				            Computed:    true,
				            Description: "Network Loadbalancer Supported values are `ROUND_ROBIN`, 
				                `WEIGHTED_ROUND_ROBIN`, `LEAST_CONNECTION`, `WEIGHTED_LEAST_CONNECTION`,`IP_HASH`",
			},
			"config": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "pool Configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"snatTranslationType": {
							Type: schema.TypeString,
				            ValidateDiagFunc: validations.StringInSlice([]string{
					            "LBSnatAutoMap", "LBSnatDisabled", "LBSnatIpPool"
				            }, false),
				            Computed:    true,
				            Description: "Network Loadbalancer Supported values are `LBSnatAutoMap`, 
				                `LBSnatDisabled`, `LBSnatIpPool`",
						},
						"snatIpAddress": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "snatIpAddress for Network loadbalancer pool",
						},
						"memberGroup": {
							Type:        schema.TypeList,
				            Computed:    true,
				            Description: "memberGroup Configuration",
				            Elem: &schema.Resource{
					            Schema: map[string]*schema.Schema{
						            "name": {
							            Type:        schema.TypeString,
										Computed:    true,
										Description: "name of the member group",
						            },
						            "path": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "path of the member group",
									},
									"ipRevisionFilter": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ipRevisionFilter of the member group",
									},
									"port": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "port of the member group",
									},
								},
							},
						},	
					},
				},
			},
		},
		ReadContext:   LBPoolReadContext,
		Description: `The ` + DSLBPool + ` data source can be used to discover the ID of a hpegl vmaas network load balancer.
		This can then be used with resources or data sources that require a ` + DSLBPool + `,
		such as the ` + ResLoadBalancerPools + ` resource.`,
		SchemaVersion:  0,
		StateUpgraders: nil,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func LBPoolReadContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(d)
	err = c.CmpClient.DSLBPool.Read(ctx, data, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
