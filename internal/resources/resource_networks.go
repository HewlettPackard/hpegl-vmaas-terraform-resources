// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"

	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/resources/schemas"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/utils"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Network() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"static_network": schemas.StaticNetworkSchema(),
			"dhcp_network":   schemas.DhcpNetworkSchema(),
		},
		SchemaVersion: 0,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		ReadContext:   resNetworkReadContext,
		CreateContext: resNetworkCreateContext,
		UpdateContext: resNetworkUpdateContext,
		DeleteContext: resNetworkDeleteContext,
		Description: `Network resource facilitates creating,
		updating and deleting NSX-T Networks.`,
	}
}

func resNetworkReadContext(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(rd)
	err = c.CmpClient.ResNetwork.Read(ctx, data, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resNetworkCreateContext(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(rd)
	err = c.CmpClient.ResNetwork.Create(ctx, data, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	return resNetworkReadContext(ctx, rd, meta)
}

func resNetworkDeleteContext(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(rd)
	err = c.CmpClient.ResNetwork.Delete(ctx, data, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resNetworkUpdateContext(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(rd)
	err = c.CmpClient.ResNetwork.Update(ctx, data, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
