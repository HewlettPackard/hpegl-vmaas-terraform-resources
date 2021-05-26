// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hpe-hcss/vmaas-terraform-resources/internal/utils"
	"github.com/hpe-hcss/vmaas-terraform-resources/pkg/client"
)

func DatastoreData() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				Description: `Name of the datastore. This needs to be exact name or
				else will return error not found`,
			},
			"cloud_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID for cloud or zone",
			},
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID for group",
			},
			"resource_pool_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Resource pool ID",
			},
		},
		ReadContext: datastoreReadContext,
		Description: "Get the datastore details",
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

func datastoreReadContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}
	data := utils.NewData(d)
	err = c.CmpClient.Datastore.Read(ctx, data)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
