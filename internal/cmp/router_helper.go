// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"
	"log"
	"net/http"

	"github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/client"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/utils"
	pkgutils "github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/utils"
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
	// Get the router with router ID and check the error
	_, err := rClient.GetSpecificRouter(ctx, routerID)
	if err != nil {
		statusCode := pkgutils.GetStatusCode(err)
		// if router not found set is_deprecated flag=true
		if statusCode == http.StatusNotFound {
			log.Printf("[ERROR] Router with %d id is not found", routerID)
			*isDeprecated = true

			return true, tftags.Set(d, tfModel)
		}
		// return error for all other status code
		return false, err
	}

	// Router exists
	return false, nil
}
