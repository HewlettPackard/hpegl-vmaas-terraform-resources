// (C) Copyright 2022 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"

	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/utils"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func LoadBalancerData() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: f(generalNamedesc, "ResLoadBalancer", "ResLoadBalancer"),
			},
			"type": {
				Type:        schema.TypeString,
				Description: "This field can be used as type for the " + ResLoadBalancer,
				Computed:    true,
			},
			"network_server_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "NSX-T Integration ID",
				ForceNew:    true,
			},
			"config": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Network Load Balancer Configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"admin_state": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "If `true` then admin State rule will be active/enabled.",
						},
						"size": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `In Filter. Supported Values are "SMALL", "MEDIUM", "LARGE"`,
						},
						"loglevel": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `In Filter. Supported Values are "DEBUG", "INFO", "WARNING", "ERROR", "CRITICAL", "ALERT", "EMERGENCY"`,
						},
						"tier1": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Network Loadbalancer NSX-T tier1 gateway",
						},
					},
				},
			},
		},
		ReadContext: LoadBalancerReadContext,
		Description: `The ` + DSLoadBalancer + ` data source can be used to discover the ID of a hpegl vmaas network load balancer.
		This can then be used with resources or data sources that require a ` + DSLoadBalancer + `,
		such as the ` + ResLoadBalancer + ` resource.`,
		SchemaVersion:  0,
		StateUpgraders: nil,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func LoadBalancerReadContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(d)
	err = c.CmpClient.LoadBalancer.Read(ctx, data, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
