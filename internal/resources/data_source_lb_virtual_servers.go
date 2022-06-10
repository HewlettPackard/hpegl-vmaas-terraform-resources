// (C) Copyright 2022 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"

	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/utils"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func LBVirtualServerData() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"vipName": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: f(generalNamedesc, "ResLoadBalancerVirtualServer", "ResLoadBalancerVirtualServer"),
			},
			"vipAddress": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "vipAddress of Network loadbalancer virtual server",
			},
			"vipPort": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "vipPort of network loadbalancer virtual server",
			},
			"vipProtocol": {
				Type: schema.TypeString,
				            ValidateDiagFunc: validations.StringInSlice([]string{
					            "http", "tcp", "udp"
				            }, false),
				            Description: "Network Loadbalancer Supported values are `http`, 
				                `tcp`, `udp`",
			},
			"pool": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "pool of Network loadbalancer virtual server",
			},
			"config": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "virtual server Configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"persistence": {
							Type: schema.TypeString,
				            ValidateDiagFunc: validations.StringInSlice([]string{
					            "SOURCE_IP", "COOKIE", "DISBALED"
				            }, false),
				            Computed:    true,
				            Description: "Network Loadbalancer Supported values are `SOURCE_IP`, 
				                `COOKIE`, `DISBALED`",
						},
						"persistenceProfile": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "persistenceProfile of virtual server Configuration",
						},
						"applicationProfile": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "applicationProfile of virtual server Configuration",
						},
						"sslClientProfile": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "sslClientProfile of virtual server Configuration",
						},
						"sslServerProfile": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "sslServerProfile of virtual server Configuration",
						},
					},
				},
			},
		},
		ReadContext: LBVirtualServerReadContext,
		Description: `The ` + DSLBVirtualServer + ` virtual server data source can be used to discover the ID of a hpegl vmaas router.
		This can then be used with resources or data sources that require a ` + DSLBVirtualServer + `,
		virtual server such as the ` + ResLoadBalancerVirtualServers + ` resource.`,
		SchemaVersion:  0,
		StateUpgraders: nil,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func LBVirtualServerReadContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(d)
	err = c.CmpClient.DSLBVirtualServer.Read(ctx, data, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
