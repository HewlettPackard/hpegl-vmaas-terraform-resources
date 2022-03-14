// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/client"
	"github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/models"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/utils"
	pkgUtils "github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/utils"
	"github.com/tshihad/tftags"
)

type instanceSharedClient struct {
	iClient *client.InstancesAPIService
	sClient *client.ServersAPIService
}

func readInstance(ctx context.Context, sharedClient instanceSharedClient, d *utils.Data, meta interface{}, isClone bool) error {
	id := d.GetID()

	log.Printf("[INFO] Get instance with ID %d", id)
	// Precheck
	if err := d.Error(); err != nil {
		return err
	}

	instance, err := sharedClient.iClient.GetASpecificInstance(ctx, id)
	if err != nil {
		return err
	}

	tfInstance := models.TFInstance{}
	if err := tftags.Get(d, &tfInstance); err != nil {
		return err
	}

	// Assign proper ID for the volume, since response may contains more
	// volumes than schema, check the name and assign id
	for i := range tfInstance.Volume {
		for _, vModel := range instance.Instance.Volumes {
			if vModel.Name == tfInstance.Volume[i].Name {
				tfInstance.Volume[i].ID = vModel.ID

				break
			}
		}
	}
	// Invoke all API request in parallel
	// Get server details
	serverRetry := &utils.CustomRetry{}
	serverRetry.RetryParallel(ctx, meta, func(ctx context.Context) (interface{}, error) {
		return sharedClient.sClient.GetSpecificServer(ctx, d.GetInt("server_id"))
	})
	// get snapshot details
	snapshotRetry := &utils.CustomRetry{}
	snapshotRetry.RetryParallel(ctx, meta, func(ctx context.Context) (interface{}, error) {
		return sharedClient.iClient.GetListOfSnapshotsForAnInstance(ctx, instance.Instance.ID)
	})
	// Get history details
	historyRetry := &utils.CustomRetry{}
	historyRetry.RetryParallel(ctx, meta, func(ctx context.Context) (interface{}, error) {
		return sharedClient.iClient.GetInstanceHistory(ctx, instance.Instance.ID)
	})

	if isClone {
		err = d.Set("layout_id", instance.Instance.Layout.ID)
		if err != nil {
			return err
		}
	}

	tfInstance.Network, err = instanceGetNetworkModel(tfInstance.Network, serverRetry)
	if err != nil {
		return err
	}

	tfInstance.Status = instance.Instance.Status
	tfInstance.Snapshot = instanceGetSnaphotModel(tfInstance.Snapshot, snapshotRetry)
	tfInstance.History = instanceGetHistoryModel(historyRetry)
	tfInstance.Containers = instance.Instance.ContainerDetails

	err = tftags.Set(d, tfInstance)
	if err != nil {
		return err
	}

	d.SetID(instance.Instance.ID)

	// post check
	return d.Error()
}

// Update instance including poweroff, powerOn, restart, suspend
// changing volumes and instance properties such as labels
// groups and tags
func updateInstance(ctx context.Context, sharedClient instanceSharedClient, d *utils.Data) error {
	log.Printf("[DEBUG] Updating the instance")

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
		_, err := sharedClient.iClient.UpdatingAnInstance(ctx, id, updateReq)
		if err != nil {
			return err
		}
	}
	if err := instanceUpdateNetworkVolumePlan(ctx, sharedClient, d, id); err != nil {
		return err
	}

	getInstance, err := sharedClient.iClient.GetASpecificInstance(ctx, id)
	if err != nil {
		return err
	}
	if d.HasChanged("power") || d.HasChanged("restart_instance") {
		// Do power operation only if backend is in different state
		// restart only if instance in actual is in power-on state
		status := utils.ParsePowerState(getInstance.Instance.Status)
		powerOp := d.GetString("power")
		if powerOp != status {
			if err := instanceDoPowerTask(ctx, sharedClient, id, d.GetString("power")); err != nil {
				return err
			}
		} else if d.HasChanged("restart_instance") {
			if err := instanceDoPowerTask(ctx, sharedClient, id, utils.Restart); err != nil {
				return err
			}
		}
	}

	if d.HasChanged("snapshot") {
		snapshot := d.GetListMap("snapshot")
		err := createInstanceSnapshot(ctx, sharedClient, getInstance.Instance.ID, models.SnapshotBody{
			Snapshot: &models.SnapshotBodySnapshot{
				Name:        snapshot[0]["name"].(string),
				Description: snapshot[0]["description"].(string),
			},
		})
		if err != nil {
			return err
		}
	}

	return d.Error()
}

