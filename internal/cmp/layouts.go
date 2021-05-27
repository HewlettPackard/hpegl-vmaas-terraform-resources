// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

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

type layout struct {
	gClient *client.LibraryApiService
}

func newLayout(gClient *client.LibraryApiService) *layout {
	return &layout{
		gClient: gClient,
	}
}

func (g *layout) Read(ctx context.Context, d *utils.Data) error {
	logger.Debug("Get Layout")

	name := d.GetString("name")
	// Pre check
	if err := d.Error(); err != nil {
		return err
	}
	resp, err := utils.Retry(func() (interface{}, error) {
		return g.gClient.GetAllInstanceTypes(ctx, map[string]string{
			nameKey:          name,
			provisionTypeKey: vmware,
		})
	})
	if err != nil {
		return err
	}
	instanceTypes := resp.(models.InstanceTypesResp)

	if len(instanceTypes.InstanceTypes) != 1 {
		return fmt.Errorf(errExactMatch, "instance type")
	}
	if len(instanceTypes.InstanceTypes[0].Instancetypelayouts) != 1 {
		return fmt.Errorf(errExactMatch, "layout type")
	}
	d.SetString("instance_code", instanceTypes.InstanceTypes[0].Code)
	d.SetID(strconv.Itoa(instanceTypes.InstanceTypes[0].Instancetypelayouts[0].ID))

	// post check
	return d.Error()
}
