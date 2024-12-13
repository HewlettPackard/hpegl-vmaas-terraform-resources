// (C) Copyright 2024 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"
	"log"
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
	layoutID := d.GetString("layout_id")
	controllerName := d.GetString("controller_name")
	busNumber := d.GetInt("bus_number")
	unitNumber := d.GetInt("interface_number")

	// Pre check
	if err := d.Error(); err != nil {
		return err
	}
	controllerName = strings.TrimSpace(strings.ToLower(controllerName))
	controllerMount, err := i.iClient.GetStorageControllerMount(ctx, layoutID, controllerName, busNumber, unitNumber)
	if err != nil {
		return err
	}
	d.SetID(controllerMount)

	return nil
}
