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

func LBPoolData() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: f(generalNamedesc, "ResLoadBalancerPool", "ResLoadBalancerPool"),
			},
			"min_active": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "minimum active members for the Network loadbalancer pool",
			},
			"vip_balance": {
				Type: schema.TypeString,
				ValidateDiagFunc: validations.StringInSlice([]string{
					"ROUND_ROBIN", "WEIGHTED_ROUND_ROBIN", " LEAST_CONNECTION",
					"WEIGHTED_LEAST_CONNECTION", "IP_HASH",
				}, false),
				Required:    true,
				Description: "Network Loadbalancer Supported values are `ROUND_ROBIN`,`WEIGHTED_ROUND_ROBIN`, `LEAST_CONNECTION`, `WEIGHTED_LEAST_CONNECTION`,`IP_HASH`",
			},
			"config": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "pool Configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"snat_translation_type": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validations.StringInSlice([]string{"LBSnatAutoMap", "LBSnatDisabled", "LBSnatIpPool"}, false),
							Description:      "Network Loadbalancer Supported values are `LBSnatAutoMap`,`LBSnatDisabled`, `LBSnatIpPool`"},
						"snat_ip_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "snat_ip_address for Network loadbalancer pool",
						},
					},
				},
			},
		},
		ReadContext: LBPoolReadContext,
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
	err = c.CmpClient.LoadBalancerPool.Read(ctx, data, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
