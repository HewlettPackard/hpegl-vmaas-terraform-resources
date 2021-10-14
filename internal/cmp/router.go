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
	routerID := d.GetID()
	routerResp, err := r.routerClient.GetSpecificRouter(ctx, routerID)
	if err != nil {
		return err
	}
	d.SetID(routerResp.NetworkRouter.ID)

	return nil
}

func (r *router) Create(ctx context.Context, d *utils.Data, meta interface{}) error {
	createReq := models.CreateRouterRequest{}
	if err := tftags.Get(d, &createReq.NetworkRouter); err != nil {
		return err
	}
	// align createReq and fill json related fields
	r.routerAlignRouterRequest(ctx, meta, &createReq)

	routerResp, err := r.routerClient.CreateRouter(ctx, createReq)
	if err != nil {
		return err
	}
	if !routerResp.Success {
		return fmt.Errorf("got success = 'false' while creating router")
	}
	d.SetID(routerResp.ID)

	return nil
}

func (r *router) Update(ctx context.Context, d *utils.Data, meta interface{}) error {
	return nil
}

func (r *router) Delete(ctx context.Context, d *utils.Data, meta interface{}) error {
	routerID := d.GetID()
	_, err := r.routerClient.DeleteRouter(ctx, routerID)
	if err != nil {
		return err
	}

	return nil
}

func (r *router) routerAlignRouterRequest(ctx context.Context, meta interface{}, routerReq *models.CreateRouterRequest) error {
	queryParam := make(map[string]string)
	// Check whether teir0 or tier1 and assign properties to proper child, so json can marshal properly
	if routerReq.NetworkRouter.TfTier0Config != nil {
		routerReq.NetworkRouter.Config.CreateRouterTier0Config.RouteRedistributionTier0 =
			routerReq.NetworkRouter.TfTier0Config.TfRRTier0
		routerReq.NetworkRouter.Config.CreateRouterTier0Config.RouteRedistributionTier1 =
			routerReq.NetworkRouter.TfTier0Config.TfRRTier1
		routerReq.NetworkRouter.Config.CreateRouterTier0Config.Bgp =
			routerReq.NetworkRouter.TfTier0Config.TfBGP

		routerReq.NetworkRouter.Config.HaMode = routerReq.NetworkRouter.TfTier0Config.TfHaMode
		routerReq.NetworkRouter.Config.FailOver = routerReq.NetworkRouter.TfTier0Config.TfFailOver
		routerReq.NetworkRouter.Config.EdgeCluster = routerReq.NetworkRouter.TfTier0Config.TfEdgeCluster
		routerReq.NetworkRouter.Config.EnableBgp = routerReq.NetworkRouter.TfTier0Config.TfEnableBgp
		queryParam[nameKey] = tier0GatewayType
	} else {
		routerReq.NetworkRouter.Config.CreateRouterTier0Config.RouteRedistributionTier1.RouteAdvertisement =
			routerReq.NetworkRouter.TfTier1Config.TfRouteAdvertisement
		routerReq.NetworkRouter.Config.EdgeCluster = routerReq.NetworkRouter.TfTier1Config.TfEdgeCluster
		routerReq.NetworkRouter.Config.Tier0Gateways = routerReq.NetworkRouter.TfTier1Config.TfTier0Gateways
		queryParam[nameKey] = tier1GatewayType
	}
	// Get router type
	rtRetry := utils.CustomRetry{}
	rtRetry.RetryParallel(ctx, meta, func(ctx context.Context) (interface{}, error) {
		return r.routerClient.GetRouterTypes(ctx, queryParam)
	})
	// Get network service ID
	nsRetry := utils.CustomRetry{}
	nsRetry.RetryParallel(ctx, meta, func(ctx context.Context) (interface{}, error) {
		return r.routerClient.GetNetworkServices(ctx, map[string]string{
			nameKey: "NSX-T",
		})
	})

	rtResp, err := rtRetry.Wait()
	if err != nil {
		return err
	}
	routerTypes := rtResp.(models.GetNetworlRouterTypes)
	if len(routerTypes.NetworkRouterTypes) != 1 {
		return fmt.Errorf(errExactMatch, "router-type")
	}
	routerReq.NetworkRouter.Type.ID = routerTypes.NetworkRouterTypes[0].ID

	nsResp, err := nsRetry.Wait()
	if err != nil {
		return err
	}
	networkService := nsResp.(models.GetNetworkServicesResp)
	if len(routerTypes.NetworkRouterTypes) != 1 {
		return fmt.Errorf(errExactMatch, "network-service")
	}
	routerReq.NetworkRouter.NetworkServer.ID = networkService.NetworkServices[0].ID
	routerReq.NetworkRouter.NetworkServerID = networkService.NetworkServices[0].ID

	routerReq.NetworkRouter.Site.ID = routerReq.NetworkRouter.GroupID

	return nil
}
