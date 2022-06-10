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

func LoadBalancerMonitors() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Network loadbalancer Monitor name",
			},
			"description": {
				Type:        schema.TypeInt,
				Description: "Creating the Network Load balancer Monitor.",
				Required:    true,
			},
			"monitorTimeout": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Timeout for Network loadbalancer Monitor",
			},
			"monitorInterval": {
				Type:        schema.TypeInt,
				Description: "Interval time for Network loadbalancer Monitor",
				Required:    true,
			},
			"sendVersion": {
				Type:        schema.TypeInt,
				Description: "Network loadbalancer Monitor http version",
				Required:    true,
			},
			"sendData": {
				Type:        schema.TypeString,
				Description: "Network loadbalancer Monitor Send info",
				Required:    true,
			},
			"receiveData": {
				Type:        schema.TypeString,
				Description: "Network loadbalancer Monitor receive info",
				Required:    true,
			},
			"receiveCode": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Network loadbalancer Monitor receive status codes like 200,300,301,302,304,307",
			},
			"monitorDestination": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Network loadbalancer Monitor destination",
			},
			"monitorReverse": {
				Type:        schema.TypeBool,
				Required:    true,
				Default:     true,
				Description: "Network loadbalancer Monitor Reverse",
			},
			"monitorTransparent": {
				Type:        schema.TypeBool,
				Required:    true,
				Default:     true,
				Description: "Network loadbalancer Monitor transparent",
			},
			"monitorAdaptive": {
				Type:        schema.TypeBool,
				Required:    true,
				Default:     true,
				Description: "Network loadbalancer Monitor adaptive",
			},
			"fallCount": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Network loadbalancer Monitor fall counts",
			},
			"riseCount": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Network loadbalancer Monitor rise counts",
			},
			"aliasPort": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Network loadbalancer Monitor alias port",
			},
			"monitorType": {
				Type: schema.TypeString,
				ValidateDiagFunc: validations.StringInSlice([]string{
					"LBHttpMonitorProfile", "LBHttpsMonitorProfile", "LBIcmpMonitorProfile",
					 "LBPassiveMonitorProfile","LBTcpMonitorProfile","LBUdpMonitorProfile"
				}, false),
				Required:    true,
				Description: "Network Loadbalancer Supported values are `LBHttpMonitorProfile`, 
				`LBHttpsMonitorProfile`, `LBIcmpMonitorProfile`, `LBPassiveMonitorProfile`,
				 `LBTcpMonitorProfile`, `LBUdpMonitorProfile`",
			},
		},
		ReadContext:   loadbalancerMonitorReadContext,
		CreateContext: loadbalancerMonitorCreateContext,
		DeleteContext: loadbalancerMonitorDeleteContext,
		Description: `loadbalancer Monitor resource facilitates creating,
		and deleting NSX-T  Network Load Balancers.`,
	}
}

func loadbalancerMonitorReadContext(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(rd)
	if err := c.CmpClient.ResLoadBalancerMonitors.Read(ctx, data, meta); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func loadbalancerMonitorCreateContext(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(rd)
	if err := c.CmpClient.ResLoadBalancerMonitors.Create(ctx, data, meta); err != nil {
		return diag.FromErr(err)
	}

	return loadbalancerMonitorReadContext(ctx, rd, meta)
}

func loadbalancerMonitorDeleteContext(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(rd)
	if err := c.CmpClient.ResLoadBalancerMonitors.Delete(ctx, data, meta); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
