// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hpe-hcss/vmaas-terraform-resources/internal/utils"
	"github.com/hpe-hcss/vmaas-terraform-resources/pkg/client"
)

const (
	// create
	instanceCloneCreateRetryTimeout    = 15 * time.Minute
	instanceCloneCreateRetryDelay      = 60 * time.Second
	instanceCloneCreateRetryMinTimeout = 30 * time.Second
	// update
	instanceCloneUpdateRetryTimeout    = 10 * time.Minute
	instanceCloneUpdateRetryDelay      = 15 * time.Second
	instanceCloneUpdateRetryMinTimeout = 15 * time.Second
	// delete
	instanceClonedeleteRetryDelay      = 15 * time.Second
	instanceClonedeleteRetryTimeout    = 5 * time.Minute
	instanceClonedeleteRetryMinTimeout = 15 * time.Second
)

func InstancesClone() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"source_instance_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
				Description: `Instance ID of the source instance. For getting source instance ID
				use 'hpeg_vmaas_instance' resource.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the instance to be provisioned.",
			},
			"cloud_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: f(generalDDesc, "cloud"),
			},
			"group_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: f(generalDDesc, "group"),
			},
			"layout_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: f(generalDDesc, "layout"),
			},
			"plan_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: f(generalDDesc, "plan"),
			},
			"instance_type_code": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: "Unique code used to identify the instance type.",
			},
			"network": {
				Type:        schema.TypeList,
				ForceNew:    true,
				Required:    true,
				MinItems:    1,
				Description: "Details of the network to which the instance should belong.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: f(generalDDesc, "network ID"),
						},
						"interface_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: f(generalDDesc, "network interface type"),
						},
					},
				},
			},
			"volume": {
				Type:     schema.TypeList,
				Optional: true,
				Description: `A list of volumes to be created inside a provisioned instance.
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
							Type:     schema.TypeString,
							Required: true,
							Description: `Datastore ID can be obtained from hpegl_vmaas_datastore
							data source. Please provide 'auto' as value to select datastore as auto.`,
							DiffSuppressFunc: utils.SkipField(),
						},
						"id": {
							Computed:    true,
							Type:        schema.TypeInt,
							Description: "ID for the volume",
						},
					},
				},
			},
			"labels": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "An array of strings used for labelling instance.",
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
			"hostname": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: "Hostname for the instance",
			},
			"config": {
				Type:        schema.TypeSet,
				ForceNew:    true,
				Optional:    true,
				Description: "Configuration details for the instance to be provisioned.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_pool_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: f(generalDDesc, "resource pool"),
						},
						"template_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Unique ID for the template",
						},
						"no_agent": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "If true agent will not be installed on the instance.",
						},
						"vm_folder": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Folder name where will be stored.",
						},
						"create_user": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "If true new user will be created",
						},
						"asset_tag": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Asset tag",
						},
					},
				},
			},
			"scale": {
				Type:        schema.TypeInt,
				ForceNew:    true,
				Optional:    true,
				Description: "Number of nodes within an instance.",
			},
			"evars": {
				ForceNew: true,
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Environment Variables to be added to the provisioned instance.",
			},
			"env_prefix": {
				ForceNew:    true,
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Environment prefix",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Status of the instance.`,
			},
			"power_schedule_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Scheduled power operations",
			},
			"environment_code": {
				Type: schema.TypeString,
				Description: `Environment code, which can be obtained via
				hpegl_vmaas_environment.code`,
				Optional: true,
			},
			"power": {
				Type:     schema.TypeString,
				Optional: true,
				Description: `Power operation for an instance. Power attribute can be
				use to update power state of an existing instance. Allowed power operations are
				'poweroff','poweron','restart' and 'suspend'. Upon creating an instance only 'poweron' operation is allowed.`,
				ValidateFunc: validation.StringInSlice([]string{
					utils.PowerOn, utils.PowerOff, utils.Restart, utils.Suspend,
				}, false),
			},
		},
		SchemaVersion:  0,
		StateUpgraders: nil,
		CreateContext:  instanceCloneCreateContext,
		ReadContext:    instanceCloneReadContext,
		UpdateContext:  instanceCloneUpdateContext,
		DeleteContext:  instanceCloneDeleteContext,
		CustomizeDiff:  instanceCustomizeDiff,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: `Instance resource facilitates creating,
		updating and deleting virtual machines.
		For creating an instance clone, provide a unique name and all the Mandatory(Required) parameters.
		It is recommend to use the Vmware type for provisioning.`,
	}
}

func instanceCloneCreateContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(d)
	if err := c.CmpClient.InstanceClone.Create(ctx, data, meta); err != nil {
		return diag.FromErr(err)
	}

	// Wait for the status to be running
	createStateConf := resource.StateChangeConf{
		Delay:      instanceCloneCreateRetryDelay,
		Pending:    []string{utils.StateProvisioning},
		Target:     []string{utils.StateRunning},
		Timeout:    instanceCloneCreateRetryTimeout,
		MinTimeout: instanceCloneCreateRetryMinTimeout,
		Refresh: func() (result interface{}, state string, err error) {
			if err := c.CmpClient.InstanceClone.Read(ctx, data, meta); err != nil {
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

func instanceCloneReadContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(d)
	err = c.CmpClient.InstanceClone.Read(ctx, data, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func instanceCloneDeleteContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(d)
	if err := c.CmpClient.InstanceClone.Delete(ctx, data, meta); err != nil {
		return diag.FromErr(err)
	}

	deleteStateConf := resource.StateChangeConf{
		Delay:      instanceClonedeleteRetryDelay,
		Pending:    []string{"deleting"},
		Target:     []string{"deleted", "Failed"},
		Timeout:    instanceClonedeleteRetryTimeout,
		MinTimeout: instanceClonedeleteRetryMinTimeout,
		Refresh: func() (result interface{}, state string, err error) {
			if err := c.CmpClient.InstanceClone.Read(ctx, data, meta); err != nil {
				// Check for status 404
				statusCode := utils.GetStatusCode(err)
				if statusCode == http.StatusNotFound {
					return d.Get("name"), "deleted", nil
				}

				return nil, "Failed", err
			}

			return d.Get("name"), "deleting", nil
		},
	}
	_, err = deleteStateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	data.SetID("")

	return nil
}

func instanceCloneUpdateContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(d)
	if err := c.CmpClient.InstanceClone.Update(ctx, data, meta); err != nil {
		return diag.FromErr(err)
	}
	// Wait for the status to be running
	updateStateConf := resource.StateChangeConf{
		Delay:      instanceCloneUpdateRetryDelay,
		Pending:    []string{utils.StateResizing},
		Target:     []string{utils.StateRunning, utils.StateStopped, utils.StateSuspended},
		Timeout:    instanceCloneUpdateRetryTimeout,
		MinTimeout: instanceCloneUpdateRetryMinTimeout,
		Refresh: func() (result interface{}, state string, err error) {
			if err := c.CmpClient.InstanceClone.Read(ctx, data, meta); err != nil {
				return nil, "", err
			}

			return d.Get("name"), data.GetString("status"), nil
		},
	}
	_, err = updateStateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	return instanceCloneReadContext(ctx, d, meta)
}
