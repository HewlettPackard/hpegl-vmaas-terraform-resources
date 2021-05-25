// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hpe-hcss/vmaas-cmp-go-sdk/pkg/client"
	"github.com/hpe-hcss/vmaas-cmp-go-sdk/pkg/models"
	"github.com/hpe-hcss/vmaas-terraform-resources/internal/logger"
	"github.com/hpe-hcss/vmaas-terraform-resources/internal/utils"
)

// instance implements functions related to cmp instances
type instance struct {
	// expose Instance API service to instances related operations
	iClient           *client.InstancesApiService
	serviceInstanceId string
	log               logger.Logger
}

func newInstance(iClient *client.InstancesApiService, serviceInstanceId string) *instance {
	return &instance{
		iClient:           iClient,
		serviceInstanceId: serviceInstanceId,
	}
}

// Create instance
func (i *instance) Create(ctx context.Context, d *utils.Data) error {
	i.log.Debug("Creating new instance")

	req := &models.CreateInstanceBody{
		ZoneId: d.GetJSONNumber("cloud_id"),
		Instance: &models.CreateInstanceBodyInstance{
			Name: d.GetString("name"),
			InstanceType: &models.CreateInstanceBodyInstanceInstanceType{
				Code: d.GetString("instance_code"),
			},
			Plan: &models.CreateInstanceBodyInstancePlan{
				Id: d.GetJSONNumber("plan_id"),
			},
			Site: &models.CreateInstanceBodyInstanceSite{
				Id: int32(d.GetInt("group_id")),
			},
			Layout: &models.CreateInstanceBodyInstanceLayout{
				Id: d.GetJSONNumber("layout_id"),
			},
		},
		Volumes:           getVolume(d.GetListMap("volumes")),
		NetworkInterfaces: getNetwork(d.GetListMap("networks")),
		Config:            getConfig(d.GetSMap("config")),
		Tags:              getTags(d.GetMap("tags")),
	}

	resp, err := i.iClient.CreateAnInstance(ctx, i.serviceInstanceId, req)
	if err != nil {
		return err
	}
	d.SetID(strconv.Itoa(int(resp.Instance.Id)))

	if d.HaveError() {
		return fmt.Errorf("%s", "error in d ")
	}
	return nil
}

// Update instance including poweroff, powerOn, restart, suspend
// changing network, volumes and instance properties such as labels
// groups and tags
func (i *instance) Update(ctx context.Context, d *utils.Data) error {
	i.log.Debug("Updating the instance")

	return nil
}

// Delete instance and set ID as ""
func (i *instance) Delete(ctx context.Context, d *utils.Data) error {
	id := d.GetID()
	i.log.Debugf("Deleting instance with ID : %d", id)

	res, err := i.iClient.DeleteAnInstance(ctx, i.serviceInstanceId, int32(id))
	if err != nil {
		return err
	}
	if !res.Success {
		return fmt.Errorf("%s", res.Message)
	}
	d.SetID("")

	if d.HaveError() {
		return fmt.Errorf("%s", "error in d ")
	}

	return nil
}

// Read instance and set state values accordingly
func (i *instance) Read(ctx context.Context, d *utils.Data) error {
	id := d.GetID()

	i.log.Debug("Get instance with ID %d", id)
	resp, err := i.iClient.GetASpecificInstance(ctx, i.serviceInstanceId, int32(id))
	if err != nil {
		return err
	}
	d.SetID(strconv.Itoa(int(resp.Instance.Id)))
	d.SetString("status", resp.Instance.Status)

	if d.HaveError() {
		return fmt.Errorf("%s", "error in d ")
	}

	return nil
}

func getVolume(volumes []map[string]interface{}) []models.CreateInstanceBodyVolumes {
	volumesModel := make([]models.CreateInstanceBodyVolumes, 0, len(volumes))
	for i := range volumes {
		vID, _ := utils.ParseInt(volumes[i]["size"].(string))
		volumesModel = append(volumesModel, models.CreateInstanceBodyVolumes{
			Id:          -1,
			Name:        volumes[i]["name"].(string),
			Size:        int32(vID),
			DatastoreId: volumes[i]["datastore_id"],
			RootVolume:  true,
		})
	}

	return volumesModel
}

func getNetwork(networksMap []map[string]interface{}) []models.CreateInstanceBodyNetworkInterfaces {
	networks := make([]models.CreateInstanceBodyNetworkInterfaces, 0, len(networksMap))
	for _, n := range networksMap {
		networks = append(networks, models.CreateInstanceBodyNetworkInterfaces{
			Network: &models.CreateInstanceBodyNetwork{
				Id: int32(n["id"].(int)),
			},
		})
	}

	return networks
}

func getConfig(c map[string]interface{}) *models.CreateInstanceBodyConfig {
	config := &models.CreateInstanceBodyConfig{
		ResourcePoolId: utils.JSONNumber(c["resource_pool_id"]),
		Template:       int32(c["template_id"].(int)),
	}

	return config
}

func getTags(t map[string]interface{}) []models.CreateInstanceBodyTag {
	tags := make([]models.CreateInstanceBodyTag, 0, len(t))
	for k, v := range t {
		tags = append(tags, models.CreateInstanceBodyTag{
			Name:  k,
			Value: v.(string),
		})
	}

	return tags
}
