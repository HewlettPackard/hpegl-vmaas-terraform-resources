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
	iClient *client.InstancesApiService
	tClient *client.VirtualImagesApiService
}

func newInstance(iClient *client.InstancesApiService, tClient *client.VirtualImagesApiService) *instance {
	return &instance{
		iClient: iClient,
		tClient: tClient,
	}
}

// Create instance
func (i *instance) Create(ctx context.Context, d *utils.Data) error {
	logger.Debug("Creating new instance")

	c := d.GetSMap("config")
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
				Id: d.GetInt("group_id"),
			},
			Layout: &models.CreateInstanceBodyInstanceLayout{
				Id: d.GetJSONNumber("layout_id"),
			},
			Type: d.GetString("instance_code"),
		},
		Volumes:           getVolume(d.GetListMap("volumes")),
		NetworkInterfaces: getNetwork(d.GetListMap("networks")),
		Config:            getConfig(c),
		Tags:              getTags(d.GetMap("tags")),
	}

	// Pre check
	if err := d.Error(); err != nil {
		return err
	}
	// Get template id
	vResp, err := i.tClient.GetAllVirtualImages(ctx, map[string]string{
		nameKey: c["template"].(string),
	})
	if err != nil {
		return err
	}
	if len(vResp.VirtualImages) != 1 {
		return fmt.Errorf(errExactMatch, "templates")
	}
	req.Config.Template = vResp.VirtualImages[0].ID

	// create instance
	resp, err := utils.Retry(func() (interface{}, error) {
		return i.iClient.CreateAnInstance(ctx, req)
	})
	if err != nil {
		return err
	}
	instance := resp.(models.GetInstanceResponse)
	d.SetID(strconv.Itoa(instance.Instance.Id))

	// post check
	return d.Error()
}

// Update instance including poweroff, powerOn, restart, suspend
// changing network, volumes and instance properties such as labels
// groups and tags
func (i *instance) Update(ctx context.Context, d *utils.Data) error {
	logger.Debug("Updating the instance")

	return nil
}

// Delete instance and set ID as ""
func (i *instance) Delete(ctx context.Context, d *utils.Data) error {
	id := d.GetID()
	logger.Debugf("Deleting instance with ID : %d", id)

	// Precheck
	if err := d.Error(); err != nil {
		return err
	}

	resp, err := utils.Retry(func() (interface{}, error) {
		return i.iClient.DeleteAnInstance(ctx, id)
	})
	deleResp := resp.(models.SuccessOrErrorMessage)
	if err != nil {
		return err
	}
	if !deleResp.Success {
		return fmt.Errorf("%s", deleResp.Message)
	}
	d.SetID("")

	// post check
	return d.Error()
}

// Read instance and set state values accordingly
func (i *instance) Read(ctx context.Context, d *utils.Data) error {
	id := d.GetID()

	logger.Debug("Get instance with ID %d", id)

	// Precheck
	if err := d.Error(); err != nil {
		return err
	}

	resp, err := utils.Retry(func() (interface{}, error) {
		return i.iClient.GetASpecificInstance(ctx, id)
	})
	if err != nil {
		return err
	}
	instance := resp.(models.GetInstanceResponse)
	d.SetID(strconv.Itoa(instance.Instance.Id))
	d.SetString("status", instance.Instance.Status)

	// post check
	return d.Error()
}

func getVolume(volumes []map[string]interface{}) []models.CreateInstanceBodyVolumes {
	volumesModel := make([]models.CreateInstanceBodyVolumes, 0, len(volumes))
	logger.Debug(volumes)
	for i := range volumes {
		volumesModel = append(volumesModel, models.CreateInstanceBodyVolumes{
			Id:          -1,
			Name:        volumes[i]["name"].(string),
			Size:        volumes[i]["size"].(int),
			DatastoreId: volumes[i]["datastore_id"],
			RootVolume:  volumes[i]["root"].(bool),
		})
	}

	return volumesModel
}

func getNetwork(networksMap []map[string]interface{}) []models.CreateInstanceBodyNetworkInterfaces {
	networks := make([]models.CreateInstanceBodyNetworkInterfaces, 0, len(networksMap))
	for _, n := range networksMap {
		networks = append(networks, models.CreateInstanceBodyNetworkInterfaces{
			Network: &models.CreateInstanceBodyNetwork{
				Id: n["id"].(int),
			},
		})
	}

	return networks
}

func getConfig(c map[string]interface{}) *models.CreateInstanceBodyConfig {
	config := &models.CreateInstanceBodyConfig{
		ResourcePoolId: utils.JSONNumber(c["resource_pool_id"]),
		NoAgent:        strconv.FormatBool(c["no_agent"].(bool)),
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
