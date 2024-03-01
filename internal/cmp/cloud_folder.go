// (C) Copyright 2021-2024 Hewlett Packard Enterprise Development LP

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
	param := map[string]string{maxKey: "2000"} // There could be many folders, and max=-1 doesn't return any data
	folders, err := f.fClient.GetAllCloudFolders(ctx, cloudID, param)
	if err != nil {
		return err
	}

	for _, cf := range folders.Folders {
		if cf.Name == name {
			err = d.Set("code", cf.ExternalID)
			if err != nil {
				return err
			}
			d.SetID(cf.ID)

			return nil
		}
	}

	return fmt.Errorf(errExactMatch, "Folder")
}
