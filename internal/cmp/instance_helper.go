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

type iClient interface {
	getIClient() *client.InstancesAPIService
}

// Update instance including poweroff, powerOn, restart, suspend
// changing volumes and instance properties such as labels
// groups and tags
func updateInstance(ctx context.Context, iclient iClient, d *utils.Data, meta interface{}) error {
	logger.Debug("Updating the instance")
	err := instanceValidateVolumeNameIsUnique(d.GetListMap("volume"))
	if err != nil {
		return err
	}

	id := d.GetID()
	if d.HasChanged("name") || d.HasChanged("group_id") || d.HasChanged("tags") ||
		d.HasChanged("labels") || d.HasChanged("environment_code") ||
		d.HasChanged("power_schedule_id") {
		addTags, removeTags := instanceCompareTags(d.GetChangedMap("tags"))
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
		_, err := utils.Retry(ctx, meta, func(ctx context.Context) (interface{}, error) {
			return iclient.getIClient().UpdatingAnInstance(ctx, id, updateReq)
		})
		if err != nil {
			return err
		}
	}

	if d.HasChanged("volume") {
		volumes := instanceCompareVolumes(d.GetChangedListMap("volume"))
		resizeReq := &models.ResizeInstanceBody{
			Instance: &models.ResizeInstanceBodyInstance{
				Plan: &models.ResizeInstanceBodyInstancePlan{
					ID: d.GetInt("plan_id"),
				},
			},
			Volumes: instanceResizeVolume(volumes),
		}
		if err := d.Error(); err != nil {
			return err
		}
		updateResp, err := utils.Retry(ctx, meta, func(ctx context.Context) (interface{}, error) {
			return iclient.getIClient().ResizeAnInstance(ctx, id, resizeReq)
		})
		if err != nil {
			return err
		}
		if !updateResp.(models.ResizeInstanceResponse).Success {
			return fmt.Errorf("%s", "failed to resize")
		}
	}

	if d.HasChanged("power") || d.HasChanged("restart_instance") {
		// Do power operation only if backend is in different state
		// restart only if instance in actual is in power-on state
		resp, err := utils.Retry(ctx, meta, func(ctx context.Context) (interface{}, error) {
			return iclient.getIClient().GetASpecificInstance(ctx, id)
		})
		if err != nil {
			return err
		}
		getInstance := resp.(models.GetInstanceResponse)
		status := utils.ParsePowerState(getInstance.Instance.Status)
		powerOp := d.GetString("power")
		if powerOp != status {
			if err := instanceDoPowerTask(ctx, iclient, id, meta, status, d.GetString("power")); err != nil {
				return err
			}
		} else if d.HasChanged("restart_instance") {
			if err := instanceDoPowerTask(ctx, iclient, id, meta, status, utils.Restart); err != nil {
				return err
			}
		}
	}

	return d.Error()
}

