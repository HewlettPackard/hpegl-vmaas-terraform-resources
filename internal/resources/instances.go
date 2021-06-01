// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hpe-hcss/vmaas-terraform-resources/internal/utils"
	"github.com/hpe-hcss/vmaas-terraform-resources/pkg/client"
)

const (
	instanceAvailableTimeout = 60 * time.Minute
	instanceReadTimeout      = 2 * time.Minute
	instanceDeleteTimeout    = 60 * time.Minute
	instanceRetryTimeout     = 10 * time.Minute
	instanceRetryDelay       = 60 * time.Second
	instanceRetryMinTimeout  = 30 * time.Second
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
				Type:        schema.TypeInt,
				Required:    true,
				Description: "ID for cloud or zone",
			},
			"group_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "ID for group",
			},
			"plan_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"layout_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"instance_code": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"networks": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
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
							Type:     schema.TypeInt,
							Required: true,
						},
						"datastore_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"root": {
							Type:        schema.TypeBool,
							Default:     true,
							Optional:    true,
							Description: "If true then the given volume as considered as root volume.",
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
						"resource_pool_id": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"public_key": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"template": {
							Type:     schema.TypeString,
							Required: true,
						},
						"no_agent": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "If true agent will not be installed on the instance",
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
			"status": {
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
			Read:   schema.DefaultTimeout(instanceReadTimeout),
		},
		Description: "Create/update/delete instance",
	}
}

func instanceCreateContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(d)
	if err := c.CmpClient.Instance.Create(ctx, data); err != nil {
		return diag.FromErr(err)
	}

	// Wait for the status to be running
	createStateConf := resource.StateChangeConf{
		Delay:      instanceRetryDelay,
		Pending:    []string{"provisioning"},
		Target:     []string{"running"},
		Timeout:    instanceRetryTimeout,
		MinTimeout: instanceRetryMinTimeout,
		Refresh: func() (result interface{}, state string, err error) {
			if err := c.CmpClient.Instance.Read(ctx, data); err != nil {
				return nil, "", err
			}

			return d.Get("name"), data.GetString("status"), nil
		},
	}
	_, err = createStateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func instanceReadContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(d)
	err = c.CmpClient.Instance.Read(ctx, data)
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

	data := utils.NewData(d)
	if err := c.CmpClient.Instance.Delete(ctx, data); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func instanceUpdateContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
