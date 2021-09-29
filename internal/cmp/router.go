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

type router struct {
	routerClient *client.RouterAPIService
}

func newRouter(routerClient *client.RouterAPIService) *router {
	return &router{
		routerClient: routerClient,
	}
}

func (r *router) Read(ctx context.Context, d *utils.Data, meta interface{}) error {

	return nil
}

func (r *router) Create(ctx context.Context, d *utils.Data, meta interface{}) error {
	createReq := models.CreateRouterRequest{}
	if err := tftags.Get(d, &createReq.NetworkRouter); err != nil {
		return err
	}
	routerResp, err := r.routerClient.CreateRouter(ctx, createReq)
	if err != nil {
		return err
	}
	fmt.Printf("\n%#v\n", routerResp)
	return nil
}

func (r *router) Update(ctx context.Context, d *utils.Data, meta interface{}) error {
	return nil
}

func (r *router) Delete(ctx context.Context, d *utils.Data, meta interface{}) error {
	return nil
}
