// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hpe-hcss/vmaas-cmp-go-sdk/pkg/client"
	"github.com/hpe-hcss/vmaas-cmp-go-sdk/pkg/models"
	"github.com/hpe-hcss/vmaas-terraform-resources/internal/utils"
)

// instance implements functions related to cmp instances
type instance struct {
	// expose Instance API service to instances related operations
	iClient         *client.InstancesApiService
	serviceInstance string
}

// Create instance
func (i *instance) Create(ctx context.Context, d *schema.ResourceData) error {
	groupID, err := utils.ParseInt(d.Get("group_id").(string))
	if err != nil {
		return err
	}

	networks, err := getNetwork(d.Get("networks"))
	if err != nil {
		return err
	}

	volumes, err := getVolume(d.Get("volumes"))
	if err != nil {
		return err
	}

	config, err := getConfig(d.Get("config"))
	if err != nil {
		return err
	}
	tags, _ := getTags(d.Get("tags"))
	req := &models.CreateInstanceBody{
		ZoneId: utils.JSONNumber(d.Get("cloud_id")),
		Instance: &models.CreateInstanceBodyInstance{
			Name: d.Get("name").(string),
			InstanceType: &models.CreateInstanceBodyInstanceInstanceType{
				Code: d.Get("instance_code").(string),
			},
			Plan: &models.CreateInstanceBodyInstancePlan{
				Id: utils.JSONNumber(d.Get("plan_id")),
			},
			Site: &models.CreateInstanceBodyInstanceSite{
				Id: int32(groupID),
			},
			Layout: &models.CreateInstanceBodyInstanceLayout{
				Id: utils.JSONNumber(d.Get("layout_id")),
			},
		},
		Volumes:           volumes,
		NetworkInterfaces: networks,
		Config:            config,
		Tags:              tags,
	}

	resp, err := i.iClient.CreateAnInstance(ctx, i.serviceInstance, req)
	if err != nil {
		return err
	}
	d.SetId(strconv.Itoa(int(resp.Instance.Id)))

	return nil
}

// Update instance including poweroff, powerOn, restart, suspend
// changing network, volumes and instance properties such as labels
// groups and tags
func (i *instance) Update(ctx context.Context, d *schema.ResourceData) error {
	return nil
}

// Delete instance and set ID as ""
func (i *instance) Delete(ctx context.Context, d *schema.ResourceData) error {
	id, err := utils.ParseInt(d.Id())
	if err != nil {
		return err
	}
	res, err := i.iClient.DeleteAnInstance(ctx, i.serviceInstance, int32(id))
	if err != nil {
		return err
	}
	if !res.Success {
		return fmt.Errorf("%s", res.Message)
	}
	d.SetId("")

	return nil
}

// Read instance and set state values accordingly
func (i *instance) Read(ctx context.Context, d *schema.ResourceData) error {
	id, err := utils.ParseInt(d.Id())
	if err != nil {
		return err
	}
	resp, err := i.iClient.GetASpecificInstance(ctx, i.serviceInstance, int32(id))
	if err != nil {
		return err
	}
	d.SetId(strconv.Itoa(int(resp.Instance.Id)))
	if err := d.Set("status", resp.Instance.Status); err != nil {
		return err
	}

	return nil
}

func getVolume(v interface{}) ([]models.CreateInstanceBodyVolumes, error) {
	volumes, err := utils.ListToMap(v)
	if err != nil {
		return nil, err
	}
	volumesModel := make([]models.CreateInstanceBodyVolumes, 0, len(volumes))
	for i := range volumes {
		vID, err := utils.ParseInt(volumes[i]["size"].(string))
		if err != nil {
			return nil, err
		}
		volumesModel = append(volumesModel, models.CreateInstanceBodyVolumes{
			Id:          -1,
			Name:        volumes[i]["name"].(string),
			Size:        int32(vID),
			DatastoreId: volumes[i]["datastore_id"],
			RootVolume:  true,
		})
	}

	return volumesModel, nil
}

func getNetwork(v interface{}) ([]models.CreateInstanceBodyNetworkInterfaces, error) {
	networksMap, err := utils.ListToMap(v)
	if err != nil {
		return nil, err
	}
	networks := make([]models.CreateInstanceBodyNetworkInterfaces, 0, len(networksMap))
	for _, n := range networksMap {
		networks = append(networks, models.CreateInstanceBodyNetworkInterfaces{
			Network: &models.CreateInstanceBodyNetwork{
				Id: int32(n["id"].(int)),
			},
			NetworkInterfaceTypeId: utils.JSONNumber(n["interface_type_id"]),
		})
	}

	return networks, nil
}

func getConfig(v interface{}) (*models.CreateInstanceBodyConfig, error) {
	c, err := utils.SetToMap(v)
	if err != nil {
		return nil, err
	}
	config := &models.CreateInstanceBodyConfig{
		ResourcePoolId: utils.JSONNumber(c["resource_pool_id"]),
		Template:       int32(c["template_id"].(int)),
	}

	return config, nil
}

func getTags(v interface{}) ([]models.CreateInstanceBodyTag, error) {
	t, err := utils.MapTopMap(v)
	if err != nil {
		return nil, err
	}
	tags := make([]models.CreateInstanceBodyTag, 0, len(t))
	for k, v := range t {
		tags = append(tags, models.CreateInstanceBodyTag{
			Name:  k,
			Value: v.(string),
		})
	}

	return tags, nil
}
