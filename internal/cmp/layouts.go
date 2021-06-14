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
	instanceTypeCode := d.GetString("instance_type_code")
	// Pre check
	if err := d.Error(); err != nil {
		return err
	}
	resp, err := utils.Retry(func() (interface{}, error) {
		return g.gClient.GetAllInstanceTypes(ctx, map[string]string{
			codeKey:          instanceTypeCode,
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
	for _, l := range instanceTypes.InstanceTypes[0].Instancetypelayouts {
		if l.Name == name {
			d.SetID(l.ID)

			return d.Error()
		}
	}

	return fmt.Errorf(errExactMatch, "layout name")
}
