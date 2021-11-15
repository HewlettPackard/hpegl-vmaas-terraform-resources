// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"
	"fmt"
	"log"

	"github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/client"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/utils"
)

type resourcePool struct {
	rClient *client.CloudsAPIService
}

func newResourcePool(rClient *client.CloudsAPIService) *resourcePool {
	return &resourcePool{rClient: rClient}
}

func (n *resourcePool) Read(ctx context.Context, d *utils.Data, meta interface{}) error {
	setMeta(meta, n.rClient.Client)
	log.Printf("[DEBUG] Get ResourcePool")

	name := d.GetString("name")
	cloudID := d.GetInt("cloud_id")
	// Pre check
	if err := d.Error(); err != nil {
		return err
	}

	flag := false
	resourcePools, err := n.rClient.GetAllCloudResourcePools(ctx, cloudID, map[string]string{
		maxKey: "100",
	})
	if err != nil {
		return err
	}

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
