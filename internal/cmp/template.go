// // (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"
	"fmt"
	"log"

	"github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/client"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/utils"
)

type template struct {
	tClient *client.VirtualImagesAPIService
}

func newTemplate(tClient *client.VirtualImagesAPIService) *template {
	return &template{
		tClient: tClient,
	}
}

func (t *template) Read(ctx context.Context, d *utils.Data, meta interface{}) error {
	setMeta(meta, t.tClient.Client)
	log.Printf("[DEBUG] Get Templates")

	name := d.GetString("name")
	if err := d.Error(); err != nil {
		return err
	}
	template, err := t.tClient.GetAllVirtualImages(ctx, map[string]string{
		nameKey:       name,
		filterTypeKey: syncedTypeValue,
	})
	if err != nil {
		return err
	}
	if len(template.VirtualImages) != 1 {
		return fmt.Errorf(errExactMatch, "templates")
	}
	d.SetID(template.VirtualImages[0].ID)

	// post check
	return d.Error()
}
