// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hpe-hcss/vmaas-terraform-resources/internal/utils"
	"github.com/hpe-hcss/vmaas-terraform-resources/pkg/client"
)

func ResourcePoolData() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				Description: "Provide ResourcePool name of a cluster" +
					fmt.Sprintf(notFoundDesc, "resourcepool"),
			},
			"cloud_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: cloudIDDesc,
			},
		},
		ReadContext: resourcePoolReadContext,
		Description: fmt.Sprintf(dsHeadingDesc, `resource pool of a cluster where the instance 
		should be provisioned`, "Infrastructure->Clouds->Resources"),
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

func resourcePoolReadContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}
	data := utils.NewData(d)
	err = c.CmpClient.ResourcePool.Read(ctx, data)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
