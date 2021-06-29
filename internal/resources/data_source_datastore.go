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
				Type:        schema.TypeString,
				Required:    true,
				Description: f(generalNamedesc, "datastore", "datastore"),
			},
			"cloud_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: f(generalDDesc, "cloud"),
			},
		},
		ReadContext: datastoreReadContext,
		Description: `The ` + DSDatastore + ` data source can be used to discover the ID of a hpegl vmaas datastore.
		This can then be used with resources or data sources that require a ` + DSDatastore + `,
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

func datastoreReadContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	client.SetScmClientToken(&ctx, meta)
	data := utils.NewData(d)
	err = c.CmpClient.Datastore.Read(ctx, data)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
