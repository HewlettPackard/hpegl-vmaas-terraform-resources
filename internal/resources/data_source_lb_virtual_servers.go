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

func LBVirtualServerData() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"vip_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: f(generalNamedesc, "loadBalancervirtualserver", "loadbalancervirtualserver"),
			},
			"vip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "vip_address of Network loadbalancer virtual server",
			},
			"vip_port": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "vip_port of network loadbalancer virtual server",
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
							Type:             schema.TypeString,
							ValidateDiagFunc: validations.StringInSlice([]string{"SOURCE_IP", "COOKIE", "DISBALED"}, false),
							Required:         true,
							Description:      "Network Loadbalancer Supported values are `SOURCE_IP`,`COOKIE`, `DISBALED`"},
						"persistence_profile": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "persistence_profile of virtual server Configuration",
						},
						"application_profile": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "application_profile of virtual server Configuration",
						},
						"ssl_client_profile": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ssl_client_profile of virtual server Configuration",
						},
						"ssl_server_profile": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ssl_server_profile of virtual server Configuration",
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
	err = c.CmpClient.LoadBalancer.Read(ctx, data, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
