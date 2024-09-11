// (C) Copyright 2024 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/utils"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/client"
)

// MorpheusDetailsBroker returns a schema.Resource for the MorpheusDetails data source
func MorpheusDetailsBroker() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"access_token": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Morpheus access_token",
				Sensitive:   true,
			},
			"refresh_token": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Morpheus refresh_token",
				Sensitive:   true,
			},
			"morpheus_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Morpheus URL",
				Sensitive:   false,
			},
		},
		ReadContext:    morpheusDetailsBrokerReadContext,
		Description:    `The ` + DSMorpheusDataSource + ` data source can be used to get a Morpheus token and URL using the IAM API Client creds provided`,
		SchemaVersion:  0,
		StateUpgraders: nil,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func morpheusDetailsBrokerReadContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(d)
	err = c.BrokerClient.DSMorpheusDetails.Read(ctx, data, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
