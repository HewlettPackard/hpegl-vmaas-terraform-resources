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
	return getSharedInstanceSchema(true)
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
