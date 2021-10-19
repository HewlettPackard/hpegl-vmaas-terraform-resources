// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"
	"fmt"
	"log"

	"github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/client"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/utils"
)

type layout struct {
	gClient *client.LibraryAPIService
}

func newLayout(gClient *client.LibraryAPIService) *layout {
	return &layout{
		gClient: gClient,
	}
}

func (g *layout) Read(ctx context.Context, d *utils.Data, meta interface{}) error {
	log.Printf("[DEBUG] Get Layout")

	name := d.GetString("name")
	instanceTypeCode := d.GetString("instance_type_code")
	// Pre check
	if err := d.Error(); err != nil {
		return err
	}
	instanceTypes, err := g.gClient.GetAllInstanceTypes(ctx, map[string]string{
		codeKey:          instanceTypeCode,
		provisionTypeKey: vmware,
	})
	if err != nil {
		return err
	}

	if len(instanceTypes.InstanceTypes) != 1 {
		return fmt.Errorf(errExactMatch, "instance type")
	}
	for _, l := range instanceTypes.InstanceTypes[0].Instancetypelayouts {
		if l.Name == name {
			d.SetID(l.ID)

			return d.Error()
		}
	}

	return fmt.Errorf(errExactMatch, "layout")
}
