// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"
	"fmt"
	"log"

	"github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/client"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/utils"
	"github.com/tshihad/tftags"
)

type routertier1ds struct {
	nClient *client.RouterAPIService
}

func newTier1RouterDS(nClient *client.RouterAPIService) *routertier1ds {
	return &routertier1ds{nClient: nClient}
}

func (n *routertier1ds) Read(ctx context.Context, d *utils.Data, meta interface{}) error {
	setMeta(meta, n.nClient.Client)
	log.Printf("[DEBUG] Get Tier1 Router")
	name := d.GetString("name")

	// Pre check
	if err := d.Error(); err != nil {
		return err
	}
	routers, err := n.nClient.GetAllRouter(ctx, nil)
	if err != nil {
		return err
	}

	for i, n := range routers.NetworkRouters {
		if n.Name == name {
			log.Print("[DEBUG]", routers.NetworkRouters[i].ID)

			return tftags.Set(d, routers.NetworkRouters[i])
		}
	}

	return fmt.Errorf(errExactMatch, "Tier1 Router")
}
