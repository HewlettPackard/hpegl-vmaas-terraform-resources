// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hpe-hcss/vmaas-cmp-go-sdk/pkg/client"
	"github.com/hpe-hcss/vmaas-cmp-go-sdk/pkg/models"
	"github.com/hpe-hcss/vmaas-terraform-resources/internal/logger"
	"github.com/hpe-hcss/vmaas-terraform-resources/internal/utils"
	"github.com/hpe-hcss/vmaas-terraform-resources/pkg/auth"
)

const (
	instanceCloneRetryDelay   = time.Second * 60
	instanceCloneRetryTimeout = time.Second * 10
	instanceCloneRetryCount   = 30
)

// instance implements functions related to cmp instances
type instance struct {
	// expose Instance API service to instances related operations
	iClient *client.InstancesAPIService
}

func newInstance(iClient *client.InstancesAPIService) *instance {
	return &instance{
		iClient: iClient,
	}
}

// Create instance
func (i *instance) Create(ctx context.Context, d *utils.Data, meta interface{}) error {
	logger.Debug("Creating new instance")

	err := validateVolumeNameIsUnique(d.GetListMap("volume"))
	if err != nil{
		return err
	}

	c := d.GetListMap("config")[0]
	req := &models.CreateInstanceBody{
		ZoneID: d.GetJSONNumber("cloud_id"),
		Instance: &models.CreateInstanceBodyInstance{
			Name: d.GetString("name"),
			InstanceType: &models.CreateInstanceBodyInstanceInstanceType{
				Code: d.GetString("instance_type_code"),
			},
			Plan: &models.CreateInstanceBodyInstancePlan{
				ID: d.GetJSONNumber("plan_id"),
			},
			Site: &models.CreateInstanceBodyInstanceSite{
				ID: d.GetInt("group_id"),
			},
			Layout: &models.CreateInstanceBodyInstanceLayout{
				ID: d.GetJSONNumber("layout_id"),
			},
			HostName:          d.GetString("hostname"),
			EnvironmentPrefix: d.GetString("env_prefix"),
		},
		Environment:       d.GetString("environment_code"),
		Ports:             getPorts(d.GetListMap("port")),
		Evars:             getEvars(d.GetMap("evars")),
		Labels:            d.GetStringList("labels"),
		Volumes:           getVolume(d.GetListMap("volume")),
		NetworkInterfaces: getNetwork(d.GetListMap("network")),
		Config:            getConfig(c),
		Tags:              getTags(d.GetMap("tags")),
		LayoutSize:        d.GetInt("scale"),
		PowerScheduleType: utils.JSONNumber(d.GetInt("power_schedule_id")),
	}

	// Get template id instance type is vmware
	if strings.ToLower(req.Instance.InstanceType.Code) == vmware {
		templateID := c["template_id"]
		if templateID == nil {
			return errors.New("error, template id is required for vmware instance type")
		}
		req.Config.Template = templateID.(int)
	}
	cData := d.GetListMap("clone")
	// Pre check
	if err := d.Error(); err != nil {
		return err
	}

	var getInstanceBody models.GetInstanceResponseInstance
	// check whether vm to be cloned?
	if len(cData) > 0 {
		cloneData := cData[0]
		req.CloneName = req.Instance.Name
		req.Instance.Name = ""
		sourceID, _ := strconv.Atoi(cloneData["source_instance_id"].(string))

		// clone the instance
		logger.Info("Cloning the instance with ", sourceID)
		respClone, err := utils.Retry(func() (interface{}, error) {
			auth.SetScmClientToken(&ctx, meta)
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
			Cond: func(resp interface{}, err error) bool {
				if err != nil {
					return false
				}

				instancesList := resp.(models.Instances)

				return len(instancesList.Instances) == 1
			},
		}
		// get cloned instance ID
		instancesResp, err := customRetry.Retry(func() (interface{}, error) {
			auth.SetScmClientToken(&ctx, meta)
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

		getInstanceBody = instancesList.Instances[0]
	} else {
		// create instance
		respVM, err := utils.Retry(func() (interface{}, error) {
			auth.SetScmClientToken(&ctx, meta)
			return i.iClient.CreateAnInstance(ctx, req)
		})
		if err != nil {
			return err
		}
		getInstanceBody = *respVM.(models.GetInstanceResponse).Instance
	}
	// Upon creation instance will be in poweron state. Check any other
	// power state provided and do accordingly
	powerOp := d.GetString("power")
	if powerOp != "" && powerOp != utils.PowerOn {
		return fmt.Errorf("power operation %s is not permitted while creating an instance", powerOp)
	}
	d.SetID(getInstanceBody.ID)

	// post check
	return d.Error()
}

// Update instance including poweroff, powerOn, restart, suspend
// changing volumes and instance properties such as labels
// groups and tags
func (i *instance) Update(ctx context.Context, d *utils.Data, meta interface{}) error {
	logger.Debug("Updating the instance")

	err := validateVolumeNameIsUnique(d.GetListMap("volume"))
	if err != nil{
		return err
	}

	id := d.GetID()
	if d.HasChanged("name") || d.HasChanged("group_id") || d.HasChanged(
		"tags") || d.HasChanged("labels") || d.HasChanged("environment_code") {
		addTags, removeTags := compareTags(d.GetChangedMap("tags"))
		updateReq := &models.UpdateInstanceBody{
			Instance: &models.UpdateInstanceBodyInstance{
				Name: d.GetString("name"),
				Site: &models.CreateInstanceBodyInstanceSite{
					ID: d.GetInt("group_id"),
				},
				AddTags:           addTags,
				RemoveTags:        removeTags,
				Labels:            d.GetStringList("labels"),
				PowerScheduleType: utils.JSONNumber(d.GetInt("power_schedule_id")),
				InstanceContext:   d.GetString("environment_code"),
			},
		}

		if err := d.Error(); err != nil {
			return err
		}
		// update instance
		_, err := utils.Retry(func() (interface{}, error) {
			auth.SetScmClientToken(&ctx, meta)
			return i.iClient.UpdatingAnInstance(ctx, id, updateReq)
		})
		if err != nil {
			return err
		}
	}

	if d.HasChanged("volume") || d.HasChanged("plan_id") {
		volumes := compareVolumes(d.GetChangedListMap("volume"))
		resizeReq := &models.ResizeInstanceBody{
			Instance: &models.ResizeInstanceBodyInstance{
				Plan: &models.ResizeInstanceBodyInstancePlan{
					ID: d.GetInt("plan_id"),
				},
			},
			Volumes: resizeVolume(volumes),
		}
		if err := d.Error(); err != nil {
			return err
		}
		_, err := utils.Retry(func() (interface{}, error) {
			auth.SetScmClientToken(&ctx, meta)
			return i.iClient.ResizeAnInstance(ctx, id, resizeReq)
		})
		if err != nil {
			return err
		}
	}

	if d.HasChanged("power") {
		// Do power operation only if backend is in different state
		resp, err := utils.Retry(func() (interface{}, error) {
			auth.SetScmClientToken(&ctx, meta)
			return i.iClient.GetASpecificInstance(ctx, id)
		})
		if err != nil {
			return err
		}
		getInstance := resp.(models.GetInstanceResponse)
		status := utils.ParsePowerState(getInstance.Instance.Status)
		powerOp := d.GetString("power")
		if powerOp != status {
			if err := i.powerOperation(ctx, id, meta, status, d.GetString("power")); err != nil {
				return err
			}
		}
	}

	return d.Error()
}

// Delete instance and set ID as ""
func (i *instance) Delete(ctx context.Context, d *utils.Data, meta interface{}) error {
	id := d.GetID()
	logger.Debugf("Deleting instance with ID : %d", id)

	// Precheck
	if err := d.Error(); err != nil {
		return err
	}

	resp, err := utils.Retry(func() (interface{}, error) {
		auth.SetScmClientToken(&ctx, meta)
		return i.iClient.DeleteAnInstance(ctx, id)
	})
	deleResp := resp.(models.SuccessOrErrorMessage)
	if err != nil {
		return err
	}
	if !deleResp.Success {
		return fmt.Errorf("%s", deleResp.Message)
	}

	// post check
	return d.Error()
}

// Read instance and set state values accordingly
func (i *instance) Read(ctx context.Context, d *utils.Data, meta interface{}) error {
	id := d.GetID()

	logger.Debug("Get instance with ID %d", id)

	// Precheck
	if err := d.Error(); err != nil {
		return err
	}

	resp, err := utils.Retry(func() (interface{}, error) {
		auth.SetScmClientToken(&ctx, meta)
		return i.iClient.GetASpecificInstance(ctx, id)
	})
	if err != nil {
		return err
	}
	instance := resp.(models.GetInstanceResponse)

	volumes := d.GetListMap("volume")
	if len(volumes) != len(instance.Instance.Volumes){
		return fmt.Errorf("volume name should be unique")
	}
	for i := range volumes {
		volumes[i]["id"] = instance.Instance.Volumes[i].ID
	}
	d.Set("volume", volumes)
	d.SetID(instance.Instance.ID)
	d.SetString("status", instance.Instance.Status)

	// post check
	return d.Error()
}

func getVolume(volumes []map[string]interface{}) []models.CreateInstanceBodyVolumes {
	volumesModel := make([]models.CreateInstanceBodyVolumes, 0, len(volumes))
	logger.Debug(volumes)
	for i := range volumes {
		volumesModel = append(volumesModel, models.CreateInstanceBodyVolumes{
			ID:          -1,
			Name:        volumes[i]["name"].(string),
			Size:        volumes[i]["size"].(int),
			DatastoreID: volumes[i]["datastore_id"],
			RootVolume:  volumes[i]["root"].(bool),
		})
	}

	return volumesModel
}

// Mapping volume data to model
func resizeVolume(volumes []map[string]interface{}) []models.ResizeInstanceBodyInstanceVolumes {
	volumesModel := make([]models.ResizeInstanceBodyInstanceVolumes, 0, len(volumes))
	logger.Debug(volumes)
	for i := range volumes {
		volumesModel = append(volumesModel, models.ResizeInstanceBodyInstanceVolumes{
			ID:          utils.JSONNumber(volumes[i]["id"]),
			Name:        volumes[i]["name"].(string),
			Size:        volumes[i]["size"].(int),
			DatastoreID: volumes[i]["datastore_id"],
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
				ID: n["id"].(int),
			},
			NetworkInterfaceTypeID: utils.JSONNumber(n["interface_id"]),
		})
	}

	return networks
}

func getConfig(c map[string]interface{}) *models.CreateInstanceBodyConfig {
	config := &models.CreateInstanceBodyConfig{
		ResourcePoolID: utils.JSONNumber(c["resource_pool_id"]),
		NoAgent:        strconv.FormatBool(c["no_agent"].(bool)),
		VMwareFolderID: c["vm_folder"].(string),
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

func getPorts(ports []map[string]interface{}) []models.CreateInstancePorts {
	pModels := make([]models.CreateInstancePorts, 0, len(ports))
	for _, p := range ports {
		pModels = append(pModels, models.CreateInstancePorts{
			Name: p["name"].(string),
			Port: p["port"].(string),
			Lb:   p["lb"].(string),
		})
	}

	return pModels
}

// Function to compare tags and based on new and old data assign to AddTags or Removetags
func compareTags(org, new map[string]interface{}) ([]models.CreateInstanceBodyTag, []models.CreateInstanceBodyTag) {
	addTags := make([]models.CreateInstanceBodyTag, 0, len(new))
	removeTags := make([]models.CreateInstanceBodyTag, 0, len(new))
	for k, v := range new {
		addTags = append(addTags, models.CreateInstanceBodyTag{
			Name:  k,
			Value: v.(string),
		})
	}

	for k, v := range org {
		if _, ok := new[k]; !ok {
			removeTags = append(removeTags, models.CreateInstanceBodyTag{
				Name:  k,
				Value: v.(string),
			})
		}
	}

	return addTags, removeTags
}

// Function to compare previous and new(from terraform) volume data and assign proper ids based on name.
// Volume name should be unique
func compareVolumes(org, new []map[string]interface{}) []map[string]interface{} {
	for i := range new {
		new[i]["id"] = -1
		for j := range org {
			if new[i]["name"] == org[j]["name"] {
				new[i]["id"] = org[j]["id"]
				new[i]["size"] = org[j]["size"]

				break
			}
		}
	}

	return new
}

func (i *instance) powerOperation(ctx context.Context, instanceID int, meta interface{}, oldOp, operation string) error {
	var err error
	err = validatePowerTransition(oldOp, operation)
	if err != nil {
		return err
	}
	switch operation {
	case utils.PowerOn:
		_, err = utils.Retry(func() (interface{}, error) {
			auth.SetScmClientToken(&ctx, meta)
			return i.iClient.StartAnInstance(ctx, instanceID)
		})
	case utils.PowerOff:
		_, err = utils.Retry(func() (interface{}, error) {
			auth.SetScmClientToken(&ctx, meta)
			return i.iClient.StopAnInstance(ctx, instanceID)
		})
	case utils.Suspend:
		_, err = utils.Retry(func() (interface{}, error) {
			auth.SetScmClientToken(&ctx, meta)
			return i.iClient.SuspendAnInstance(ctx, instanceID)
		})
	case utils.Restart:
		_, err = utils.Retry(func() (interface{}, error) {
			auth.SetScmClientToken(&ctx, meta)
			return i.iClient.RestartAnInstance(ctx, instanceID)
		})
	default:
		return fmt.Errorf("power operation not allowed from %s state", operation)
	}

	return err
}

func validatePowerTransition(oldPower, newPower string) error {
	if oldPower == utils.PowerOn || oldPower == utils.Restart {
		if newPower == utils.PowerOff || newPower == utils.Suspend || newPower == utils.Restart {
			return nil
		}
	} else {
		if newPower == utils.PowerOn {
			return nil
		}
	}

	return fmt.Errorf("power operation not allowed from %s state to %s state", oldPower, newPower)
}

func validateVolumeNameIsUnique(vol []map[string]interface{}) error{
	volumes  := make(map[string]bool)
	for _, v := range vol{

		if _, ok := volumes[v["name"].(string)]; !ok{
			volumes[v["name"].(string)] = true
			continue
		}

		return fmt.Errorf("volume names should be unique")
	}
	return nil
}