// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"

	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/utils"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Tier1RouterData() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: f(generalNamedesc, "Router", "Router"),
			},
			"provider_id": {
				Type:     schema.TypeString,
				Computed: true,
				Description: "Provider ID of the given router/gateway. This field can be used as connected_gateway in " +
					ResNetwork,
			},
			"interfaces": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Interface Configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "ID of the Uplink Interface",
						},
						"source_addresses": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Interface IP Address of the Uplink Interface",
						},
						"cidr": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CIDR of the network the Uplink Interface",
						},
					},
				},
			},
		},
		ReadContext: Tier1RouterReadContext,
		Description: `The ` + DSTier1Router + ` data source can be used to discover the ID of a hpegl vmaas router.
		This can then be used with resources or data sources that require a ` + DSTier1Router + `,
		such as the ` + ResRouter + ` resource.`,
		SchemaVersion:  0,
		StateUpgraders: nil,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func Tier1RouterReadContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(d)
	err = c.CmpClient.DSTier1Router.Read(ctx, data, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
