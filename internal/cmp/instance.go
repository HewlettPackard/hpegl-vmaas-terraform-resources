// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/hpe-hcss/vmaas-cmp-go-sdk/pkg/client"
	"github.com/hpe-hcss/vmaas-cmp-go-sdk/pkg/models"
	"github.com/hpe-hcss/vmaas-terraform-resources/internal/logger"
	"github.com/hpe-hcss/vmaas-terraform-resources/internal/utils"
)

const (
	instanceCloneRetryDelay   = time.Second * 60
	instanceCloneRetryTimeout = time.Second * 10
	instanceCloneRetryCount   = 10
)

// instance implements functions related to cmp instances
type instance struct {
	// expose Instance API service to instances related operations
	iClient *client.InstancesApiService
}

func newInstance(iClient *client.InstancesApiService) *instance {
	return &instance{
		iClient: iClient,
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
				Code: d.GetString("instance_type_code"),
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
			HostName: d.GetString("hostname"),
		},
		Evars:             getEvars(d.GetMap("evars")),
		Labels:            d.GetStringList("labels"),
		Volumes:           getVolume(d.GetListMap("volume")),
		NetworkInterfaces: getNetwork(d.GetListMap("network")),
		Config:            getConfig(c),
		Tags:              getTags(d.GetMap("tags")),
		LayoutSize:        d.GetInt("scale"),
		PowerScheduleType: utils.JSONNumber(d.GetInt("power_schedule_id")),
	}
	if req.Instance.InstanceType.Code == vmware {
		templateID := c["template_id"]
		if templateID == nil {
			return errors.New("error, template id is required for vmware instance type")
		}
		req.Config.Template = templateID.(int)
	}
	cloneData := d.GetSMap("clone", true)
	// Pre check
	if err := d.Error(); err != nil {
		return err
	}
	// Get template id

	var getInstanceBody models.GetInstanceResponseInstance
	// check whether vm to be cloned?
	if cloneData != nil {
		req.CloneName = req.Instance.Name
		req.Instance.Name = ""
		sourceID, _ := strconv.Atoi(cloneData["source_instance_id"].(string))

		// clone the instance
		logger.Info("Cloning the instance with ", sourceID)
		respClone, err := utils.Retry(func() (interface{}, error) {
			return i.iClient.CloneAnInstance(ctx, sourceID, req)
		})
		if err != nil {
			return err
		}

		logger.Info("Check clone success")
		isCloneSuccess := respClone.(models.SuccessOrErrorMessage)
		if !isCloneSuccess.Success {
			return fmt.Errorf("failed to clone, error: %s", isCloneSuccess.Message)
		}

		logger.Info("Get all instances")
		customRetry := &utils.CustomRetry{
			Delay:        instanceCloneRetryDelay,
			RetryTimeout: instanceCloneRetryTimeout,
			RetryCount:   instanceCloneRetryCount,
		}
		// get cloned instance ID
		instancesResp, err := customRetry.Retry(func() (interface{}, error) {
			return i.iClient.GetAllInstances(ctx, map[string]string{
				nameKey: req.CloneName,
			})
		})
		if err != nil {
			return err
		}

		instancesList := instancesResp.(models.Instances)
		if len(instancesList.Instances) != 1 {
			return errors.New("get cloned instance is failed")
		}
		logger.Info("Instance id = ", instancesList.Instances[0].Id)
		getInstanceBody = instancesList.Instances[0]
	} else {
		// create instance
		respVM, err := utils.Retry(func() (interface{}, error) {
			return i.iClient.CreateAnInstance(ctx, req)
		})
		if err != nil {
			return err
		}
		getInstanceBody = *respVM.(models.GetInstanceResponse).Instance
	}
	d.SetID(getInstanceBody.Id)

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

	volumes := d.GetListMap("volume")
	for i := range volumes {
		volumes[i]["id"] = instance.Instance.Volumes[i].Id
	}
	d.Set("volume", volumes)
	d.SetID(instance.Instance.Id)
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
		VMwareFolderId: c["vm_folder"].(string),
		CreateUser:     c["create_user"].(bool),
		SmbiosAssetTag: c["asset_tag"].(string),
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

func getEvars(evars map[string]interface{}) []models.GetInstanceResponseInstanceEvars {
	evarModel := make([]models.GetInstanceResponseInstanceEvars, 0, len(evars))
	for k, v := range evars {
		evarModel = append(evarModel, models.GetInstanceResponseInstanceEvars{
			Name:   k,
			Value:  v.(string),
			Export: true,
			Masked: false,
		})
	}

	return evarModel
}
