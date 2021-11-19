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

type routerRoute struct {
	rClient *client.RouterAPIService
}

func newRouterRoute(routerNatClient *client.RouterAPIService) *routerRoute {
	return &routerRoute{
		rClient: routerNatClient,
	}
}

func (r *routerRoute) Read(ctx context.Context, d *utils.Data, meta interface{}) error {
	setMeta(meta, r.rClient.Client)
	var tfRoute models.RouterRouteBody
	if err := tftags.Get(d, &tfRoute); err != nil {
		return err
	}
	// Get the router, if the router not exists, return warning
	if check, err := checkRouterDeprecated(
		ctx, r.rClient, d, tfRoute.RouterID, &tfRoute.IsDeprecated, &tfRoute,
	); err != nil || check {
		return err
	}

	resp, err := r.rClient.GetSpecificRouterRoute(ctx, tfRoute.RouterID, tfRoute.ID)
	if err != nil {
		return err
	}

	return tftags.Set(d, resp.NetworkRoute)
}

func (r *routerRoute) Create(ctx context.Context, d *utils.Data, meta interface{}) error {
	setMeta(meta, r.rClient.Client)
	var tfRoute models.RouterRouteBody
	err := tftags.Get(d, &tfRoute)
	if err != nil {
		return err
	}
	natRes, err := r.rClient.CreateRouterRoute(ctx, tfRoute.RouterID,
		models.CreateRouterRoute{
			NetworkRoute: tfRoute,
		},
	)
	if err != nil {
		return err
	}

	if !natRes.Success {
		return fmt.Errorf(successErr, "creating route for the router")
	}
	tfRoute.ID = natRes.ID

	return tftags.Set(d, tfRoute)
}

func (r *routerRoute) Update(ctx context.Context, d *utils.Data, meta interface{}) error {
	return nil
}

func (r *routerRoute) Delete(ctx context.Context, d *utils.Data, meta interface{}) error {
	setMeta(meta, r.rClient.Client)
	var tfRoute models.RouterRouteBody
	if err := tftags.Get(d, &tfRoute); err != nil {
		return err
	}

	// if parent router got deleted, NAT is already deleted
	if tfRoute.IsDeprecated {
		log.Printf("[WARNING] Router route already deleted since router is deleted")

		return nil
	}

	resp, err := r.rClient.DeleteRouterRoute(ctx, tfRoute.RouterID, tfRoute.ID)
	if err != nil {
		return err
	}
	if !resp.Success {
		return fmt.Errorf("got success = 'false' while deleting NAT rule")
	}

	return nil
}
