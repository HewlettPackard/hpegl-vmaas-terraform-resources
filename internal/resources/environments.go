// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hpe-hcss/vmaas-terraform-resources/internal/utils"
	"github.com/hpe-hcss/vmaas-terraform-resources/pkg/client"
)

func EnvironmentData() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: f(generalNamedesc, "Environment", "Environment"),
			},
			"code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "code of each environment",
			},
		},
		ReadContext: environmentReadContext,
		Description: `The hpegl_vmaas_environment data source can be used to discover the ID/Code of a hpegl vmaas environment.
		This can then be used with resources or data sources that require a hpegl_vmaas_environment,
		such as the hpegl_vmaas_instance resources etc.`,
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

func environmentReadContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}
	data := utils.NewData(d)
	err = c.CmpClient.Environment.Read(ctx, data)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
