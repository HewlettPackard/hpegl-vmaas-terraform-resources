// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"
	"fmt"
	"log"

	"github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/client"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/utils"
)

type cloudFolder struct {
	fClient *client.CloudsAPIService
}

func newCloudFolder(fClient *client.CloudsAPIService) *cloudFolder {
	return &cloudFolder{
		fClient: fClient,
	}
}

func (f *cloudFolder) Read(ctx context.Context, d *utils.Data, meta interface{}) error {
	setMeta(meta, f.fClient.Client)
	log.Printf("[INFO] Get cloud folder")

	name := d.GetString("name")
	cloudID := d.GetInt("cloud_id")
	// Pre check
	if err := d.Error(); err != nil {
		return err
	}
	folders, err := f.fClient.GetAllCloudFolders(ctx, cloudID, nil)
	if err != nil {
		return err
	}

	for _, cf := range folders.Folders {
		if cf.Name == name {
			d.Set("code", cf.ExternalID)
			d.SetID(cf.ID)

			return nil
		}
	}

	return fmt.Errorf(errExactMatch, "Folder")
}
