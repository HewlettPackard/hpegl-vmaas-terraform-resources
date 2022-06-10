// (C) Copyright 2022 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"

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
			"sendVersion": {
				Type:        schema.TypeInt,
				Description: "Network loadbalancer Monitor http version",
				Computed:    true,
			},
			"sendData": {
				Type:        schema.TypeString,
				Description: "Network loadbalancer Monitor Send info",
				Computed:    true,
			},
			"receiveData": {
				Type:        schema.TypeString,
				Description: "Network loadbalancer Monitor receive info",
				Computed:    true,
			},
			"receiveCode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Network loadbalancer Monitor receive status codes like 200,300,301,302,304,307",
			},
			"monitorDestination": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Network loadbalancer Monitor destination",
			},
			"monitorReverse": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Network loadbalancer Monitor Reverse",
			},
			"monitorTransparent": {
				Type:        schema.TypeBool,
				Computed:    true,		
				Description: "Network loadbalancer Monitor transparent",
			},
			"monitorAdaptive": {
				Type:        schema.TypeBool,
				Computed:    true,		
				Description: "Network loadbalancer Monitor adaptive",
			},
			"fallCount": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Network loadbalancer Monitor fall counts",
			},
			"riseCount": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Network loadbalancer Monitor rise counts",
			},
			"aliasPort": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Network loadbalancer Monitor alias port",
			},
			"monitorType": {
				Type: schema.TypeString,
				ValidateDiagFunc: validations.StringInSlice([]string{
					"LBHttpMonitorProfile", "LBHttpsMonitorProfile", "LBIcmpMonitorProfile",
					 "LBPassiveMonitorProfile","LBTcpMonitorProfile","LBUdpMonitorProfile"
				}, false),
				Computed:    true,
				Description: "Network Loadbalancer Supported values are `LBHttpMonitorProfile`, 
				`LBHttpsMonitorProfile`, `LBIcmpMonitorProfile`, `LBPassiveMonitorProfile`,
				 `LBTcpMonitorProfile`, `LBUdpMonitorProfile`",
			},
		},
		ReadContext:   LBMonitorReadContext,
		Description: `The ` + DSLoadBalancer + ` monitor data source can be used to discover the ID of a hpegl vmaas network load balancer.
		This can then be used with resources or data sources that require a ` + DSLoadBalancer + `,
		monitor such as the ` + ResLoadBalancer + ` monitor resource.`,
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
	err = c.CmpClient.DSLBMonitor.Read(ctx, data, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
