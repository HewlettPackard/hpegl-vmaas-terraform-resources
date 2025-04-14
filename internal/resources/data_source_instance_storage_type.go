// (C) Copyright 2024-2025 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"

	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/utils"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ReadInstanceStorageType() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: f(generalNamedesc, "disk type", "disk type"),
			},
			"cloud_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: f(generalDDesc, "cloud"),
			},
			"layout_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: f(generalDDesc, "layout"),
			},
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: f(generalDDesc, "group"),
			},
		},
		ReadContext: readInstanceStorageTypeContext,
		Description: `The ` + DSInstanceStorageType + ` data source can be used to discover the ID of a disk type.
		This can then be used with resources or data sources that require a ` + DSInstanceStorageType + `,
		such as the ` + ResInstance + ` resource.`,
		SchemaVersion:  0,
		StateUpgraders: nil,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func readInstanceStorageTypeContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(d)
	err = c.CmpClient.InstanceStorageType.Read(ctx, data, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