// Delete instance and set ID as ""
func deleteInstance(ctx context.Context, sharedClient instanceSharedClient, d *utils.Data, meta interface{}) error {
	id := d.GetID()
	log.Printf("[DEBUG] Deleting instance with ID : %d", id)

	// Precheck
	if err := d.Error(); err != nil {
		return err
	}

	deleResp, err := sharedClient.iClient.DeleteAnInstance(ctx, id)
	if err != nil {
		return err
	}
	if !deleResp.Success {
		return fmt.Errorf("failed to delete instance with error: %s", deleResp.Message)
	}

	errCount := 0
	cRetry := utils.CustomRetry{
		RetryDelay: time.Second * 15,
		Timeout:    maxTimeout,
		Cond: func(response interface{}, ResponseErr error) (bool, error) {
			if ResponseErr != nil {
				if pkgUtils.GetStatusCode(ResponseErr) == http.StatusNotFound {
					return true, nil
				}
				errCount++
				if errCount == 3 {
					return false, err
				}
			}

			return false, nil
		},
	}
	_, err = cRetry.Retry(ctx, meta, func(ctx context.Context) (interface{}, error) {
		return sharedClient.iClient.GetASpecificInstance(ctx, id)
	})

	// post check
	return d.Error()
}

func instanceGetVolume(volumes []map[string]interface{}) []models.CreateInstanceBodyVolumes {
	volumesModel := make([]models.CreateInstanceBodyVolumes, 0, len(volumes))
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

func instanceGetConfig(c map[string]interface{}, isVmware bool) *models.CreateInstanceBodyConfig {
	config := &models.CreateInstanceBodyConfig{
		ResourcePoolID: utils.JSONNumber(c["resource_pool_id"]),
		NoAgent:        strconv.FormatBool(c["no_agent"].(bool)),
		SmbiosAssetTag: c["asset_tag"].(string),
		VMwareFolderID: c["folder_code"].(string),
		Template:       c["template_id"].(int),
		CreateUser:     c["create_user"].(bool),
	}
	if !isVmware {
		config.Template = 0
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
	sharedClient instanceSharedClient,
	instanceID int,
	newOp string) error {
	var err error

	switch newOp {
	case utils.PowerOn:
		_, err = sharedClient.iClient.StartAnInstance(ctx, instanceID)
	case utils.PowerOff:
		_, err = sharedClient.iClient.StopAnInstance(ctx, instanceID)
	case utils.Suspend:
		_, err = sharedClient.iClient.SuspendAnInstance(ctx, instanceID)
	case utils.Restart:
		_, err = sharedClient.iClient.RestartAnInstance(ctx, instanceID)
	}

	return err
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

func createInstanceSnapshot(
	ctx context.Context,
	sharedClient instanceSharedClient,
	instanceID int,
	snapshot models.SnapshotBody,
) error {
	instanceModel, err := sharedClient.iClient.SnapshotAnInstance(ctx, instanceID, &snapshot)
	if err != nil {
		return err
	}
	if !instanceModel.Success {
		return fmt.Errorf("%s", "failed to create snapshot, API returns status as false")
	}

	return nil
}

func instanceGetSnaphotModel(snapshot models.TFInstanceSnapshot, retry *utils.CustomRetry) models.TFInstanceSnapshot {
	if utils.IsEmpty(snapshot) {
		return snapshot
	}

	snaphostResp, err := retry.Wait()
	if err != nil {
		if pkgUtils.GetStatusCode(err) != http.StatusNotFound {
			return snapshot
		}
		snapshot.IsSnapshotExists = false

		return snapshot
	}
	id := instanceCheckSnaphotByName(snapshot.Name, snaphostResp)
	snapshot.ID = id
	snapshot.IsSnapshotExists = !(id == -1)

	return snapshot
}

func instanceCheckSnaphotByName(name string, snapshotResp interface{}) int {
	snapshots := snapshotResp.(models.ListSnapshotResponse)
	for _, snapshot := range snapshots.Snapshots {
		if snapshot.Name == name {
			return snapshot.ID
		}
	}

	return -1
}

func instanceWaitUntilCreated(ctx context.Context, sharedClient instanceSharedClient, meta interface{}, instanceID int) error {
	errCount := 0
	cRetry := utils.CustomRetry{
		Timeout:      maxTimeout,
		RetryDelay:   time.Second * 15,
		InitialDelay: time.Minute,
		Cond: func(response interface{}, err error) (bool, error) {
			if err != nil {
				errCount++
				// return false as condition if same error returns 3 times.
				if errCount == 3 {
					return false, err
				}

				return false, nil
			}

			instance, ok := response.(models.GetInstanceResponse)
			if !ok {
				errCount++
				if errCount == 3 {
					return false, fmt.Errorf("%s", "error while getting instance")
				}

				return false, nil
			}
			errCount = 0

			if instance.Instance.Status == utils.StateFailed ||
				instance.Instance.Status == utils.StateRunning {
				return true, nil
			}

			return false, nil
		},
	}

	_, err := cRetry.Retry(ctx, meta, func(ctx context.Context) (interface{}, error) {
		return sharedClient.iClient.GetASpecificInstance(ctx, instanceID)
	})
	if err != nil {
		return err
	}

	return nil
}

func instanceGetHistoryModel(retry *utils.CustomRetry) []models.GetInstanceHistoryProcesses {
	resp, err := retry.Wait()
	if err != nil {
		log.Printf("[WARN] Failed to retrieve the history with error %v", err)

		return nil
	}
	historyModel := resp.(models.GetInstanceHistory)

	return historyModel.Processes
}

func instanceGetResizeNetwork(network []map[string]interface{}) []models.CreateInstanceBodyNetworkInterfaces {
	nics := make([]models.CreateInstanceBodyNetworkInterfaces, 0, len(network))
	for _, n := range network {
		nics = append(nics, models.CreateInstanceBodyNetworkInterfaces{
			Name: n["name"].(string),
			ID:   n["internal_id"].(int),
			Network: &models.CreateInstanceBodyNetwork{
				ID: n["id"].(int),
			},
			NetworkInterfaceTypeID: utils.JSONNumber(n["interface_id"]),
		})
	}

	return nics
}

func instanceSetServerID(ctx context.Context, d *utils.Data, sharedClient instanceSharedClient) error {
	servers, err := sharedClient.sClient.GetAllServers(ctx, map[string]string{
		externalNameKey: d.GetString("name"),
	})
	if err != nil {
		return err
	}
	if len(servers.Server) != 1 {
		return fmt.Errorf(errExactMatch, "server")
	}
	err = d.Set("server_id", servers.Server[0].ID)
	if err != nil {
		return err
	}

	return nil
}

func instanceGetNetworkModel(
	networks []models.TFInstanceNetwork, retry *utils.CustomRetry) ([]models.TFInstanceNetwork, error) {
	resp, err := retry.Wait()
	if err != nil {
		return nil, err
	}
	serverInterface := resp.(models.GetSpecificServerResponse).Server.Interfaces
	if len(serverInterface) != len(networks) {
		return nil, fmt.Errorf("failed to set network. There is mismatch on created network and the terraform state")
	}
	for i, s := range serverInterface {
		networks[i].InternalID = s.ID
		networks[i].IsPrimary = s.PrimaryInterface
		networks[i].Name = s.Name
	}

	return networks, nil
}

func instanceUpdateNetworkVolumePlan(
	ctx context.Context,
	sharedClient instanceSharedClient,
	d *utils.Data,
	instanceID int,
) error {
	var resizeReq models.ResizeInstanceBody
	if d.HasChanged("volume") {
		volumes := instanceCompareVolumes(d.GetChangedListMap("volume"))
		resizeReq = models.ResizeInstanceBody{
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
	} else if d.HasChanged("plan_id") {
		resizeReq = models.ResizeInstanceBody{
			Instance: &models.ResizeInstanceBodyInstance{
				Plan: &models.ResizeInstanceBodyInstancePlan{
					ID: d.GetInt("plan_id"),
				},
			},
		}
	}

	if d.HasChanged("network") {
		schemaNetwork := d.GetListMap("network")
		resizeReq.NetworkInterfaces = instanceGetResizeNetwork(schemaNetwork)
	}
	if d.HasChanged("volume") || d.HasChanged("network") || d.HasChanged("plan_id") {
		updateResp, err := sharedClient.iClient.ResizeAnInstance(ctx, instanceID, &resizeReq)
		if err != nil {
			return err
		}
		if !updateResp.Success {
			return fmt.Errorf("%s", "failed to resize")
		}
	}

	return nil
}
