// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	client "github.com/hpe-hcss/vmaas-cmp-go-sdk/pkg/client"
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
	groupID, err := strconv.Atoi(d.Get("group_id").(string))
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
	// config, err := utils.SetToMap(d.Get("config").([]interface{}))
	log.Println("[DEBUG] config =  ", d.Get("config"))
	req := &models.CreateInstanceBody{
		ZoneId: utils.JsonNumber(d.Get("cloud_id")),
		Instance: &models.CreateInstanceBodyInstance{
			Name:  d.Get("name").(string),
			Cloud: "HPE GreenLake VMaaS Cloud",
			InstanceType: &models.CreateInstanceBodyInstanceInstanceType{
				Code: d.Get("instance_code").(string),
			},
			Plan: &models.CreateInstanceBodyInstancePlan{
				Id: utils.JsonNumber(d.Get("plan_id")),
			},
			Site: &models.CreateInstanceBodyInstanceSite{
				Id: int32(groupID),
			},
			Type: d.Get("instance_code").(string),
			Layout: &models.CreateInstanceBodyInstanceLayout{
				Id: utils.JsonNumber(d.Get("layout_id")),
			},
		},
		Volumes:           volumes,
		NetworkInterfaces: networks,
		Config: &models.CreateInstanceBodyConfig{
			ResourcePoolId: "2",
		},
	}

	resp, err := i.iClient.CreateAnInstance(ctx, i.serviceInstance, req)
	if err != nil {
		return err
	}
	d.SetId(strconv.Itoa(int(resp.Instance.Id)))

	return nil
}

// Update instance including poweroff, powerOn, restart, suspend
// changing network, volumes and instance properies such as labels
// groups and tags
func (i *instance) Update(ctx context.Context, d *schema.ResourceData) error {
	return nil
}

// Delete instance and set ID as ""
func (i *instance) Delete(ctx context.Context, d *schema.ResourceData) error {
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	res, err := i.iClient.DeleteAnInstance(ctx, i.serviceInstance, int32(id))
	if err != nil {
		return err
	}
	if res.Success {
		d.SetId("")
	} else {
		return fmt.Errorf("%s", res.Message)
	}
	return nil
}

// Read instance and set state values accordingly
func (i *instance) Read(ctx context.Context, d *schema.ResourceData) error {
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	resp, err := i.iClient.GetASpecificInstance(ctx, i.serviceInstance, int32(id))
	if err != nil {
		return err
	}
	d.Set("name", resp.Instance.Name)
	d.Set("layout", resp.Instance.Cloud.Id)
	d.Set("cloud_id", resp.Instance.Cloud.Id)
	d.Set("group_id", resp.Instance.Group.Id)
	d.Set("plan_id", resp.Instance.Plan.Id)
	d.Set("instance_type", resp.Instance.InstanceType)
	d.Set("networks", resp.Instance)
	d.Set("volumes", resp.Instance.Volumes)
	d.Set("size", resp.Instance.Volumes[0].Size)
	d.Set("datastore_id", resp.Instance.Volumes[0].DatastoreId)
	d.Set("labels", resp.Instance.Labels)
	d.Set("tags", resp.Instance.Tags)
	d.Set("config", resp.Instance.Config)
	d.Set("vmware_resource_pool", resp.Instance.Config.ResourcePoolID)
	d.Set("public_key", resp.Instance.Config)
	d.Set("copies", resp.Instance)
	d.Set("evars", resp.Instance.Evars)
	return nil
}

func getVolume(v interface{}) ([]models.CreateInstanceBodyVolumes, error) {
	log.Printf("\n[INFO] Volumes V :  %+v", v)
	volumes, err := utils.ListToMap(v)
	if err != nil {
		return nil, err
	}
	log.Printf("\n[INFO] Volumes :  %+v", volumes)
	var volumesModel []models.CreateInstanceBodyVolumes
	for i := range volumes {
		vID, err := strconv.Atoi(volumes[i]["size"].(string))
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
	log.Printf("\n[INFO] Networks :  %+v", networksMap)
	var networks []models.CreateInstanceBodyNetworkInterfaces
	for _, n := range networksMap {
		networks = append(networks, models.CreateInstanceBodyNetworkInterfaces{
			Network: &models.CreateInstanceBodyNetwork{
				Id: int32(n["network_id"].(int)),
			},
			NetworkInterfaceTypeId: utils.JsonNumber(n["interface_type_id"]),
		})
	}
	return networks, nil
}
