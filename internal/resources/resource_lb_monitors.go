// (C) Copyright 2022 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"

	diffvalidation "github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/resources/diffValidation"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/resources/schemas"
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
			"http_monitor":    schemas.HttpMonitorSchema(),
			"https_monitor":   schemas.HttpsMonitorSchema(),
			"icmp_monitor":    schemas.IcmpMonitorSchema(),
			"passive_monitor": schemas.PassiveMonitorSchema(),
			"tcp_monitor":     schemas.TcpMonitorSchema(),
			"udp_monitor":     schemas.UdpMonitorSchema(),
		},
		ReadContext:   loadbalancerMonitorReadContext,
		UpdateContext: loadbalancerMonitorUpdateContext,
		CreateContext: loadbalancerMonitorCreateContext,
		DeleteContext: loadbalancerMonitorDeleteContext,
		CustomizeDiff: monitorCustomDiff,
		Description: `loadbalancer Monitor resource facilitates creating,updating
		and deleting NSX-T Network Load Balancers.`,
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

func monitorCustomDiff(ctx context.Context, diff *schema.ResourceDiff, meta interface{}) error {
	return diffvalidation.NewLoadBalancerMonitorValidate(diff).DiffValidate()
}
