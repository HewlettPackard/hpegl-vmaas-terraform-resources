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
			"valid_till": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Unix timestamp of when the access_token expires, in seconds",
				Sensitive:   false,
			},
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Morpheus URL",
				Sensitive:   false,
			},
		},
		ReadContext: MorpheusDetailsBrokerReadContext,
		Description: `The ` + DSMorpheusDataSource + ` data source can be used to get a details of the Morpheus instance
		used by VMaaS.  The details that can be retrieved are the access_token, valid_till (the Unix timestamp of
		access_token expiration) and the URL of the Morpheus instance.`,
		SchemaVersion:  0,
		StateUpgraders: nil,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func MorpheusDetailsBrokerReadContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
