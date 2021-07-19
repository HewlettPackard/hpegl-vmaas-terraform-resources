// // (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"
	"fmt"

	"github.com/hpe-hcss/vmaas-cmp-go-sdk/pkg/client"
	"github.com/hpe-hcss/vmaas-cmp-go-sdk/pkg/models"
	"github.com/hpe-hcss/vmaas-terraform-resources/internal/logger"
	"github.com/hpe-hcss/vmaas-terraform-resources/internal/utils"
)

type template struct {
	tClient *client.VirtualImagesAPIService
}

func newTemplate(tClient *client.VirtualImagesAPIService) *template {
	return &template{
		tClient: tClient,
	}
}

func (c *template) Read(ctx context.Context, d *utils.Data, meta interface{}) error {
	logger.Debug("Get Templates")

	name := d.GetString("name")
	if err := d.Error(); err != nil {
		return err
	}
	resp, err := utils.Retry(ctx, meta, func(ctx context.Context) (interface{}, error) {
		return c.tClient.GetAllVirtualImages(ctx, map[string]string{
			nameKey: name,
		})
	})
	if err != nil {
		return err
	}
	template := resp.(models.VirtualImages)
	if len(template.VirtualImages) != 1 {
		return fmt.Errorf(errExactMatch, "templates")
	}
	d.SetID(template.VirtualImages[0].ID)

	// post check
	return d.Error()
}
