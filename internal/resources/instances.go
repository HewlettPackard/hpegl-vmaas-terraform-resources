// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hpe-hcss/vmaas-cmp-go-sdk/pkg/models"
	"github.com/hpe-hcss/vmaas-terraform-resources/pkg/client"
)

const (
	vmAvailableTimeout = 60 * time.Minute
	vmDeleteTimeout    = 60 * time.Minute
)

func Instances() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the instance",
			},
			"cloud_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID for cloud or zone",
			},
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID for group",
			},
			"plan_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"networks": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"volumes": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size": {
							Type:     schema.TypeString,
							Required: true,
						},
						"datastore_id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"labels": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"config": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vmware_resource_pool": {
							Type:     schema.TypeString,
							Required: true,
						},
						"public_key": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"copies": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"evars": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		SchemaVersion:  0,
		StateUpgraders: nil,
		CreateContext:  instanceCreateContext,
		ReadContext:    instanceReadContext,
		// TODO figure out if a VM can be updated
		UpdateContext: instanceUpdate,
		DeleteContext: instanceDeleteContext,
		CustomizeDiff: nil,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		DeprecationMessage: "",
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(vmAvailableTimeout),
			// Update: schema.DefaultTimeout(vmAvailableTimeout),
			Delete: schema.DefaultTimeout(vmDeleteTimeout),
		},
		Description: "Create/update/delete instance",
	}
}

func instanceCreateContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	if c.IAMToken == "" {
		return diag.Errorf("Empty token")
	}
	if err := c.CmpClient.CreateInstance(models.CreateInstanceBody{}); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("1")

	return instanceReadContext(ctx, d, meta)
}

func instanceReadContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}
	token := c.IAMToken

	println(" Read Context IAM Token : " + token)

	var diags diag.Diagnostics
	id := d.Id()
	println(" ID : " + id)

	if token == "" {
		diags = append(diags, diag.Errorf("Empty token")...)
	}

	return diags
}

func instanceDeleteContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}
	token := c.IAMToken
	print(" Delete IAM Token : " + token)
	var diags diag.Diagnostics
	id := d.Id()

	if id == "" {
		diags = append(diags, diag.Errorf("Empty ID")...)
	}
	d.SetId("")

	return diags
}

func instanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
