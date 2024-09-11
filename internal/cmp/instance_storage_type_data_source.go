// (C) Copyright 2024 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/client"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/utils"
	"github.com/tshihad/tftags"
)

type instanceStorageType struct {
	iClient *client.InstancesAPIService
}

func newInstanceStorageType(instanceClient *client.InstancesAPIService) *instanceStorageType {
	return &instanceStorageType{iClient: instanceClient}
}

func (i *instanceStorageType) Read(ctx context.Context, d *utils.Data, meta interface{}) error {
	setMeta(meta, i.iClient.Client)
	log.Printf("[DEBUG] Get Instance Storage Volume Type")
	name := d.GetString("name")
	name = strings.ToLower(name)
	name = strings.TrimSpace(name)
	cloudID := d.GetString("cloud_id")
	layoutID := d.GetString("layout_id")

	// Pre check
	if err := d.Error(); err != nil {
		return err
	}

	storageType, err := i.iClient.GetStorageVolTypeID(ctx, cloudID, layoutID)
	if err != nil {
		return err
	}

	for _, n := range storageType.Plans[0].StorageTypes {
		volType := strings.ToLower(n.Name)
		volType = strings.TrimSpace(volType)
		if volType == name {
			log.Print("[DEBUG] Storage type ID ", n.ID)

			return tftags.Set(d, n)
		}
	}

	return fmt.Errorf(errExactMatch, "Instance Storage Type")
}
