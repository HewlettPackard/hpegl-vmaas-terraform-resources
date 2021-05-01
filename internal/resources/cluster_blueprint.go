// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hpe-hcss/vmaas-terraform-resources/pkg/client"
)

func ClusterBlueprint() *schema.Resource {
	return &schema.Resource{
		Schema:         nil,
		SchemaVersion:  0,
		StateUpgraders: nil,
		CreateContext:  clusterBlueprintCreateContext,
		ReadContext:    clusterBlueprintReadContext,
		// TODO figure out if and how a blueprint can be updated
		// Update:             clusterBlueprintUpdate,
		DeleteContext:      clusterBlueprintDeleteContext,
		CustomizeDiff:      nil,
		Importer:           nil,
		DeprecationMessage: "",
		Timeouts:           nil,
		Description:        "",
	}
}

func clusterBlueprintCreateContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	_ = client.GetClientFromMetaMap(meta)

	return nil
}

func clusterBlueprintReadContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	_ = client.GetClientFromMetaMap(meta)

	return nil
}

func clusterBlueprintDeleteContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	_ = client.GetClientFromMetaMap(meta)

	return nil
}
