// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"
	"github.com/hpe-hcss/vmaas-terraform-resources/pkg/auth"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hpe-hcss/vmaas-terraform-resources/internal/utils"
	"github.com/hpe-hcss/vmaas-terraform-resources/pkg/client"
)

func NetworkInterfaceData() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: f(generalNamedesc, "Network Interface Type", "Network Interface Type"),
			},
			"code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "code corresponding to the network interface type",
			},
			"cloud_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: f(generalDDesc, "cloud"),
			},
		},
		ReadContext: networkInterfaceReadContext,
		Description: `The ` + DSNetworkInterface + ` data source can be used to discover the ID of a hpegl vmaas NetworkInterface.
		This can then be used with resources or data sources that require a ` + DSNetworkInterface + `,
		such as the ` + ResInstance + ` resource.`,
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

func networkInterfaceReadContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	auth.SetScmClientToken(&ctx, meta)
	data := utils.NewData(d)
	err = c.CmpClient.NetworkInterface.Read(ctx, data, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
