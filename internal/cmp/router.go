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
	rClient *client.RouterAPIService
}

func newRouter(routerClient *client.RouterAPIService) *router {
	return &router{
		rClient: routerClient,
	}
}

func (r *router) Read(ctx context.Context, d *utils.Data, meta interface{}) error {
	var tfRouter models.GetNetworkRouter
	if err := tftags.Get(d, &tfRouter); err != nil {
		return err
	}
	getRouter, err := r.rClient.GetSpecificRouter(ctx, tfRouter.ID)
	if err != nil {
		return err
	}

	return tftags.Set(d, getRouter.NetworkRouter)
}

func (r *router) Create(ctx context.Context, d *utils.Data, meta interface{}) error {
	createReq := models.CreateRouterRequest{}
	if err := tftags.Get(d, &createReq.NetworkRouter); err != nil {
		return err
	}
	// align createReq and fill json related fields
	if err := r.routerAlignRouterRequest(ctx, meta, &createReq); err != nil {
		return err
	}
	routerResp, err := r.rClient.CreateRouter(ctx, createReq)
	if err != nil {
		return err
	}
	if !routerResp.Success {
		return fmt.Errorf(successErr, "creating router")
	}
	createReq.NetworkRouter.ID = routerResp.ID

	// wait until created
	retry := &utils.CustomRetry{
		RetryDelay:   1,
		InitialDelay: 1,
		Cond: func(response interface{}, ResponseErr error) (bool, error) {
			return response.(models.GetSpecificRouterResp).NetworkRouter.Status == "ok", nil
		},
	}
	_, err = retry.Retry(ctx, meta, func(ctx context.Context) (interface{}, error) {
		return r.rClient.GetSpecificRouter(ctx, routerResp.ID)
	})
	if err != nil {
		return err
	}

	return tftags.Set(d, createReq.NetworkRouter)
}

func (r *router) Update(ctx context.Context, d *utils.Data, meta interface{}) error {
	createReq := models.CreateRouterRequest{}
	if err := tftags.Get(d, &createReq.NetworkRouter); err != nil {
		return err
	}
	// align createReq and fill json related fields
	if err := r.routerAlignRouterRequest(ctx, meta, &createReq); err != nil {
		return err
	}

	// HaMode cannot be updated, setting it to empty so that it is ignored in the API Payload.
	createReq.NetworkRouter.Config.HaMode = ""

	resp, err := utils.Retry(ctx, meta, func(ctx context.Context) (interface{}, error) {
		return r.rClient.UpdateRouter(ctx, createReq.NetworkRouter.ID, createReq)
	})
	if err != nil {
		return err
	}
	routerResp := resp.(models.SuccessOrErrorMessage)
	if !routerResp.Success {
		return fmt.Errorf("got success = 'false' while updating router")
	}

	d.SetID(createReq.NetworkRouter.ID)

	return nil
}

func (r *router) Delete(ctx context.Context, d *utils.Data, meta interface{}) error {
	routerID := d.GetID()
	_, err := r.rClient.DeleteRouter(ctx, routerID)
	if err != nil {
		return err
	}

	return nil
}

