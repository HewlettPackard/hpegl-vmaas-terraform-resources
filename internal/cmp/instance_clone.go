// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/client"
	"github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/models"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/utils"
)

const (
	instanceCloneRetryDelay = time.Second * 15
)

// instanceClone implements functions related to cmp instanceClones
type instanceClone struct {
	// expose Instance API service to instanceClones related operations
	instanceSharedClient
}

func newInstanceClone(iClient *client.InstancesAPIService, sClient *client.ServersAPIService) *instanceClone {
	return &instanceClone{
		instanceSharedClient: instanceSharedClient{
			iClient: iClient,
			sClient: sClient,
		},
	}
}

// Create instanceClone
func (i *instanceClone) Create(ctx context.Context, d *utils.Data, meta interface{}) error {
	setMeta(meta, i.iClient.Client)
	log.Printf("[INFO] Cloning instance")

	req := models.CreateInstanceCloneBody{
		Name:  d.GetString("name"),
		Cloud: models.IDModel{ID: d.GetInt("cloud_id")},
		Group: models.IDModel{ID: d.GetInt("group_id")},
		InstanceType: models.CreateInstanceCloneInstanceTypeBody{
			Code: d.GetString("instance_type_code"),
		},
		Instance: models.CreateInstanceCloneInstanceBody{
			EnvironmentPrefix: d.GetString("env_prefix"),
			Tags:              d.GetStringList("labels"),
			InstanceContext:   d.GetString("environment_code"),
			PowerScheduleType: d.GetInt("power_schedule_id"),
		},
		Plan:              models.IDModel{ID: d.GetInt("plan_id")},
		LayoutSize:        d.GetInt("scale"),
		NetworkInterfaces: instanceGetNetwork(d.GetListMap("network")),
		Evars:             instanceGetEvars(d.GetMap("evars")),
		Metadata:          instanceGetTags(d.GetMap("tags")),
	}

	c := d.GetListMap("config")
	if len(c) > 0 {
		req.Config = *instanceGetConfig(c[0], strings.ToLower(req.InstanceType.Code) == vmware)
	}
	// Pre check
	if err := d.Error(); err != nil {
		return err
	}

	// Get source instance
	sourceID := d.GetInt("source_instance_id")
	err := copyInstanceAttribsToClone(ctx, i, &req, d.GetListMap("volume"), sourceID)
	if err != nil {
		return err
	}

	// clone the instance
	log.Printf("[INFO] Cloning the instance with %d", sourceID)
	err = cloneInstance(ctx, i, meta, req, sourceID)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Check history")
	err = checkInstanceCloneHistory(ctx, i, meta, sourceID)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Get all instances")
	getInstanceRetry := &utils.CustomRetry{
		RetryDelay: instanceCloneRetryDelay,
		Timeout:    time.Minute * 2,
		Cond: func(resp interface{}, err error) (bool, error) {
			if err != nil {
				return false, nil
			}
			instancesList := resp.(models.Instances)

			return len(instancesList.Instances) == 1, nil
		},
	}
	// get cloned instance ID
	instancesResp, err := getInstanceRetry.Retry(ctx, meta, func(ctx context.Context) (interface{}, error) {
		return i.iClient.GetAllInstances(ctx, map[string]string{
			nameKey: req.Name,
		})
	})
	if err != nil {
		return err
	}

	instancesList := instancesResp.(models.Instances)
	if len(instancesList.Instances) != 1 {
		return errors.New("get cloned instance is failed")
	}

	if err := instanceWaitUntilCreated(ctx, i.instanceSharedClient, meta, instancesList.Instances[0].ID); err != nil {
		return err
	}

	if snapshot := d.GetListMap("snapshot"); len(snapshot) == 1 {
		err := createInstanceSnapshot(ctx, i.instanceSharedClient, instancesList.Instances[0].ID, models.SnapshotBody{
			Snapshot: &models.SnapshotBodySnapshot{
				Name:        snapshot[0]["name"].(string),
				Description: snapshot[0]["description"].(string),
			},
		})
		if err != nil {
			return err
		}
	}
	err = instanceSetServerID(ctx, d, i.instanceSharedClient)
	if err != nil {
		return err
	}
	d.SetID(instancesList.Instances[0].ID)

	// post check
	return d.Error()
}

// Update instance including poweroff, powerOn, restart, suspend
// changing volumes and instance properties such as labels
// groups and tags
func (i *instanceClone) Update(ctx context.Context, d *utils.Data, meta interface{}) error {
	setMeta(meta, i.iClient.Client)

	return updateInstance(ctx, i.instanceSharedClient, d)
}

// Delete instance and set ID as ""
func (i *instanceClone) Delete(ctx context.Context, d *utils.Data, meta interface{}) error {
	setMeta(meta, i.iClient.Client)

	return deleteInstance(ctx, i.instanceSharedClient, d, meta)
}

// Read instance and set state values accordingly
func (i *instanceClone) Read(ctx context.Context, d *utils.Data, meta interface{}) error {
	setMeta(meta, i.iClient.Client)

	return readInstance(ctx, i.instanceSharedClient, d, meta, true)
}

