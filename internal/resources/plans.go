// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hpe-hcss/vmaas-terraform-resources/internal/utils"
	"github.com/hpe-hcss/vmaas-terraform-resources/pkg/client"
)

func PlanData() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				Description: `Name of the Plan. Provide appropriate name as appears on the GLC` +
					fmt.Sprintf(notFoundDesc, "plan."),
			},
		},
		ReadContext: planReadContext,
		Description: fmt.Sprintf(dsHeadingDesc, "plans for vmaas", "Administration->Plans and Pricing") +
			fmt.Sprintf(notFoundDesc, "plan"),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(readTimeout),
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
