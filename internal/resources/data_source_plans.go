// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hpe-hcss/vmaas-terraform-resources/internal/utils"
	"github.com/hpe-hcss/vmaas-terraform-resources/pkg/client"
)

func PlanData() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: f(generalNamedesc, "Plan", "Plan"),
			},
		},
		ReadContext: planReadContext,
		Description: `The ` + DSPlan + ` data source can be used to discover the ID of a hpegl vmaas plan.
		This can then be used with resources or data sources that require a ` + DSPlan + `,
		such as the ` + ResInstance + ` resources etc.`,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(readTimeout),
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

	client.SetScmClientToken(&ctx, meta)
	data := utils.NewData(d)
	err = c.CmpClient.Plan.Read(ctx, data)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
