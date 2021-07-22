// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hpe-hcss/vmaas-terraform-resources/internal/cmp"
	"github.com/hpe-hcss/vmaas-terraform-resources/pkg/client"
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
	All optional parameters will be inherits from parent resource if not provided.`

	instanceCloneSchema.CreateContext = instanceCloneCreateContext
	instanceCloneSchema.ReadContext = instanceCloneReadContext
	instanceCloneSchema.UpdateContext = instanceCloneUpdateContext
	instanceCloneSchema.DeleteContext = instanceCloneDeleteContext

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
