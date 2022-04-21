// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"

	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/cmp"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func InstancesClone() *schema.Resource {
	instanceCloneSchema := getInstanceDefaultSchema(true)

	instanceCloneSchema.Schema["source_instance_id"] = &schema.Schema{
		Type:     schema.TypeInt,
		Required: true,
		ForceNew: true,
		Description: `Instance ID of the source instance. For getting source instance ID
		use 'hpeg_vmaas_instance' resource.`,
	}
	instanceCloneSchema.Description = `Instance clone resource facilitates creating,
	updating and deleting cloned virtual machines.
	For creating an instance clone, provide a unique name and all the Mandatory(Required) parameters.
	All optional parameters will be inherited from parent resource if not provided.`

	instanceCloneSchema.CreateWithoutTimeout = instanceCloneCreateContext
	instanceCloneSchema.ReadWithoutTimeout = instanceCloneReadContext
	instanceCloneSchema.UpdateWithoutTimeout = instanceCloneUpdateContext
	instanceCloneSchema.DeleteWithoutTimeout = instanceCloneDeleteContext
	instanceCloneSchema.CustomizeDiff = instanceCustomizeDiff

	return instanceCloneSchema
}

type instanceCloneResourceObj struct{}

func (*instanceCloneResourceObj) getClient(c *client.Client) cmp.Resource {
	return c.CmpClient.InstanceClone
}

func instanceCloneCreateContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return instanceHelperCreateContext(ctx, &instanceCloneResourceObj{}, d, meta)
}

func instanceCloneReadContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return instanceHelperReadContext(ctx, &instanceCloneResourceObj{}, d, meta)
}

func instanceCloneDeleteContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return instanceHelperDeleteContext(ctx, &instanceCloneResourceObj{}, d, meta)
}

func instanceCloneUpdateContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return instanceHelperUpdateContext(ctx, &instanceCloneResourceObj{}, d, meta)
}
