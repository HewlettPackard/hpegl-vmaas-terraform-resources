// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	client "github.com/hpe-hcss/vmaas-cmp-go-sdk/pkg/client"
)

// instance implements functions related to cmp instances
type instance struct {
	// expose Instance API service to instances related operations
	iClient         *client.InstancesApiService
	serviceInstance string
}

func (i *instance) Create(ctx context.Context, d *schema.ResourceData) error {
	return nil
}
func (i *instance) Update(ctx context.Context, d *schema.ResourceData) error {
	return nil
}
func (i *instance) Delete(ctx context.Context, d *schema.ResourceData) error {
	return nil
}
func (i *instance) Read(ctx context.Context, d *schema.ResourceData) error {
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	resp, err := i.iClient.GetASpecificInstance(ctx, i.serviceInstance, int32(id))
	if err != nil {
		return err
	}
	d.Set("name", resp.Instance.Name)
	d.Set("cloud_id", resp.Instance.Cloud.Id)
	d.Set("group_id", resp.Instance.Group.Id)
	d.Set("plan_id", resp.Instance.Plan.Id)
	d.Set("instance_type", resp.Instance.InstanceType)
	d.Set("networks", resp.Instance)
	d.Set("volumes", resp.Instance.Volumes)
	d.Set("size", resp.Instance.Volumes[0].Size)
	d.Set("datastore_id", resp.Instance.Volumes[0].DatastoreId)
	d.Set("labels", resp.Instance.Labels)
	d.Set("tags", resp.Instance.Tags)
	d.Set("config", resp.Instance.Config)
	d.Set("vmware_resource_pool", resp.Instance.Config.ResourcePoolID)
	// d.Set("public_key", resp.Instance.Config)
	// d.Set("copies", resp.Instance)
	d.Set("evars", resp.Instance.Evars)
	return nil
}
