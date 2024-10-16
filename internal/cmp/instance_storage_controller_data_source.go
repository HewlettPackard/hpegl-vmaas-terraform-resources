// (C) Copyright 2024 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/client"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/utils"
)

type instanceStorageController struct {
	iClient *client.InstancesAPIService
}

func newInstanceStorageController(instanceClient *client.InstancesAPIService) *instanceStorageController {
	return &instanceStorageController{iClient: instanceClient}
}

func (i *instanceStorageController) Read(ctx context.Context, d *utils.Data, meta interface{}) error {
	setMeta(meta, i.iClient.Client)
	log.Printf("[DEBUG] Get Instance Storage Controller")
	instanceID := d.GetInt("instance_id")
	if instanceID == 0 {
		return nil
	}
	controllerType := d.GetString("controller_type")
	busNumber := d.GetInt("bus_number")
	unitNumber := d.GetInt("interface_number")

	// Pre check
	if err := d.Error(); err != nil {
		return err
	}
	controllerType = strings.TrimSpace(strings.ToLower(controllerType))
	supportedControllerType := []string{"ide", "scsi"}
	if !slices.Contains(supportedControllerType, controllerType) {
		err := fmt.Errorf("storage controller '%s' is not supported", controllerType)
		return err
	}
	controllerMount, err := i.iClient.GetStorageControllerMount(ctx, instanceID, controllerType, busNumber, unitNumber)
	if err != nil {
		return err
	}
	d.SetID(controllerMount)

	return nil
}
