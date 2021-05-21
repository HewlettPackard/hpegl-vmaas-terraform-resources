// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hpe-hcss/vmaas-terraform-resources/pkg/client"
)

const (
	instanceAvailableTimeout = 60 * time.Minute
	instanceReadTimeout      = 2 * time.Minute
	instanceDeleteTimeout    = 60 * time.Minute
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
			"layout_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_code": {
				Type:     schema.TypeString,
				Required: true,
			},
			"networks": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"interface_type_id": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"network_id": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
			"volumes": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
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
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		SchemaVersion:  0,
		StateUpgraders: nil,
		CreateContext:  instanceCreateContext,
		ReadContext:    instanceReadContext,
		// TODO figure out if a VM can be updated
		UpdateContext: instanceUpdateContext,
		DeleteContext: instanceDeleteContext,
		CustomizeDiff: nil,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(instanceAvailableTimeout),
			Update: schema.DefaultTimeout(instanceAvailableTimeout),
			Delete: schema.DefaultTimeout(instanceDeleteTimeout),
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
	if err := c.CmpClient.Instance.Create(ctx, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func instanceReadContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	err = c.CmpClient.Instance.Read(ctx, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
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

func instanceUpdateContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
