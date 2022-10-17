// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"

	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/utils"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DhcpServerData() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: f(generalNamedesc, "DHCP Server", "DHCP Server"),
			},
			"id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "id can be obtained by using DHCP Server datasource/resource.",
			},
			"provider_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ProviderId of the DHCP Server. Use the provider_id  while creating DHCP NSX-T Segment Network",
			},
		},
		ReadContext: DHCPServerReadContext,
		Description: `The ` + DSDhcpServer + ` data source can be used to discover the ID of a hpegl vmaas DHCP server.
		This can then be used with resources or data sources that require a ` + DSDhcpServer + `,
		such as the ` + ResDhcpServer + ` resource.`,
		SchemaVersion:  0,
		StateUpgraders: nil,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func DHCPServerReadContext(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(rd)
	if err := c.CmpClient.DhcpServer.Read(ctx, data, meta); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
