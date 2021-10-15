// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"
	"fmt"
	"log"

	"github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/client"
	"github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/models"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/utils"
	"github.com/tshihad/tftags"
)

type routerds struct {
	nClient *client.RouterAPIService
}

func newRouterDS(nClient *client.RouterAPIService) *routerds {
	return &routerds{nClient: nClient}
}

func (n *routerds) Read(ctx context.Context, d *utils.Data, meta interface{}) error {
	log.Printf("[DEBUG] Get Router")
	name := d.GetString("name")

	// Pre check
	if err := d.Error(); err != nil {
		return err
	}
	resp, err := utils.Retry(ctx, meta, func(ctx context.Context) (interface{}, error) {
		return n.nClient.GetAllRouter(ctx, nil)
	})
	if err != nil {
		return err
	}

	routers := resp.(models.GetAllNetworkRouter)
	for i, n := range routers.NetworkRouters {
		if n.Name == name {
			log.Print("[DEBUG]", routers.NetworkRouters[i].ID)

			return tftags.Set(d, routers.NetworkRouters[i])
		}
	}

	return fmt.Errorf(errExactMatch, "Router")
}
