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

func LBMonitorData() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: f(generalNamedesc, "ResLoadBalancerMonitor", "ResLoadBalancerMonitor"),
			},
			"send_version": {
				Type:        schema.TypeString,
				Description: "Network loadbalancer Monitor http version",
				Computed:    true,
			},
			"send_data": {
				Type:        schema.TypeString,
				Description: "Network loadbalancer Monitor Send info",
				Computed:    true,
			},
			"receive_data": {
				Type:        schema.TypeString,
				Description: "Network loadbalancer Monitor receive info",
				Computed:    true,
			},
			"receive_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Network loadbalancer Monitor receive status codes like 200,300,301,302,304,307",
			},
			"monitor_destination": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Network loadbalancer Monitor destination",
			},
			"monitor_reverse": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Network loadbalancer Monitor Reverse",
			},
			"monitor_transparent": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Network loadbalancer Monitor transparent",
			},
			"monitor_adaptive": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Network loadbalancer Monitor adaptive",
			},
			"fall_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Network loadbalancer Monitor fall counts",
			},
			"rise_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Network loadbalancer Monitor rise counts",
			},
			"alias_port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Network loadbalancer Monitor alias port",
			},
			"monitor_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateDiagFunc: validations.StringInSlice([]string{
					"LBHttpMonitorProfile", "LBHttpsMonitorProfile", "LBIcmpMonitorProfile",
					"LBPassiveMonitorProfile", "LBTcpMonitorProfile", "LBUdpMonitorProfile",
				}, false),
				Description: "Network Loadbalancer Supported values are `LBHttpMonitorProfile`, `LBHttpsMonitorProfile`, `LBIcmpMonitorProfile`, `LBPassiveMonitorProfile`, `LBTcpMonitorProfile`, `LBUdpMonitorProfile`",
			},
		},
		ReadContext: LBMonitorReadContext,
		Description: `The ` + DSLBMonitor + ` data source can be used to discover the ID of a hpegl vmaas network load balancer.
		This can then be used with resources or data sources that require a ` + DSLBMonitor + `,
		such as the ` + ResLoadBalancerMonitors + ` resource.`,
		SchemaVersion:  0,
		StateUpgraders: nil,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func LBMonitorReadContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
