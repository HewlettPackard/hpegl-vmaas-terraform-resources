// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"
	"fmt"

	"github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/client"
	"github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/models"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/utils"
	"github.com/tshihad/tftags"
)

type routerBgpNeighbor struct {
	rClient *client.RouterAPIService
}

func newRouterBgpNeighbor(routerBgpNeighborClient *client.RouterAPIService) *routerBgpNeighbor {
	return &routerBgpNeighbor{
		rClient: routerBgpNeighborClient,
	}
}

func (r *routerBgpNeighbor) Read(ctx context.Context, d *utils.Data, meta interface{}) error {
	setMeta(meta, r.rClient.Client)
	var tfBgpNeighbor models.CreateRouterRequestBgpNeighborBody
	if err := tftags.Get(d, &tfBgpNeighbor); err != nil {
		return err
	}

	_, err := r.rClient.GetSpecificRouterBgpNeighbor(ctx, tfBgpNeighbor.RouterID, tfBgpNeighbor.ID)
	if err != nil {
		return err
	}

	return tftags.Set(d, tfBgpNeighbor)
}

func (r *routerBgpNeighbor) Create(ctx context.Context, d *utils.Data, meta interface{}) error {
	setMeta(meta, r.rClient.Client)
	var tfBgpNeighbor models.CreateRouterRequestBgpNeighborBody
	err := tftags.Get(d, &tfBgpNeighbor)
	if err != nil {
		return err
	}
	bgpNeighborRes, err := r.rClient.CreateRouterBgpNeighbor(ctx, tfBgpNeighbor.RouterID,
		models.CreateNetworkRouterBgpNeighborRequest{NetworkRouterBgpNeighbor: tfBgpNeighbor},
	)
	if err != nil {
		return err
	}

	if !bgpNeighborRes.Success {
		return fmt.Errorf(successErr, "creating BGPNEIGHBOR rule for the router")
	}
	tfBgpNeighbor.ID = bgpNeighborRes.ID

	return tftags.Set(d, tfBgpNeighbor)
}

func (r *routerBgpNeighbor) Update(ctx context.Context, d *utils.Data, meta interface{}) error {
	setMeta(meta, r.rClient.Client)
	var tfBgpNeighbor models.CreateRouterRequestBgpNeighborBody
	err := tftags.Get(d, &tfBgpNeighbor)
	if err != nil {
		return err
	}
	bgpNeighborRes, err := r.rClient.UpdateRouterBgpNeighbor(ctx, tfBgpNeighbor.RouterID, tfBgpNeighbor.ID,
		models.CreateNetworkRouterBgpNeighborRequest{NetworkRouterBgpNeighbor: tfBgpNeighbor},
	)
	if err != nil {
		return err
	}

	if !bgpNeighborRes.Success {
		return fmt.Errorf(successErr, "creating BGPNEIGHBOR rule for the router")
	}
	tfBgpNeighbor.ID = bgpNeighborRes.ID

	return tftags.Set(d, tfBgpNeighbor)
}

func (r *routerBgpNeighbor) Delete(ctx context.Context, d *utils.Data, meta interface{}) error {
	setMeta(meta, r.rClient.Client)
	var tfBgpNeighbor models.CreateRouterRequestBgpNeighborBody
	if err := tftags.Get(d, &tfBgpNeighbor); err != nil {
		return err
	}

	resp, err := r.rClient.DeleteRouterBgpNeighbor(ctx, tfBgpNeighbor.RouterID, tfBgpNeighbor.ID)
	if err != nil {
		return err
	}
	if !resp.Success {
		return fmt.Errorf("got success = 'false' while deleting BGPNEIGHBOR rule")
	}

	return nil
}
