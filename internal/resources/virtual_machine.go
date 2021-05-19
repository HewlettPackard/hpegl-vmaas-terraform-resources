// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hpe-hcss/vmaas-terraform-resources/internal/utils"
	"github.com/hpe-hcss/vmaas-terraform-resources/pkg/client"
)

const (
	vmAvailableTimeout = 60 * time.Minute
	vmDeleteTimeout    = 60 * time.Minute
)

func VirtualMachine() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the instance",
			},
			"cloud_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
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
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Required: true,
			},
			"volumes": {
				Type: schema.TypeList,
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
				Required: true,
			},
			"labels": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"tags": utils.ListOfMap(),
			"config": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vmware_resource_pool": {
							Type:     schema.TypeString,
							Required: true,
						},
						"public_key": {
							Type: schema.TypeString,
						},
					},
				},
			},
			"copies": {
				Type:    schema.TypeInt,
				Default: 1,
			},
			"evars": utils.ListOfMap(),
		},
		SchemaVersion:  0,
		StateUpgraders: nil,
		CreateContext:  vmCreateContext,
		ReadContext:    vmReadContext,
		// TODO figure out if a VM can be updated
		// Update:             vmUpdate,
		DeleteContext: vmDeleteContext,
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

func vmCreateContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}
	token := c.IAMToken
	url := c.VMaaSAPIUrl

	println(" Create Context IAM Token : " + token + " URL : " + url)

	diags := new(diag.Diagnostics)

	if c.IAMToken == "" {
		*diags = append(*diags, diag.Errorf("Empty token")...)
	}
	// instanceCreateOpts := models.CreateInstanceBodyInstance{}
	// cmp_client.APIClient{}.InstancesApi.CreateAnInstance(ctx, sid, instanceCreateOpts)
	d.SetId(string(1))

	return vmReadContext(ctx, d, meta)
}

func vmReadContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

func vmDeleteContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
