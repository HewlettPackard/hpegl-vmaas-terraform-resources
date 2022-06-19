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

func LoadBalancerVirtualServers() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"vip_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "vip_name of Network loadbalancer virtual server name",
			},
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "description of Network loadbalancer virtual server",
				ForceNew:    true,
			},
			"vip_address": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "vip_address of Network loadbalancer virtual server",
				ForceNew:    true,
			},
			"vip_port": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "vip_port of network loadbalancer virtual server",
			},
			"pool": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "pool of Network loadbalancer virtual server",
			},
			"ssl_cert": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "ssl_cert of Network loadbalancer virtual server",
			},
			"ssl_server_cert": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "ssl_server_cert of the Network loadbalancer virtual server",
			},
			"config": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "virtual server Configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"persistence": {
							Type:             schema.TypeString,
							ValidateDiagFunc: validations.StringInSlice([]string{"SOURCE_IP", "COOKIE", "DISBALED"}, false),
							Required:         true,
							Description:      "Network Loadbalancer Supported values are `SOURCE_IP`,`COOKIE`, `DISBALED`"},
						"persistence_profile": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "persistence_profile of virtual server Configuration",
						},
						"application_profile": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "application_profile of virtual server Configuration",
						},
						"ssl_client_profile": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ssl_client_profile of virtual server Configuration",
						},
						"ssl_server_profile": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ssl_server_profile of virtual server Configuration",
						},
					},
				},
			},
		},
		ReadContext:   loadbalancerVirtualServerReadContext,
		UpdateContext: loadbalancerVirtualServerReadContext,
		CreateContext: loadbalancerVirtualServerCreateContext,
		DeleteContext: loadbalancerVirtualServerDeleteContext,
		Description: `loadbalancer Virtual Server resource facilitates creating,
		and deleting NSX-T  Network Load Balancers.`,
	}
}

func loadbalancerVirtualServerReadContext(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(rd)
	if err := c.CmpClient.LoadBalancerVirtualServer.Read(ctx, data, meta); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func loadbalancerVirtualServerCreateContext(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(rd)
	if err := c.CmpClient.LoadBalancerVirtualServer.Create(ctx, data, meta); err != nil {
		return diag.FromErr(err)
	}

	return loadbalancerVirtualServerReadContext(ctx, rd, meta)
}

func loadbalancerVirtualServerDeleteContext(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(rd)
	if err := c.CmpClient.LoadBalancerVirtualServer.Delete(ctx, data, meta); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
