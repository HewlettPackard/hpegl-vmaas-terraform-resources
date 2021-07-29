// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hpe-hcss/vmaas-terraform-resources/internal/cmp"
	"github.com/hpe-hcss/vmaas-terraform-resources/internal/utils"
	"github.com/hpe-hcss/vmaas-terraform-resources/pkg/client"
)

const (
	// update
	instanceUpdateRetryTimeout    = 10 * time.Minute
	instanceUpdateRetryDelay      = 15 * time.Second
	instanceUpdateRetryMinTimeout = 15 * time.Second
)

type resourceObject interface {
	getClient(*client.Client) cmp.Resource
}

func getInstanceDefaultSchema(isClone bool) *schema.Resource {
	layoutID := &schema.Schema{
		Type:        schema.TypeInt,
		Description: f(generalDDesc, "layout"),
	}
	if isClone {
		layoutID.Computed = true
		layoutID.Optional = true
	} else {
		layoutID.Required = true
	}

	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the instance to be provisioned.",
			},
			"cloud_id": {
				Type:        schema.TypeInt,
				Optional:    isClone,
				Required:    !isClone,
				ForceNew:    true,
				Description: f(generalDDesc, "cloud"),
			},
			"group_id": {
				Type:        schema.TypeInt,
				Optional:    isClone,
				Required:    !isClone,
				ForceNew:    true,
				Description: f(generalDDesc, "group"),
			},
			"layout_id": layoutID,
			"plan_id": {
				Type:        schema.TypeInt,
				Optional:    isClone,
				Required:    !isClone,
				ForceNew:    true,
				Description: f(generalDDesc, "plan"),
			},
			"instance_type_code": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    isClone,
				Required:    !isClone,
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
				Optional: isClone,
				Required: !isClone,
				MinItems: 1,
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
						"root": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "true if volume is root",
							Optional:    true,
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
				Computed:true,
				Optional:    true,
				Description: "Hostname for the instance",
			},
			"ip": {
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				Description: "IP assigned to instance",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"config": {
				Type:        schema.TypeSet,
				ForceNew:    true,
				Optional:    isClone,
				Required:    !isClone,
				Description: "Configuration details for the instance to be provisioned.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_pool_id": {
							Type:        schema.TypeInt,
							Optional:    isClone,
							Required:    !isClone,
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
				'poweroff', 'poweron' and 'suspend'. Upon creating an instance only 'poweron' operation is allowed.`,
				ValidateFunc: validation.StringInSlice([]string{
					utils.PowerOn, utils.PowerOff, utils.Suspend,
				}, false),
			},
			"restart_instance": {
				Type:     schema.TypeInt,
				Optional: true,
				Description: `Restarts the instance if set to any positive integer.
				Restart works only on pre-created instance.`,
				ValidateFunc: validation.IntAtLeast(1),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.HasChange("power")
				},
			},
			"snapshot": {
				Type:     schema.TypeSet,
				MaxItems: 1,
				Description: `Snapshot details to be created. Snapshot name and description
				 should be unique. Any change in those will results into creation of new snapshot,
				 with preserving previous snapshot(s).`,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Description: `ID of the snapshot.`,
							Optional:    true,
							Computed:    true,
						},
						"name": {
							Type:        schema.TypeString,
							Description: `Name of the snapshot.`,
							Required:    true,
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Description of the snapshot",
						},
						"is_snapshot_exists": {
							Type: schema.TypeBool,
							Description: `Flag which will be set to be true if the snapshot with the name
							exists.`,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		},
		SchemaVersion:  0,
		StateUpgraders: nil,
		CustomizeDiff:  nil,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func instanceHelperCreateContext(
	ctx context.Context,
	ro resourceObject,
	d *schema.ResourceData,
	meta interface{},
) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(d)
	if err := ro.getClient(c).Create(ctx, data, meta); err != nil {
		return diag.FromErr(err)
	}

	return instanceHelperReadContext(ctx, ro, d, meta)
}

func instanceHelperReadContext(
	ctx context.Context,
	ro resourceObject,
	d *schema.ResourceData,
	meta interface{},
) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(d)
	err = ro.getClient(c).Read(ctx, data, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func instanceHelperDeleteContext(
	ctx context.Context,
	ro resourceObject,
	d *schema.ResourceData,
	meta interface{},
) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(d)
	if err := ro.getClient(c).Delete(ctx, data, meta); err != nil {
		return diag.FromErr(err)
	}
	data.SetID("")

	return nil
}

func instanceHelperUpdateContext(
	ctx context.Context,
	ro resourceObject,
	d *schema.ResourceData,
	meta interface{},
) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(d)
	if err := ro.getClient(c).Update(ctx, data, meta); err != nil {
		return diag.FromErr(err)
	}
	// Wait for the status to be running
	updateStateConf := resource.StateChangeConf{
		Delay:      instanceUpdateRetryDelay,
		Pending:    []string{utils.StateResizing, utils.StateStopping, utils.StateSuspending, utils.StateRestarting},
		Target:     []string{utils.StateRunning, utils.StateStopped, utils.StateSuspended},
		Timeout:    instanceUpdateRetryTimeout,
		MinTimeout: instanceUpdateRetryMinTimeout,
		Refresh: func() (result interface{}, state string, err error) {
			if err := ro.getClient(c).Read(ctx, data, meta); err != nil {
				return nil, "", err
			}

			return d.Get("name"), data.GetString("status"), nil
		},
	}
	_, err = updateStateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	return instanceReadContext(ctx, d, meta)
}
