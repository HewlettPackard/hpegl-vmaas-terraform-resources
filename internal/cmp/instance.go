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

	networks, err := getNetwork(ctx, d.Get("networks"))
	if err != nil {
		return err
	}

	volumes, err := getVolume(ctx, d.Get("volumes"))
	if err != nil {
		return err
	}

	config, err := getConfig(ctx, d.Get("config"))
	if err != nil {
		return err
	}
	tags, _ := getTags(ctx, d.Get("tags"))
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
			Layout: &models.CreateInstanceBodyInstanceLayout{
				Id: utils.JsonNumber(d.Get("layout_id")),
			},
		},
		Volumes:           volumes,
		NetworkInterfaces: networks,
		Config:            config,
		Tags:              tags,
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
	if !res.Success {
		return fmt.Errorf("%s", res.Message)
	}
	d.SetId("")
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
	d.SetId(strconv.Itoa(int(resp.Instance.Id)))
	d.Set("status", resp.Instance.Status)
	d.Set("state", "poweron")

	return nil
}

func getVolume(ctx context.Context, v interface{}) ([]models.CreateInstanceBodyVolumes, error) {
	log.Printf("[INFO] Volumes V :  %+v", v)
	volumes, err := utils.ListToMap(v)
	if err != nil {
		return nil, err
	}
	log.Printf("[INFO] Volumes :  %+v", volumes)
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

func getNetwork(ctx context.Context, v interface{}) ([]models.CreateInstanceBodyNetworkInterfaces, error) {
	networksMap, err := utils.ListToMap(v)
	if err != nil {
		return nil, err
	}
	log.Printf("[INFO] Networks :  %+v", networksMap)
	var networks []models.CreateInstanceBodyNetworkInterfaces
	for _, n := range networksMap {
		networks = append(networks, models.CreateInstanceBodyNetworkInterfaces{
			Network: &models.CreateInstanceBodyNetwork{
				Id: int32(n["id"].(int)),
			},
			NetworkInterfaceTypeId: utils.JsonNumber(n["interface_type_id"]),
		})
	}
	return networks, nil
}

func getConfig(ctx context.Context, v interface{}) (*models.CreateInstanceBodyConfig, error) {
	c, err := utils.SetToMap(v)
	if err != nil {
		return nil, err
	}
	config := &models.CreateInstanceBodyConfig{
		ResourcePoolId: utils.JsonNumber(c["resource_pool_id"]),
		Template:       int32(c["template_id"].(int)),
	}
	return config, nil
}

func getTags(ctx context.Context, v interface{}) ([]models.CreateInstanceBodyTag, error) {
	t, err := utils.MapTopMap(v)
	if err != nil {
		return nil, err
	}
	tags := make([]models.CreateInstanceBodyTag, 0, len(t))
	for k, v := range t {
		tags = append(tags, models.CreateInstanceBodyTag{
			Name:  k,
			Value: v.(string),
		})
	}
	return tags, nil
}
