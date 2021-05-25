// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hpe-hcss/vmaas-terraform-resources/internal/utils"
	"github.com/hpe-hcss/vmaas-terraform-resources/pkg/client"
)

const (
	planReadTimeout = 30 * time.Second
)

func PlanData() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				Description: `Name of the Plan. This needs to be exact name or
				else will return error not found`,
			},
			"provision_name": {
				Type:     schema.TypeString,
				Required: true,
				Description: `Name of the provision. This needs to be exact name or
				else will return error not found`,
			},
		},
		ReadContext: planReadContext,
		Description: "Get the plan details",
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(planReadTimeout),
		},
		SchemaVersion:  0,
		StateUpgraders: nil,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func planReadContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}
	data := utils.NewData(d)
	err = c.CmpClient.Plan.Read(ctx, data)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
