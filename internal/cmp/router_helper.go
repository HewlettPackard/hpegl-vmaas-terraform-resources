// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"
	"log"

	"github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/client"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/utils"
	"github.com/tshihad/tftags"
)

func checkRouterDeprecated(
	ctx context.Context,
	rClient *client.RouterAPIService,
	d *utils.Data,
	routerID int,
	isDeprecated *bool,
	tfModel interface{},
) (bool, error) {
	*isDeprecated = false
	// Get the router, if the router not exists, return warning
	router, err := rClient.GetSpecificRouter(ctx, routerID)
	if err != nil {
		return false, err
	}
	// if router not found set is_deprecated flag=true
	if router.ID == 0 {
		log.Printf("[ERROR] Router with %d id is not found on NAT plan", routerID)
		*isDeprecated = true

		return true, tftags.Set(d, tfModel)
	}

	return false, nil
}
