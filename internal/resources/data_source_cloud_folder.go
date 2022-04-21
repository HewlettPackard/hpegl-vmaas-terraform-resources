// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"

	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/utils"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func CloudFolderData() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: f(generalNamedesc, "cloud folder", "cloud folder"),
			},
			"cloud_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: f(generalDDesc, "cloud"),
			},
			"code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "External ID or code for the folder.",
			},
		},
		ReadContext: cloudFolderReadContext,
		Description: `The ` + DSCloudFolder + ` data source can be used to discover the ID for a folder.
		` + DSCloudFolder + ` can be used along with hpegl_vmaas_instance, If it is used, all instances/VMs
		spawned will be stored in the specified folder.`,
		SchemaVersion:  0,
		StateUpgraders: nil,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func cloudFolderReadContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(d)
	err = c.CmpClient.CloudFolder.Read(ctx, data, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
