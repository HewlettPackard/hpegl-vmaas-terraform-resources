// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"
	"fmt"

	"github.com/hpe-hcss/vmaas-cmp-go-sdk/pkg/client"
	"github.com/hpe-hcss/vmaas-cmp-go-sdk/pkg/models"
	"github.com/hpe-hcss/vmaas-terraform-resources/internal/logger"
	"github.com/hpe-hcss/vmaas-terraform-resources/internal/utils"
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
	logger.Debug("Get cloud folder")

	name := d.GetString("name")
	cloudID := d.GetInt("cloud_id")
	// Pre check
	if err := d.Error(); err != nil {
		return err
	}
	resp, err := utils.Retry(ctx, meta, func(ctx context.Context) (interface{}, error) {
		return f.fClient.GetAllCloudFolders(ctx, cloudID, map[string]string{
			nameKey: name,
		})
	})
	folders := resp.(models.GetAllCloudFolders)
	if err != nil {
		return err
	}
	if len(folders.Folders) != 1 {
		return fmt.Errorf(errExactMatch, "folder")
	}
	d.Set("code", folders.Folders[0].ExternalID)
	d.SetID(folders.Folders[0].ID)
	// post check
	return d.Error()
}
