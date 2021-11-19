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

type routerNat struct {
	rClient *client.RouterAPIService
}

func newRouterNat(routerNatClient *client.RouterAPIService) *routerNat {
	return &routerNat{
		rClient: routerNatClient,
	}
}

func (r *routerNat) Read(ctx context.Context, d *utils.Data, meta interface{}) error {
	setMeta(meta, r.rClient.Client)
	var tfNat models.CreateRouterNat
	if err := tftags.Get(d, &tfNat); err != nil {
		return err
	}
	// Get the router, if the router not exists, return warning
	if check, err := checkRouterDeprecated(
		ctx, r.rClient, d, tfNat.RouterID, &tfNat.IsDeprecated, &tfNat,
	); err != nil || check {
		return err
	}

	_, err := r.rClient.GetSpecificRouterNat(ctx, tfNat.RouterID, tfNat.ID)
	if err != nil {
		return err
	}
	tfNat.IsDeprecated = false

	return tftags.Set(d, tfNat)
}

func (r *routerNat) Create(ctx context.Context, d *utils.Data, meta interface{}) error {
	setMeta(meta, r.rClient.Client)
	var tfNat models.CreateRouterNat
	err := tftags.Get(d, &tfNat)
	if err != nil {
		return err
	}
	natRes, err := r.rClient.CreateRouterNat(ctx, tfNat.RouterID,
		models.CreateRouterNatRequest{CreateRouterNat: tfNat},
	)
	if err != nil {
		return err
	}

	if !natRes.Success {
		return fmt.Errorf(successErr, "creating NAT rule for the router")
	}
	tfNat.ID = natRes.ID

	return tftags.Set(d, tfNat)
}

func (r *routerNat) Update(ctx context.Context, d *utils.Data, meta interface{}) error {
	setMeta(meta, r.rClient.Client)
	var tfNat models.CreateRouterNat
	err := tftags.Get(d, &tfNat)
	if err != nil {
		return err
	}
	natRes, err := r.rClient.UpdateRouterNat(ctx, tfNat.RouterID, tfNat.ID,
		models.CreateRouterNatRequest{CreateRouterNat: tfNat},
	)
	if err != nil {
		return err
	}

	if !natRes.Success {
		return fmt.Errorf(successErr, "creating NAT rule for the router")
	}
	tfNat.ID = natRes.ID

	return tftags.Set(d, tfNat)
}

func (r *routerNat) Delete(ctx context.Context, d *utils.Data, meta interface{}) error {
	setMeta(meta, r.rClient.Client)
	var tfNat models.CreateRouterNat
	if err := tftags.Get(d, &tfNat); err != nil {
		return err
	}

	// if parent router got deleted, NAT is already deleted
	if tfNat.IsDeprecated {
		log.Printf("[WARNING] NAT already deleted since router is deleted")

		return nil
	}

	resp, err := r.rClient.DeleteRouterNat(ctx, tfNat.RouterID, tfNat.ID)
	if err != nil {
		return err
	}
	if !resp.Success {
		return fmt.Errorf("got success = 'false' while deleting NAT rule")
	}

	return nil
}