// Delete instance and set ID as ""
func deleteInstance(ctx context.Context, iclient iClient, d *utils.Data, meta interface{}) error {
	id := d.GetID()
	logger.Debugf("Deleting instance with ID : %d", id)

	// Precheck
	if err := d.Error(); err != nil {
		return err
	}

	resp, err := utils.Retry(ctx, meta, func(ctx context.Context) (interface{}, error) {
		return iclient.getIClient().DeleteAnInstance(ctx, id)
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

func instanceGetVolume(volumes []map[string]interface{}) []models.CreateInstanceBodyVolumes {
	volumesModel := make([]models.CreateInstanceBodyVolumes, 0, len(volumes))
	logger.Debug(volumes)
	for i := range volumes {
		volumesModel = append(volumesModel, models.CreateInstanceBodyVolumes{
			ID:          -1,
			Name:        volumes[i]["name"].(string),
			Size:        volumes[i]["size"].(int),
			DatastoreID: volumes[i]["datastore_id"],
		})
	}
	volumesModel[0].RootVolume = true

	return volumesModel
}

// Mapping volume data to model
func instanceResizeVolume(volumes []map[string]interface{}) []models.ResizeInstanceBodyInstanceVolumes {
	volumesModel := make([]models.ResizeInstanceBodyInstanceVolumes, 0, len(volumes))
	for i := range volumes {
		volumesModel = append(volumesModel, models.ResizeInstanceBodyInstanceVolumes{
			ID:          utils.JSONNumber(volumes[i]["id"]),
			Name:        volumes[i]["name"].(string),
			Size:        volumes[i]["size"].(int),
			DatastoreID: volumes[i]["datastore_id"],
		})
	}

	return volumesModel
}

func instanceGetNetwork(networksMap []map[string]interface{}) []models.CreateInstanceBodyNetworkInterfaces {
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

func instanceGetConfig(c map[string]interface{}) *models.CreateInstanceBodyConfig {
	config := &models.CreateInstanceBodyConfig{
		ResourcePoolID: utils.JSONNumber(c["resource_pool_id"]),
		NoAgent:        strconv.FormatBool(c["no_agent"].(bool)),
		VMwareFolderID: c["vm_folder"].(string),
		CreateUser:     c["create_user"].(bool),
		SmbiosAssetTag: c["asset_tag"].(string),
	}

	return config
}

func instanceGetTags(t map[string]interface{}) []models.CreateInstanceBodyTag {
	tags := make([]models.CreateInstanceBodyTag, 0, len(t))
	for k, v := range t {
		tags = append(tags, models.CreateInstanceBodyTag{
			Name:  k,
			Value: v.(string),
		})
	}

	return tags
}

func instanceGetEvars(evars map[string]interface{}) []models.GetInstanceResponseInstanceEvars {
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

func instanceGetPorts(ports []map[string]interface{}) []models.CreateInstancePorts {
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
func instanceCompareTags(org, new map[string]interface{}) ([]models.CreateInstanceBodyTag, []models.CreateInstanceBodyTag) {
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
func instanceCompareVolumes(org, new []map[string]interface{}) []map[string]interface{} {
	for i := range new {
		new[i]["id"] = -1
		for j := range org {
			if new[i]["name"] == org[j]["name"] {
				new[i]["id"] = org[j]["id"]

				break
			}
		}
	}

	return new
}

func instanceDoPowerTask(
	ctx context.Context,
	iclient iClient,
	instanceID int,
	meta interface{},
	currState,
	newOp string) error {
	var err error
	err = instanceValidatePowerTransition(currState, newOp)
	if err != nil {
		return err
	}
	switch newOp {
	case utils.PowerOn:
		_, err = utils.Retry(ctx, meta, func(ctx context.Context) (interface{}, error) {
			_, err := iclient.getIClient().StartAnInstance(ctx, instanceID)

			return nil, err
		})
	case utils.PowerOff:
		_, err = utils.Retry(ctx, meta, func(ctx context.Context) (interface{}, error) {
			_, err := iclient.getIClient().StopAnInstance(ctx, instanceID)

			return nil, err
		})
	case utils.Suspend:
		_, err = utils.Retry(ctx, meta, func(ctx context.Context) (interface{}, error) {
			_, err := iclient.getIClient().SuspendAnInstance(ctx, instanceID)

			return nil, err
		})
	case utils.Restart:
		_, err = utils.Retry(ctx, meta, func(ctx context.Context) (interface{}, error) {
			_, err := iclient.getIClient().RestartAnInstance(ctx, instanceID)

			return nil, err
		})
	}

	return err
}

func instanceValidatePowerTransition(oldPower, newPower string) error {
	if oldPower == utils.PowerOn {
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

func instanceValidateVolumeNameIsUnique(vol []map[string]interface{}) error {
	volumes := make(map[string]bool)
	for _, v := range vol {
		if _, ok := volumes[v["name"].(string)]; !ok {
			volumes[v["name"].(string)] = true

			continue
		}

		return fmt.Errorf("volume names should be unique")
	}

	return nil
}

func instanceValidatePower(powerOp string) error {
	if powerOp != "" && powerOp != utils.PowerOn {
		return fmt.Errorf("power operation %s is not permitted while creating an instance", powerOp)
	}

	return nil
}

func instanceCloneCompareVolume(
	vSchemas []map[string]interface{},
	vModels []models.GetInstanceResponseInstanceVolumes,
) []models.CreateInstanceBodyVolumes {
	newVolumes := make([]models.CreateInstanceBodyVolumes, 0, len(vSchemas))

	// convert schema volume to model
	for i := range newVolumes {
		newVolumes = append(newVolumes, models.CreateInstanceBodyVolumes{
			ID:          -1,
			Size:        vSchemas[i]["size"].(int),
			DatastoreID: vSchemas[i]["datastore_id"],
		})
	}

	// check parent instance have same volume name, if so use same id
	for _, VModel := range vModels {
		volumeExist := false
		for i, v := range newVolumes {
			if VModel.Name == v.Name {
				newVolumes[i].ID = VModel.ID
				volumeExist = true
			}
		}
		// if parent instance volume not exist in schema add it in request
		if !volumeExist {
			newVolumes = append(newVolumes, models.CreateInstanceBodyVolumes{
				ID:          VModel.ID,
				Size:        VModel.Size,
				DatastoreID: VModel.DatastoreID,
				Name:        VModel.Name,
			})
		}
	}
	newVolumes[0].RootVolume = true

	return newVolumes
}

func createInstanceSnapshot(ctx context.Context, iclient iClient, meta interface{}, instanceID int, snapshot models.SnapshotBody) error {
	snapshotResponse, err := utils.Retry(ctx, meta, func(ctx context.Context) (interface{}, error) {
		return iclient.getIClient().SnapshotAnInstance(ctx, instanceID, &snapshot)
	})
	if err != nil {
		return err
	}
	instanceModel := snapshotResponse.(models.Instances)
	if instanceModel.Success {
		return fmt.Errorf("%s", "failed to create snapshot with status as false")
	}

	return nil
func instanceSetIP(d *utils.Data, instance models.GetInstanceResponse) {
	ip := make([]string, len(instance.Instance.ConnectionInfo))
	for i := range instance.Instance.ConnectionInfo {
		ip[i] = instance.Instance.ConnectionInfo[i].IP
	}
	d.Set("ip", ip)
}
