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

func LoadBalancerMonitor() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"lb_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Parent lb ID, lb_id can be obtained by using LB datasource/resource.",
				ForceNew:    true,
			},
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
			"type": {
				Type:             schema.TypeString,
				ValidateDiagFunc: validations.StringInSlice([]string{"LBHttpMonitorProfile", "LBHttpsMonitorProfile", "LBIcmpMonitorProfile", "LBPassiveMonitorProfile", "LBTcpMonitorProfile", "LBUdpMonitorProfile"}, false),
				Required:         true,
				Description:      "Network Loadbalancer Supported values are `LBHttpMonitorProfile`,`LBHttpsMonitorProfile`, `LBIcmpMonitorProfile`, `LBPassiveMonitorProfile`,`LBTcpMonitorProfile`, `LBUdpMonitorProfile`",
			},
			"fall_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     3,
				Description: "Network loadbalancer Monitor fall counts",
			},
			"interval": {
				Type:        schema.TypeInt,
				Default:     5,
				Description: "Interval time for Network loadbalancer Monitor",
				Optional:    true,
			},
			"monitor_port": {
				Type:        schema.TypeInt,
				Description: "Interval time for Network loadbalancer Monitor",
				Optional:    true,
			},
			"rise_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     3,
				Description: "Network loadbalancer Monitor rise counts",
			},
			"timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     15,
				Description: "Timeout for Network loadbalancer Monitor",
			},
			"request_body": {
				Type:        schema.TypeString,
				Description: "request body to send the monitor details",
				Optional:    true,
			},
			"request_method": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateDiagFunc: validations.StringInSlice([]string{"GET", "POST", "OPTIONS",
					"HEAD", "PUT"}, false),
				Default:     "GET",
				Description: "Supported values are `GET`,`POST`,`OPTIONS`, `HEAD`,`PUT`",
			},
			"request_url": {
				Type:        schema.TypeString,
				Description: "request url to send the monitor urls",
				Optional:    true,
			},
			"request_version": {
				Type: schema.TypeString,
				ValidateDiagFunc: validations.StringInSlice([]string{"HTTP_VERSION_1_0",
					"HTTP_VERSION_1_1"}, false),
				Description: "Supported values are `HTTP_VERSION_1_0`,`HTTP_VERSION_1_1`",
				Optional:    true,
			},
			"response_data": {
				Type:        schema.TypeString,
				Description: "response data to get the monitor data",
				Optional:    true,
			},
			"response_status_codes": {
				Type:        schema.TypeString,
				Description: "response status codes for the monitor calls",
				Optional:    true,
			},
			"data_length": {
				Type:        schema.TypeInt,
				Default:     56,
				Description: "data length is for the ICMP monitor type",
				Optional:    true,
			},
			"max_fail": {
				Type:        schema.TypeInt,
				Default:     5,
				Description: "maximum failure for the ICMP monitor type",
				Optional:    true,
			},
			"send_type": {
				Type:        schema.TypeString,
				Description: "send type method like GET,POST",
				Optional:    true,
			},
			// "monitor_port": {
			// 	Type:        schema.TypeInt,
			// 	Optional:    true,
			// 	Description: "Network loadbalancer Monitor alias port",
			// },
		},
		ReadContext:   loadbalancerMonitorReadContext,
		UpdateContext: loadbalancerMonitorUpdateContext,
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

func loadbalancerMonitorUpdateContext(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(rd)
	if err := c.CmpClient.LoadBalancerMonitor.Update(ctx, data, meta); err != nil {
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
