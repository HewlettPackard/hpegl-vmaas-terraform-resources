// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"
	"errors"
	"strings"

	"github.com/hpe-hcss/vmaas-cmp-go-sdk/pkg/client"
	"github.com/hpe-hcss/vmaas-cmp-go-sdk/pkg/models"
	"github.com/hpe-hcss/vmaas-terraform-resources/internal/logger"
	"github.com/hpe-hcss/vmaas-terraform-resources/internal/utils"
)

// instance implements functions related to cmp instances
type instance struct {
	// expose Instance API service to instances related operations
	iClient *client.InstancesAPIService
}

func (i *instance) getIClient() *client.InstancesAPIService {
	return i.iClient
}

func newInstance(iClient *client.InstancesAPIService) *instance {
	return &instance{
		iClient: iClient,
	}
}

// Create instance
func (i *instance) Create(ctx context.Context, d *utils.Data, meta interface{}) error {
	logger.Debug("Creating new instance")

	err := instanceValidateVolumeNameIsUnique(d.GetListMap("volume"))
	if err != nil {
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
		Ports:             instanceGetPorts(d.GetListMap("port")),
		Evars:             instanceGetEvars(d.GetMap("evars")),
		Labels:            d.GetStringList("labels"),
		Volumes:           instanceGetVolume(d.GetListMap("volume")),
		NetworkInterfaces: instanceGetNetwork(d.GetListMap("network")),
		Config:            instanceGetConfig(c),
		Tags:              instanceGetTags(d.GetMap("tags")),
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
	// Pre check
	if err := d.Error(); err != nil {
		return err
	}

	// create instance
	respVM, err := utils.Retry(ctx, meta, func(ctx context.Context) (interface{}, error) {
		return i.iClient.CreateAnInstance(ctx, req)
	})
	if err != nil {
		return err
	}
	getInstanceBody := *respVM.(models.GetInstanceResponse).Instance

	// Upon creation instance will be in poweron state. Check any other
	// power state provided and do accordingly
	err = instanceValidatePower(d.GetString("power"))
	if err != nil {
		return err
	}

	d.SetID(getInstanceBody.ID)

	// post check
	return d.Error()
}

func (i *instance) Update(ctx context.Context, d *utils.Data, meta interface{}) error {
	return updateInstance(ctx, i, d, meta)
}

func (i *instance) Delete(ctx context.Context, d *utils.Data, meta interface{}) error {
	return deleteInstance(ctx, i, d, meta)
}

func (i *instance) Read(ctx context.Context, d *utils.Data, meta interface{}) error {
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
	volumeLen := len(volumes)
	if volumeLen > len(instance.Instance.Volumes) {
		volumeLen = len(instance.Instance.Volumes)
	}
	for i := 0; i < volumeLen; i++ {
		volumes[i]["id"] = instance.Instance.Volumes[i].ID
	}
	d.Set("volume", volumes)
	instanceSetIP(ctx, d, instance)
	d.SetID(instance.Instance.ID)
	d.SetString("status", instance.Instance.Status)

	// post check
	return d.Error()
}
