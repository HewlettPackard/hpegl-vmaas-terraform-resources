// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"

	"github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/client"
	"github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/models"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/utils"
	"github.com/tshihad/tftags"
)

type routerNat struct {
	routerNatClient *client.RouterAPIService
}

func newRouterNat(routerNatClient *client.RouterAPIService) *routerNat {
	return &routerNat{
		routerNatClient: routerNatClient,
	}
}

func (r *routerNat) Read(ctx context.Context, d *utils.Data, meta interface{}) error {
	var tfNat models.GetSpecificRouterNat
	if err := tftags.Get(d, &tfNat); err != nil {
		return err
	}
	r.routerNatClient.GetSpecificRouterNat(ctx, tfNat.ID, tfNat.ID)
	return nil
}

func (r *routerNat) Create(ctx context.Context, d *utils.Data, meta interface{}) error {
	return nil

}

func (r *routerNat) Update(ctx context.Context, d *utils.Data, meta interface{}) error {

	return nil
}

func (r *routerNat) Delete(ctx context.Context, d *utils.Data, meta interface{}) error {

	return nil
}
