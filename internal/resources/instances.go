// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hpe-hcss/vmaas-terraform-resources/internal/utils"
	"github.com/hpe-hcss/vmaas-terraform-resources/pkg/client"
)

const (
	instanceCreateRetryTimeout    = 10 * time.Minute
	instanceCreateRetryDelay      = 60 * time.Second
	instanceCreateRetryMinTimeout = 30 * time.Second
	instanceUpdateRetryTimeout    = 10 * time.Minute
	instanceUpdateRetryDelay      = 15 * time.Second
	instanceUpdateRetryMinTimeout = 15 * time.Second
)

func Instances() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the instance to be provisioned.",
			},
			"cloud_id": {
				Type:        schema.TypeInt,
				ForceNew:    true,
				Required:    true,
				Description: f(generalDDesc, "cloud"),
			},
			"group_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: f(generalDDesc, "group"),
			},
			"plan_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: f(generalDDesc, "plan"),
			},
			"layout_id": {
				Type:        schema.TypeInt,
				ForceNew:    true,
				Required:    true,
				Description: f(generalDDesc, "layout"),
			},
			"instance_type_code": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "Unique code used to identify the instance type.",
			},
			"network": {
				Type:        schema.TypeSet,
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
				Required: true,
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
							Type:        schema.TypeString,
							Required:    true,
							Description: f(generalDDesc, "datastore"),
						},
						"root": {
							Type:        schema.TypeBool,
							Default:     true,
							Optional:    true,
							Description: "If true then the given volume as considered as root volume.",
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
			"port": {
				Type:        schema.TypeList,
				ForceNew:    true,
				Optional:    true,
				Description: "Provide port",
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
				Required:    true,
				Description: "Configuration details for the instance to be provisioned.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_pool_id": {
							Type:        schema.TypeInt,
							Required:    true,
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
							Default:     true,
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
							Default:     false,
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
				Default:     1,
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
			"clone": {
				Type:        schema.TypeSet,
				Optional:    true,
				ForceNew:    true,
				Description: "If Clone is provided, this instance will created from cloning an existing instance",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source_instance_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Instance ID of the source.",
						},
					},
				},
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
				Description: `Power operation for an instance. Allower values are
				'poweroff','poweron','restart' and 'suspend'. On creating an instance only 'poweron' operation is allowed.`,
				ValidateFunc: validation.StringInSlice([]string{
					utils.PowerOn, utils.PowerOff, utils.Restart, utils.Suspend,
				}, false),
			},
		},
		SchemaVersion:  0,
		StateUpgraders: nil,
		CreateContext:  instanceCreateContext,
		ReadContext:    instanceReadContext,
		UpdateContext:  instanceUpdateContext,
		DeleteContext:  instanceDeleteContext,
		CustomizeDiff:  nil,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: `Instance resource facilitates creating,
		updating and deleting virtual machines.
		For creating an instance, provide a unique name and all the Mandatory(Required) parameters.
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
		Delay:      instanceCreateRetryDelay,
		Pending:    []string{utils.StateProvisioning},
		Target:     []string{utils.StateRunning},
		Timeout:    instanceCreateRetryTimeout,
		MinTimeout: instanceCreateRetryMinTimeout,
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
	log.Print("[INFO] this a log")
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
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}
	data := utils.NewData(d)
	if err := c.CmpClient.Instance.Update(ctx, data); err != nil {
		return diag.FromErr(err)
	}
	// Wait for the status to be running
	createStateConf := resource.StateChangeConf{
		Delay:      instanceUpdateRetryDelay,
		Pending:    []string{utils.StateResizing},
		Target:     []string{utils.StateRunning, utils.StateStopped, utils.StateSuspended},
		Timeout:    instanceUpdateRetryTimeout,
		MinTimeout: instanceUpdateRetryMinTimeout,
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

	return instanceReadContext(ctx, d, meta)
}
