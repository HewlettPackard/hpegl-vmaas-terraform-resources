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

func LoadBalancerMonitors() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Network loadbalancer Monitor name",
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Creating the Network Load balancer Monitor.",
				Optional:    true,
			},
			"send_type": {
				Type:        schema.TypeString,
				Description: "send type method like GET,POST",
				Optional:    true,
			},
			"monitor_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Timeout for Network loadbalancer Monitor",
			},
			"monitor_interval": {
				Type:        schema.TypeInt,
				Description: "Interval time for Network loadbalancer Monitor",
				Optional:    true,
			},
			"send_version": {
				Type:        schema.TypeString,
				Description: "Network loadbalancer Monitor http version",
				Optional:    true,
			},
			"send_data": {
				Type:        schema.TypeString,
				Description: "Network loadbalancer Monitor Send info",
				Optional:    true,
			},
			"receive_data": {
				Type:        schema.TypeString,
				Description: "Network loadbalancer Monitor receive info",
				Optional:    true,
			},
			"receive_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Network loadbalancer Monitor receive status codes like 200,300,301,302,304,307",
			},
			"monitor_destination": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Network loadbalancer Monitor destination",
			},
			"monitor_reverse": {
				Type:        schema.TypeBool,
				Default:     true,
				Optional:    true,
				Description: "Network loadbalancer Monitor Reverse",
			},
			"monitor_transparent": {
				Type:        schema.TypeBool,
				Default:     true,
				Optional:    true,
				Description: "Network loadbalancer Monitor transparent",
			},
			"monitor_adaptive": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Network loadbalancer Monitor adaptive",
			},
			"fall_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Network loadbalancer Monitor fall counts",
			},
			"rise_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Network loadbalancer Monitor rise counts",
			},
			"alias_port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Network loadbalancer Monitor alias port",
			},
			"monitor_type": {
				Type:             schema.TypeString,
				ValidateDiagFunc: validations.StringInSlice([]string{"LBHttpMonitorProfile", "LBHttpsMonitorProfile", "LBIcmpMonitorProfile", "LBPassiveMonitorProfile", "LBTcpMonitorProfile", "LBUdpMonitorProfile"}, false),
				Required:         true,
				Description:      "Network Loadbalancer Supported values are `LBHttpMonitorProfile`,`LBHttpsMonitorProfile`, `LBIcmpMonitorProfile`, `LBPassiveMonitorProfile`,`LBTcpMonitorProfile`, `LBUdpMonitorProfile`",
			},
		},
		ReadContext:   loadbalancerMonitorReadContext,
		UpdateContext: loadbalancerMonitorReadContext,
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
	if err := c.CmpClient.LoadBalancerMonitor.Read(ctx, data, meta); err != nil {
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
	if err := c.CmpClient.LoadBalancerMonitor.Create(ctx, data, meta); err != nil {
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
	if err := c.CmpClient.LoadBalancerMonitor.Delete(ctx, data, meta); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
