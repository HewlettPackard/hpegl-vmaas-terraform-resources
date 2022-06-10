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
			"vipName": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "vipName of Network loadbalancer virtual server name",
			},
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "description of Network loadbalancer virtual server",
				ForceNew:    true,
			},
			"vipAddress": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "vipAddress of Network loadbalancer virtual server",
				ForceNew:    true,
			},
			"vipPort": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "vipPort of network loadbalancer virtual server",
				Default:     true,
			},
			"vipProtocol": {
				Type:             schema.TypeString,
				ValidateDiagFunc: validations.StringInSlice([]string{"http", "tcp", "udp"}, false),
				Description:      "Network Loadbalancer Supported values are `http`,`tcp`, `udp`"},
			"pool": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "pool of Network loadbalancer virtual server",
			},
			"sslCert": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "sslCert of Network loadbalancer virtual server",
			},
			"sslServerCert": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "sslServerCert of the Network loadbalancer virtual server",
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
						"persistenceProfile": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "persistenceProfile of virtual server Configuration",
						},
						"applicationProfile": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "applicationProfile of virtual server Configuration",
						},
						"sslClientProfile": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "sslClientProfile of virtual server Configuration",
						},
						"sslServerProfile": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "sslServerProfile of virtual server Configuration",
						},
					},
				},
			},
		},
		ReadContext:   loadbalancerVirtualServerReadContext,
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
