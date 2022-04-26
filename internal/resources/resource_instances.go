// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"

	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/cmp"
	diffvalidation "github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/resources/diffValidation"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func Instances() *schema.Resource {
	instanceSchema := getInstanceDefaultSchema(false)
	instanceSchema.Schema["port"] = &schema.Schema{
		Type:        schema.TypeList,
		ForceNew:    true,
		Optional:    true,
		Description: "Provide port for the instance",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Name of the port",
				},
				"port": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Port value in string",
				},
				"lb": {
					Type:     schema.TypeString,
					Required: true,
					Description: `Load balancing configuration for ports.
					 Supported values are "No LB", "HTTP", "HTTPS", "TCP"`,
					ValidateFunc: validation.StringInSlice([]string{
						"No LB", "HTTP", "HTTPS", "TCP",
					}, false),
				},
			},
		},
	}
	instanceSchema.Description = `This Instance resource facilitates creating,
		updating and deleting virtual machines. HPE recommends that you use the VMware as type for provisioning.`
	instanceSchema.CreateWithoutTimeout = instanceCreateContext
	instanceSchema.ReadWithoutTimeout = instanceReadContext
	instanceSchema.DeleteWithoutTimeout = instanceDeleteContext
	instanceSchema.UpdateWithoutTimeout = instanceUpdateContext
	instanceSchema.CustomizeDiff = instanceCustomizeDiff

	return instanceSchema
}

type instanceResourceObj struct{}

func instanceCustomizeDiff(ctx context.Context, diff *schema.ResourceDiff, meta interface{}) error {
	instance := diffvalidation.NewInstanceValidate(diff)

	return instance.DiffValidate()
}

func (i *instanceResourceObj) getClient(c *client.Client) cmp.Resource {
	return c.CmpClient.Instance
}

func instanceCreateContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return instanceHelperCreateContext(ctx, &instanceResourceObj{}, d, meta)
}

func instanceReadContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return instanceHelperReadContext(ctx, &instanceResourceObj{}, d, meta)
}

func instanceDeleteContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return instanceHelperDeleteContext(ctx, &instanceResourceObj{}, d, meta)
}

func instanceUpdateContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return instanceHelperUpdateContext(ctx, &instanceResourceObj{}, d, meta)
}
