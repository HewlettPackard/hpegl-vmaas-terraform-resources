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
	snapshotAvailableTimeout = 1 * time.Minute
	snapshotReadTimeout      = 2 * time.Minute
	// snapshotDeleteTimeout    = 20 * time.Second
	snapshotRetryTimeout    = 10 * time.Minute
	snapshotRetryDelay      = 10 * time.Second
	snapshotRetryMinTimeout = 30 * time.Second
)

func Snapshots() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeInt,
				ForceNew:    true,
				Required:    true,
				Description: "Instance ID of which VMware snapshot to be taken",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Name of the Instance Snapshot",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Description for VMware Snapshot",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "status of instance snapshot",
			},
			"timestamp": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Timestamp when instance snapshot is created",
			},
		},
		SchemaVersion:  0,
		StateUpgraders: nil,
		CreateContext:  snapshotCreateContext,
		ReadContext:    snapshotReadContext,
		DeleteContext:  snapshotDeleteContext,
		CustomizeDiff:  nil,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(snapshotAvailableTimeout),
			Read:   schema.DefaultTimeout(snapshotReadTimeout),
		},
		Description: `Snapshot resource facilitates creating,
			VMware snapshot of insatnce.For creating an VMware snapshot of instance, 
			provide a unique name and all the Mandatory(Required) parameters.`,
	}
}

func snapshotCreateContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(d)
	if err := c.CmpClient.Snapshot.Create(ctx, data); err != nil {
		return diag.FromErr(err)
	}

	// Wait for the status to be complete
	createStateConf := resource.StateChangeConf{
		Delay:      snapshotRetryDelay,
		Pending:    []string{"creating"},
		Target:     []string{"complete"},
		Timeout:    snapshotRetryTimeout,
		MinTimeout: snapshotRetryMinTimeout,
		Refresh: func() (result interface{}, state string, err error) {
			if err := c.CmpClient.Snapshot.Read(ctx, data); err != nil {
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

func snapshotReadContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(d)
	err = c.CmpClient.Snapshot.Read(ctx, data)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func snapshotDeleteContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return diag.Errorf("Delete of Snapshot is not supported from terraform")
}
