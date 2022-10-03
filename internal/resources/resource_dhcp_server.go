// (C) Copyright 2022 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"

	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/utils"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DhcpServer() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the DHCP server",
			},
			"server_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Server ID of Network DhcpServer",
			},
			"lease_time": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Lease time for the DHCP server",
			},
			"server_address": {
				Type:        schema.TypeString,
				Description: "server address for the DHCP server",
				Optional:    true,
			},
			"config": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "DHCP Server Configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"edge_cluster": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Edge cluster for the DHCP server",
						},
					},
				},
			},
		},
		SchemaVersion: 0,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		ReadContext:   DhcpServerReadContext,
		UpdateContext: DhcpServerUpdateContext,
		CreateContext: DhcpServerCreateContext,
		DeleteContext: DhcpServerDeleteContext,
		Description: `DhcpServer resource facilitates creating, updating
		and deleting DhcpServer.`,
	}
}

func DhcpServerReadContext(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

func DhcpServerCreateContext(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(rd)
	if err := c.CmpClient.DhcpServer.Create(ctx, data, meta); err != nil {
		return diag.FromErr(err)
	}

	return DhcpServerReadContext(ctx, rd, meta)
}

func DhcpServerUpdateContext(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(rd)
	if err := c.CmpClient.DhcpServer.Update(ctx, data, meta); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func DhcpServerDeleteContext(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(rd)
	if err := c.CmpClient.DhcpServer.Delete(ctx, data, meta); err != nil {
		return diag.FromErr(err)
	}

	return nil
}