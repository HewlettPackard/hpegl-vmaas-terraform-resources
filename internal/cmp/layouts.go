// (C) Copyright 2021 Hewlett Packard Enterprise Development LP
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

type layout struct {
	gClient           *client.LibraryApiService
	serviceInstanceID string
}

func newLayout(gClient *client.LibraryApiService, serviceInstanceID string) *layout {
	return &layout{
		gClient:           gClient,
		serviceInstanceID: serviceInstanceID,
	}
}

func (g *layout) Read(ctx context.Context, d *utils.Data) error {
	logger.Debug("Get Layout")

	name := d.GetString("name")
	layouts, err := g.gClient.GetAllLayouts(ctx, g.serviceInstanceID, map[string]string{
		"name": name,
	})
	if err != nil {
		return err
	}

	if len(layouts.InstanceTypeLayouts) != 1 {
		return errors.New("error coudn't find exact layout, please check the name")
	}

	d.SetID(strconv.Itoa(layouts.InstanceTypeLayouts[0].ID))

	// post check
	if err := d.Error(); err != nil {
		return err
	}

	return nil
}
