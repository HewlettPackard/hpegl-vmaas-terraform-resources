// // (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hpe-hcss/vmaas-cmp-go-sdk/pkg/client"
	"github.com/hpe-hcss/vmaas-cmp-go-sdk/pkg/models"
	"github.com/hpe-hcss/vmaas-terraform-resources/internal/logger"
	"github.com/hpe-hcss/vmaas-terraform-resources/internal/utils"
)

type cloud struct {
	cloudClient *client.CloudsApiService
}

func newCloud(cloudClient *client.CloudsApiService) *cloud {
	return &cloud{
		cloudClient: cloudClient,
	}
}

func (c *cloud) Read(ctx context.Context, d *utils.Data) error {
	logger.Debug("Get Cloud")

	name := d.GetString("name")
	if err := d.Error(); err != nil {
		return err
	}
	resp, err := utils.Retry(func() (interface{}, error) {
		return c.cloudClient.GetAllClouds(ctx, map[string]string{
			nameKey: name,
		})
	})
	if err != nil {
		return err
	}
	cloud := resp.(models.CloudsResp)
	if len(cloud.Clouds) != 1 {
		return fmt.Errorf(errExactMatch, "clouds")
	}
	d.SetID(strconv.Itoa(cloud.Clouds[0].ID))

	// post check
	return d.Error()
}
