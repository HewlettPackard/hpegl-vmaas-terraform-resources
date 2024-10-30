// // (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"
	"fmt"
	"log"

	"github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/client"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/utils"
)

type cloud struct {
	cloudClient *client.CloudsAPIService
}

func newCloud(cloudClient *client.CloudsAPIService) *cloud {
	return &cloud{
		cloudClient: cloudClient,
	}
}

func (c *cloud) Read(ctx context.Context, d *utils.Data, meta interface{}) error {
	setMeta(meta, c.cloudClient.Client)
	log.Printf("[INFO] Get Cloud")

	name := d.GetString("name")
	if err := d.Error(); err != nil {
		return err
	}
	cloud, err := c.cloudClient.GetAllClouds(context.WithoutCancel(ctx), map[string]string{
		nameKey: name,
	})
	if err != nil {
		return err
	}
	if len(cloud.Clouds) != 1 {
		return fmt.Errorf(errExactMatch, "clouds")
	}
	d.SetID(cloud.Clouds[0].ID)

	// post check
	return d.Error()
}
