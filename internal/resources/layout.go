// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hpe-hcss/vmaas-terraform-resources/internal/utils"
	"github.com/hpe-hcss/vmaas-terraform-resources/pkg/client"
)

func LayoutData() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: f(generalNamedesc, "layout", "layout"),
			},
			"instance_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Type for the instance. This should be vmware for vmaas resource.`,
			},
			"instance_code": {
				Type:     schema.TypeString,
				Computed: true,
				Description: "Unique code used to identify the instance type. " +
					"Instance_code can use as ID for instance type.",
			},
		},
		ReadContext: layoutReadContext,
		Description: `The hpegl_vmaas_layout data source can be used to discover the ID of a hpegl vmaas layout.
		This can then be used with resources or data sources that require a ` + DSLayout + `,
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

func layoutReadContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}
	data := utils.NewData(d)
	err = c.CmpClient.Layout.Read(ctx, data)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
