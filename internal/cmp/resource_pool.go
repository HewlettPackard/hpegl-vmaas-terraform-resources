// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"
	"errors"
	"strconv"

	"github.com/hpe-hcss/vmaas-cmp-go-sdk/pkg/client"
	"github.com/hpe-hcss/vmaas-terraform-resources/internal/logger"
	"github.com/hpe-hcss/vmaas-terraform-resources/internal/utils"
)

type resourcePool struct {
	nClient           *client.CloudsApiService
	serviceInstanceID string
}

func newResourcePool(nClient *client.CloudsApiService, serviceInstanceID string) *resourcePool {
	return &resourcePool{nClient: nClient, serviceInstanceID: serviceInstanceID}
}

func (n *resourcePool) Read(ctx context.Context, d *utils.Data) error {
	logger.Debug("Get ResourcePool")

	name := d.GetString("name")
	cloudID := d.GetInt("cloud_id")
	resourcePools, err := n.nClient.GetAllCloudResourcePools(ctx, n.serviceInstanceID, int(cloudID), map[string]string{
		"name": name,
	})
	if err != nil {
		return err
	}

	if len(resourcePools.ResourcePools) != 1 {
		return errors.New("error coudn't find exact resourcePool, please check the name")
	}
	d.SetID(strconv.Itoa(resourcePools.ResourcePools[0].ID))

	// post check
	if err := d.Error(); err != nil {
		return err
	}

	return nil
}