func checkInstanceCloneHistory(ctx context.Context, i *instanceClone, meta interface{}, instanceID int) error {
	errCount := 0
	historyRetry := utils.CustomRetry{
		InitialDelay: time.Second * 15,
		RetryDelay:   time.Second * 30,
		Timeout:      maxTimeout,
		Cond: func(response interface{}, ResponseErr error) (bool, error) {
			if ResponseErr != nil {
				errCount++
				if errCount == 3 {
					return false, ResponseErr
				}

				return false, nil
			}
			errCount = 0

			instanceHistory := response.(models.GetInstanceHistory)
			for _, processes := range instanceHistory.Processes {
				if processes.ProcessType.Code == "cloning" {
					if processes.Status == "success" || processes.Status == "complete" {
						return true, nil
					}
					if processes.Status == "failed" {
						return false, fmt.Errorf("failed to clone instance")
					}

					break
				}
			}

			return false, nil
		},
	}
	_, err := historyRetry.Retry(ctx, meta, func(ctx context.Context) (interface{}, error) {
		return i.iClient.GetInstanceHistory(ctx, instanceID)
	})

	return err
}

func cloneInstance(
	ctx context.Context,
	i *instanceClone,
	meta interface{},
	req models.CreateInstanceCloneBody,
	sourceID int,
) error {
	cloneRetry := &utils.CustomRetry{
		Cond: func(response interface{}, ResponseErr error) (bool, error) {
			if ResponseErr != nil {
				return false, nil
			}
			if !response.(models.SuccessOrErrorMessage).Success {
				return false, fmt.Errorf("failed to clone instance")
			}

			return true, nil
		},
	}
	_, err := cloneRetry.Retry(ctx, meta, func(ctx context.Context) (interface{}, error) {
		val, err := json.Marshal(&req)
		if err != nil {
			return nil, err
		}
		log.Printf("value: %s", string(val))

		return i.iClient.CloneAnInstance(ctx, sourceID, req)
	})

	return err
}

func copyInstanceAttribsToClone(
	ctx context.Context,
	i *instanceClone,
	req *models.CreateInstanceCloneBody,
	volumes []map[string]interface{},
	sourceID int,
) error {
	sourceInstance, err := i.iClient.GetASpecificInstance(ctx, sourceID)
	if err != nil {
		return err
	}

	if utils.IsEmpty(req.Cloud.ID) {
		req.Cloud.ID = sourceInstance.Instance.Cloud.ID
	}
	if utils.IsEmpty(req.Plan.ID) {
		req.Plan.ID = sourceInstance.Instance.Plan.ID
	}
	if utils.IsEmpty(req.InstanceType.Code) {
		req.InstanceType.Code = sourceInstance.Instance.InstanceType.Code
	}
	if utils.IsEmpty(req.Group.ID) {
		req.Group.ID = sourceInstance.Instance.Group.ID
	}
	if utils.IsEmpty(req.LayoutSize) {
		req.LayoutSize = sourceInstance.Instance.Config.Layoutsize
	}
	if utils.IsEmpty(req.Instance.Tags) {
		req.Instance.Tags = sourceInstance.Instance.Labels
	}
	if utils.IsEmpty(req.Metadata) {
		req.Metadata = sourceInstance.Instance.Tags
	}
	if utils.IsEmpty(req.Instance.PowerScheduleType) {
		req.Instance.PowerScheduleType = sourceInstance.Instance.Config.PowerScheduleType
	}
	if utils.IsEmpty(req.Instance.EnvironmentPrefix) {
		req.Instance.EnvironmentPrefix = sourceInstance.Instance.EnvironmentPrefix
	}
	if utils.IsEmpty(req.Instance.InstanceContext) {
		req.Instance.InstanceContext = sourceInstance.Instance.InstanceContext
	}
	if utils.IsEmpty(req.LayoutSize) {
		req.LayoutSize = sourceInstance.Instance.Config.Layoutsize
	}
	instanceCloneCopyConfig(req, sourceInstance)

	req.Volumes = instanceCloneCompareVolume(volumes, sourceInstance.Instance.Volumes)
	req.Layout = models.IDModel{
		ID: sourceInstance.Instance.Layout.ID,
	}

	return nil
}

func instanceCloneCopyConfig(req *models.CreateInstanceCloneBody, sourceInstance models.GetInstanceResponse) {
	if utils.IsEmpty(req.Config.ResourcePoolID) {
		req.Config.ResourcePoolID = sourceInstance.Instance.Config.ResourcePoolID
	}
	if utils.IsEmpty(req.Config.Template) {
		req.Config.Template = sourceInstance.Instance.Config.Template
	}
	if utils.IsEmpty(req.Config.VMwareFolderID) {
		req.Config.VMwareFolderID = sourceInstance.Instance.Config.Vmwarefolderid
	}
	if utils.IsEmpty(req.Config.SmbiosAssetTag) {
		req.Config.SmbiosAssetTag = sourceInstance.Instance.Config.Smbiosassettag
	}
}
