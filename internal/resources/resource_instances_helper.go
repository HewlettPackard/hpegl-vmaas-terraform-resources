// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hpe-hcss/vmaas-terraform-resources/internal/cmp"
	"github.com/hpe-hcss/vmaas-terraform-resources/internal/utils"
	"github.com/hpe-hcss/vmaas-terraform-resources/pkg/client"
)

const (
	// create
	instanceCreateRetryTimeout    = 10 * time.Minute
	instanceCreateRetryDelay      = 60 * time.Second
	instanceCreateRetryMinTimeout = 30 * time.Second
	// update
	instanceUpdateRetryTimeout    = 10 * time.Minute
	instanceUpdateRetryDelay      = 15 * time.Second
	instanceUpdateRetryMinTimeout = 15 * time.Second
	// delete
	instancedeleteRetryDelay      = 15 * time.Second
	instancedeleteRetryTimeout    = 60 * time.Second
	instancedeleteRetryMinTimeout = 15 * time.Second
)

type resourceObject interface {
	getResourceObject(*client.Client) cmp.Resource
}

func instanceHelperCreateContext(ctx context.Context, ro resourceObject, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(d)
	if err := ro.getResourceObject(c).Create(ctx, data, meta); err != nil {
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
			if err := ro.getResourceObject(c).Read(ctx, data, meta); err != nil {
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

func instanceHelperReadContext(ctx context.Context, ro resourceObject, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(d)
	err = ro.getResourceObject(c).Read(ctx, data, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func instanceHelperDeleteContext(ctx context.Context, ro resourceObject, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(d)
	if err := ro.getResourceObject(c).Delete(ctx, data, meta); err != nil {
		return diag.FromErr(err)
	}

	deleteStateConf := resource.StateChangeConf{
		Delay:      instancedeleteRetryDelay,
		Pending:    []string{"deleting"},
		Target:     []string{"deleted", "Failed"},
		Timeout:    instancedeleteRetryTimeout,
		MinTimeout: instancedeleteRetryMinTimeout,
		Refresh: func() (result interface{}, state string, err error) {
			if err := ro.getResourceObject(c).Read(ctx, data, meta); err != nil {
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

func instanceHelperUpdateContext(ctx context.Context, ro resourceObject, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(d)
	if err := ro.getResourceObject(c).Update(ctx, data, meta); err != nil {
		return diag.FromErr(err)
	}
	// Wait for the status to be running
	updateStateConf := resource.StateChangeConf{
		Delay:      instanceUpdateRetryDelay,
		Pending:    []string{utils.StateResizing},
		Target:     []string{utils.StateRunning, utils.StateStopped, utils.StateSuspended},
		Timeout:    instanceUpdateRetryTimeout,
		MinTimeout: instanceUpdateRetryMinTimeout,
		Refresh: func() (result interface{}, state string, err error) {
			if err := ro.getResourceObject(c).Read(ctx, data, meta); err != nil {
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