func (r *router) routerAlignRouterRequest(ctx context.Context, meta interface{}, routerReq *models.CreateRouterRequest) error {
	queryParam := make(map[string]string)
	// Check whether teir0 or tier1 and assign properties to proper child, so json can marshal properly
	tier0Config := routerReq.NetworkRouter.Config.CreateRouterTier0Config
	if routerReq.NetworkRouter.TfTier0Config != nil {
		tier0Config.RouteRedistributionTier0 = routerReq.NetworkRouter.TfTier0Config.TfRRTier0

		tfRRTier1 := routerReq.NetworkRouter.TfTier0Config.TfRRTier1
		tier0Config.RouteRedistributionTier1 = tfRRTier1
		tier0Config.RouteRedistributionTier1.RouteAdvertisement.TIER1DNSFORWARDERIP = tfRRTier1.TIER1DNSFORWARDERIP
		tier0Config.RouteRedistributionTier1.RouteAdvertisement.TIER1LBSNAT = tfRRTier1.TIER1LBSNAT
		tier0Config.RouteRedistributionTier1.RouteAdvertisement.TIER1NAT = tfRRTier1.TIER1NAT
		tier0Config.RouteRedistributionTier1.RouteAdvertisement.TIER1LBVIP = tfRRTier1.TIER1LBVIP
		tier0Config.RouteRedistributionTier1.RouteAdvertisement.TIER1IPSECLOCALENDPOINT = tfRRTier1.TIER1IPSECLOCALENDPOINT
		tier0Config.RouteRedistributionTier1.RouteAdvertisement.TIER1STATIC = tfRRTier1.TIER1STATIC
		tier0Config.RouteRedistributionTier1.RouteAdvertisement.Tier1Connected = tfRRTier1.Tier1Connected
		tier0Config.RouteRedistributionTier1.RouteAdvertisement.Tier1StaticRoutes = tfRRTier1.Tier1StaticRoutes

		tier0Config.Bgp = routerReq.NetworkRouter.TfTier0Config.TfBGP
		routerReq.NetworkRouter.Config.CreateRouterTier0Config = tier0Config
		routerReq.NetworkRouter.Config.HaMode = routerReq.NetworkRouter.TfTier0Config.TfHaMode
		routerReq.NetworkRouter.Config.FailOver = routerReq.NetworkRouter.TfTier0Config.TfFailOver
		routerReq.NetworkRouter.Config.EdgeCluster = routerReq.NetworkRouter.TfTier0Config.TfEdgeCluster
		routerReq.NetworkRouter.EnableBGP = routerReq.NetworkRouter.TfTier0Config.TfBGP.TfEnableBgp
		queryParam[nameKey] = tier0GatewayType
	} else {
		routerReq.NetworkRouter.Config.CreateRouterTier0Config.RouteRedistributionTier1.RouteAdvertisement = routerReq.NetworkRouter.TfTier1Config.TfRouteAdvertisement
		routerReq.NetworkRouter.Config.EdgeCluster = routerReq.NetworkRouter.TfTier1Config.TfEdgeCluster
		routerReq.NetworkRouter.Config.FailOver = routerReq.NetworkRouter.TfTier1Config.TfFailOver
		routerReq.NetworkRouter.Config.Tier0Gateways = routerReq.NetworkRouter.TfTier1Config.TfTier0Gateways
		queryParam[nameKey] = tier1GatewayType
	}
	// Get Router type
	rtRetry := utils.CustomRetry{}
	rtRetry.RetryParallel(ctx, meta, func(ctx context.Context) (interface{}, error) {
		return r.rClient.GetRouterTypes(ctx, queryParam)
	})
	// Get network service ID
	nsRetry := utils.CustomRetry{}
	nsRetry.RetryParallel(ctx, meta, func(ctx context.Context) (interface{}, error) {
		return r.rClient.GetNetworkServices(ctx, nil)
	})
	// Align Router Type
	rtResp, err := rtRetry.Wait()
	if err != nil {
		return err
	}
	routerTypes := rtResp.(models.GetNetworlRouterTypes)
	if len(routerTypes.NetworkRouterTypes) != 1 {
		return fmt.Errorf(errExactMatch, "router-type")
	}
	routerReq.NetworkRouter.Type.ID = routerTypes.NetworkRouterTypes[0].ID
	routerReq.NetworkRouter.TypeID = routerTypes.NetworkRouterTypes[0].ID
	// Align Network Server
	nsResp, err := nsRetry.Wait()
	if err != nil {
		return err
	}
	networkService := nsResp.(models.GetNetworkServicesResp)
	if len(routerTypes.NetworkRouterTypes) == 0 {
		return fmt.Errorf(errExactMatch, "network-service")
	}
	for i, n := range networkService.NetworkServices {
		if n.TypeName == nsxt {
			routerReq.NetworkRouter.NetworkServer.ID = networkService.NetworkServices[i].ID
			routerReq.NetworkRouter.NetworkServerID = networkService.NetworkServices[i].ID

			break
		}
	}
	routerReq.NetworkRouter.Site.ID = routerReq.NetworkRouter.GroupID

	return nil
}
