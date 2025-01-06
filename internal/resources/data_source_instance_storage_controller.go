// (C) Copyright 2024 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"

	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/utils"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ReadInstanceStorageController() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"layout_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The layout ID of an instance.",
			},
			"controller_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The controller name displayed in an instance storage controller section.",
			},
			"bus_number": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The Bus sequence for a storage controller type",
			},
			"interface_number": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The interface number to be allocated. Enter `0` to automatically pick the next available interface number.",
			},
		},
		ReadContext: readInstanceStorageControllerContext,
		Description: `The ` + DSInstanceStorageController + ` data source can be used to discover the ID of a storage controller mount.
		This can then be used with resources or data sources that require a ` + DSInstanceStorageType + `,
		such as the ` + ResInstance + ` resource.`,
		SchemaVersion:  0,
		StateUpgraders: nil,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func readInstanceStorageControllerContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(d)
	err = c.CmpClient.InstanceStorageController.Read(ctx, data, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
