// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"

	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/utils"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourcePoolData() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: f(generalNamedesc, "Resource pool of a cluster", "Resource pool"),
			},
			"cloud_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: f(generalDDesc, "cloud"),
			},
		},
		ReadContext: resourcePoolReadContext,
		Description: `The ` + DSResourcePool + ` data source can be used to discover the ID of an hpegl vmaas resource pool.
		This can then be used with resources or data sources that require an ` + DSResourcePool + `,
		such as the ` + ResInstance + ` resource.`,
		SchemaVersion:  0,
		StateUpgraders: nil,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourcePoolReadContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(d)
	err = c.CmpClient.ResourcePool.Read(ctx, data, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
