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
	instanceRetryDelay       = 120 * time.Second
	instanceRetryMinTimeout  = 30 * time.Second
)

func Instances() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the instance to be provisioned",
			},
			"cloud_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: cloudIDDesc,
			},
			"group_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: groupIDDesc,
			},
			"plan_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: planIDDesc,
			},
			"layout_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: layoutIDDesc,
			},
			"instance_code": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique code used to identify the instance type",
			},
			"instance_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type for the instance type. This should be vmware for vmaas resource",
			},
			"networks": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Network details of which network the instance should belong to",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: networkIDDesc,
						},
					},
				},
			},
			"volumes": {
				Type:     schema.TypeList,
				Required: true,
				Description: `A list of volumes which consist of the volumes to be created inside a provisioned instance.
				It can have a root volume and other secondary volumes.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Unique name for the volume.",
						},
						"size": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Size of the volume in GB.",
						},
						"datastore_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Unique ID for the datastore.",
						},
					},
				},
			},
			"labels": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A list of strings used for labelling instances.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "A list of key and value pairs used to tag instances of similar type.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"config": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "Configuration details for the instance to be provisioned'",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_pool_id": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Unique ID of resource pool.",
						},
						"public_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ID of a public key to add to the instance.",
						},
						"template": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Unique ID of virtual image to be used.",
						},
					},
				},
			},
			"copies": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "Number of instance copies to be provisioned.",
			},
			"evars": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Environment Variables to be added to the provisioned instance.",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "State of the instance provisioned.",
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `Status of the instance .It can be one among these:
				 Provisioning/Failed/Unknown/Running.`,
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
		Description: `Instance resource facilitates creating,
		updating and deleting virtual machines.
		For creating an instance, provide a unique name and all the Mandatory(Required) parameters.,
		It is recommend to use the Vmware type for provisioning.`,
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
