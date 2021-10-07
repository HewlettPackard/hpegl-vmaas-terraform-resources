// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"
	"fmt"
	"log"

	"github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/client"
	"github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/models"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/utils"
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

	isMatch := false
	routers := resp.(models.GetAllNetworkRouter)
	for i, n := range routers.NetworkRouters {
		if n.Name == name {
			isMatch = true
			d.SetID(routers.NetworkRouters[i].ID)
			log.Print("[DEBUG]", routers.NetworkRouters[i].ID)
			break
		}
	}
	if !isMatch {
		return fmt.Errorf(errExactMatch, "Router")
	}

	// post check
	return d.Error()
}
