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

type resourcePool struct {
	rClient *client.CloudsAPIService
}

func newResourcePool(rClient *client.CloudsAPIService) *resourcePool {
	return &resourcePool{rClient: rClient}
}

func (n *resourcePool) Read(ctx context.Context, d *utils.Data, meta interface{}) error {
	logger.Debug("Get ResourcePool")

	name := d.GetString("name")
	cloudID := d.GetInt("cloud_id")
	// Pre check
	if err := d.Error(); err != nil {
		return err
	}

	flag := false
	resp, err := utils.Retry(ctx, meta, func(ctx context.Context) (interface{}, error) {
		return n.rClient.GetAllCloudResourcePools(ctx, cloudID, map[string]string{
			maxKey: "100",
		})
	})
	if err != nil {
		return err
	}

	resourcePools := resp.(models.ResourcePoolsResp)
	for i, r := range resourcePools.ResourcePools {
		if r.Name == name {
			flag = true
			d.SetID(resourcePools.ResourcePools[i].ID)

			break
		}
	}
	if !flag {
		return fmt.Errorf(errExactMatch, "resource pool")
	}

	// post check
	return d.Error()
}
