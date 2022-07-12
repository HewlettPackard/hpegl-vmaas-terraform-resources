// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"

	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/utils"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ActiveMonitorData() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: f(generalNamedesc, "ActiveMonitor", "ActiveMonitor"),
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: f(generalNamedesc, "MonitorType", "MonitorType"),
			},
			"lb_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Parent lb ID, lb_id can be obtained by using LB datasource/resource.",
			},
		},
		ReadContext: ActiveMonitorReadContext,
		Description: `The ` + DSActiveMonitor + ` data source can be used to discover the ID of a hpegl vmaas network load balancer.
		This can then be used with resources or data sources that require a ` + DSActiveMonitor + `,
		such as the ` + ResActiveMonitor + ` resource.`,
		SchemaVersion:  0,
		StateUpgraders: nil,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func ActiveMonitorReadContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(d)
	err = c.CmpClient.DSActiveMonitor.Read(ctx, data, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
