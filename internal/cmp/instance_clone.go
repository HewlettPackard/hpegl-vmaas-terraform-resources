// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hpe-hcss/vmaas-cmp-go-sdk/pkg/client"
	"github.com/hpe-hcss/vmaas-cmp-go-sdk/pkg/models"
	"github.com/hpe-hcss/vmaas-terraform-resources/internal/logger"
	"github.com/hpe-hcss/vmaas-terraform-resources/internal/utils"
)

const (
	instanceCloneRetryDelay   = time.Second * 60
	instanceCloneRetryTimeout = time.Second * 30
	instanceCloneRetryCount   = 20
)

// instanceClone implements functions related to cmp instanceClones
type instanceClone struct {
	// expose Instance API service to instanceClones related operations
	iClient *client.InstancesAPIService
}

func (i *instanceClone) getIClient() *client.InstancesAPIService {
	return i.iClient
}

func newInstanceClone(iClient *client.InstancesAPIService) *instanceClone {
	return &instanceClone{
		iClient: iClient,
	}
}

// Create instanceClone
func (i *instanceClone) Create(ctx context.Context, d *utils.Data, meta interface{}) error {
	logger.Debug("Cloning instance")

	volumes := d.GetListMap("volume")
	err := instanceValidateVolumeNameIsUnique(volumes)
	if err != nil {
		return err
	}

	req := &models.CreateInstanceBody{
		CloneName: d.GetString("name"),
		ZoneID:    d.GetJSONNumber("cloud_id"),
		Instance: &models.CreateInstanceBodyInstance{
			InstanceType: &models.CreateInstanceBodyInstanceInstanceType{
				Code: d.GetString("instance_type_code"),
			},
			Plan: &models.CreateInstanceBodyInstancePlan{
				ID: d.GetJSONNumber("plan_id"),
			},
			Site: &models.CreateInstanceBodyInstanceSite{
				ID: d.GetInt("group_id"),
			},
			HostName:          d.GetString("hostname"),
			EnvironmentPrefix: d.GetString("env_prefix"),
		},
		Environment:       d.GetString("environment_code"),
		Evars:             instanceGetEvars(d.GetMap("evars")),
		Labels:            d.GetStringList("labels"),
		NetworkInterfaces: instanceGetNetwork(d.GetListMap("network")),
		Tags:              instanceGetTags(d.GetMap("tags")),
		LayoutSize:        d.GetInt("scale"),
		PowerScheduleType: utils.JSONNumber(d.GetInt("power_schedule_id")),
	}

	c := d.GetListMap("config")
	if len(c) > 0 {
		req.Config = instanceGetConfig(c[0])

		// Get template id instance type is vmware
		if strings.ToLower(req.Instance.InstanceType.Code) == vmware {
			templateID := c[0]["template_id"]
			if templateID == nil {
				return errors.New("error, template id is required for vmware instance type")
			}
			req.Config.Template = templateID.(int)
		}
	}
	// Pre check
	if err := d.Error(); err != nil {
		return err
	}

	// Get source instance
	sourceID := d.GetInt("source_instance_id")
	sourceInstanceResp, err := utils.Retry(ctx, meta, func(ctx context.Context) (interface{}, error) {
		return i.iClient.GetASpecificInstance(ctx, sourceID)
	})
	if err != nil {
		return err
	}
	sourceInstance := sourceInstanceResp.(models.GetInstanceResponse)

	if utils.IsEmpty(req.ZoneID) {
		req.ZoneID = utils.JSONNumber(sourceInstance.Instance.Cloud.ID)
	}
	if utils.IsEmpty(req.Instance.Plan.ID) {
		req.Instance.Plan.ID = utils.JSONNumber(sourceInstance.Instance.Plan.ID)
	}
	if req.Instance.InstanceType.Code == "" {
		req.Instance.InstanceType.Code = sourceInstance.Instance.InstanceType.Code
	}
	if req.Instance.Site.ID == 0 {
		req.Instance.Site.ID = sourceInstance.Instance.Group.ID
	}

	req.Volumes = instanceCloneCompareVolume(volumes, sourceInstance.Instance.Volumes)
	req.Instance.Layout = &models.CreateInstanceBodyInstanceLayout{
		ID: utils.JSONNumber(sourceInstance.Instance.Layout.ID),
	}

	// clone the instance
	logger.Info("Cloning the instance with ", sourceID)
	respClone, err := utils.Retry(ctx, meta, func(ctx context.Context) (interface{}, error) {
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
	instancesResp, err := customRetry.Retry(ctx, meta, func(ctx context.Context) (interface{}, error) {
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

	// Upon creation instance will be in poweron state. Check any other
	// power state provided and do accordingly
	err = instanceValidatePower(d.GetString("power"))
	if err != nil {
		return err
	}

	if snapshotName := d.GetString("snapshot"); snapshotName != "" {
		createInstanceSnapshot(ctx, i, meta, instancesList.Instances[0].ID, models.SnapshotBody{
			Snapshot: &models.SnapshotBodySnapshot{
				Name: snapshotName,
			},
		})
	}

	d.SetID(instancesList.Instances[0].ID)

	// post check
	return d.Error()
}

// Update instance including poweroff, powerOn, restart, suspend
// changing volumes and instance properties such as labels
// groups and tags
func (i *instanceClone) Update(ctx context.Context, d *utils.Data, meta interface{}) error {
	return updateInstance(ctx, i, d, meta)
}

// Delete instance and set ID as ""
func (i *instanceClone) Delete(ctx context.Context, d *utils.Data, meta interface{}) error {
	return deleteInstance(ctx, i, d, meta)
}

// Read instance and set state values accordingly
func (i *instanceClone) Read(ctx context.Context, d *utils.Data, meta interface{}) error {
	id := d.GetID()

	logger.Debug("Get instance with ID %d", id)

	// Precheck
	if err := d.Error(); err != nil {
		return err
	}

	resp, err := utils.Retry(ctx, meta, func(ctx context.Context) (interface{}, error) {
		return i.iClient.GetASpecificInstance(ctx, id)
	})
	if err != nil {
		return err
	}
	instance := resp.(models.GetInstanceResponse)

	volumes := d.GetListMap("volume")
	// Assign proper ID for the volume, since response may contains more
	// volumes than schema, check the name and assign ip
	for i := range volumes {
		for _, vModel := range instance.Instance.Volumes {
			if vModel.Name == volumes[i]["name"].(string) {
				volumes[i]["id"] = vModel.ID
			}
		}
	}

	d.Set("volume", volumes)

	// Write IPs in to state file
	instanceSetIP(d, instance)

	d.Set("layout_id", instance.Instance.Layout.ID)
	d.SetString("status", instance.Instance.Status)
	d.SetID(instance.Instance.ID)

	// post check
	return d.Error()
}
