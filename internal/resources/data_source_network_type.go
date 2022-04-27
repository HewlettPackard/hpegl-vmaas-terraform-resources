// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"

	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/utils"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func NetworkTypeData() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the network type. Example 'NSX-T Segment' or 'Custom Network'",
			},
		},
		ReadContext:    networkTypeReadContext,
		SchemaVersion:  0,
		StateUpgraders: nil,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: `The ` + DSNetworkType + ` data source can be used to discover the ID of an hpegl vmaas network types.
		This can then be used with resources or data sources that require an ` + DSNetworkType + `,
		such as the ` + ResNetwork + ` resource.`,
	}
}

func networkTypeReadContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(d)
	err = c.CmpClient.NetworkType.Read(ctx, data, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
