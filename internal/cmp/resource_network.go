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

type resNetwork struct {
	nClient *client.NetworksAPIService
	rClient *client.RouterAPIService
}

func newResNetwork(nclient *client.NetworksAPIService, rclient *client.RouterAPIService) *resNetwork {
	return &resNetwork{
		nClient: nclient,
		rClient: rclient,
	}
}

func (r *resNetwork) Read(ctx context.Context, d *utils.Data, meta interface{}) error {
	// get network details
	var tfNetwork models.GetSpecificNetwork
	if err := tftags.Get(d, &tfNetwork); err != nil {
		return err
	}

	// Get network details with ID
	response, err := utils.Retry(ctx, meta, func(ctx context.Context) (interface{}, error) {
		return r.nClient.GetSpecificNetwork(ctx, tfNetwork.ID)
	})
	if err != nil {
		return err
	}
	getNetwork := response.(models.GetSpecificNetworkBody)

	return tftags.Set(d, getNetwork.Network)
}

func (r *resNetwork) Create(ctx context.Context, d *utils.Data, meta interface{}) error {
	var createReq models.CreateNetwork
	if err := tftags.Get(d, &createReq); err != nil {
		return err
	}

	// Get network type id for NSX-T
	typeRetry := utils.CustomRetry{}
	typeRetry.RetryParallel(ctx, meta, func(ctx context.Context) (interface{}, error) {
		return r.nClient.GetNetworkType(ctx, map[string]string{
			nameKey: nsxtSegment,
		})
	})
	// Get network server ID for nsx-t
	serverRetry := utils.CustomRetry{}
	serverRetry.RetryParallel(ctx, meta, func(ctx context.Context) (interface{}, error) {
		return r.rClient.GetNetworkServices(ctx, map[string]string{
			nameKey: "NSX-T",
		})
	})
	typeResp, err := typeRetry.Wait()
	if err != nil {
		return err
	}
	networkTypeResp := typeResp.(models.GetNetworkTypesResponse)
	if len(networkTypeResp.NetworkTypes) != 1 {
		return fmt.Errorf(errExactMatch, "network type")
	}

	serverResp, err := serverRetry.Wait()
	if err != nil {
		return err
	}
	networkServer := serverResp.(models.GetNetworkServicesResp)
	if len(networkServer.NetworkServices) != 1 {
		return fmt.Errorf(errExactMatch, "network server")
	}

	// Align request
	createReq.NetworkServer.ID = networkServer.NetworkServices[0].ID
	createReq.Type.ID = networkTypeResp.NetworkTypes[0].ID
	alignNetworkReq(&createReq)

	// Create network
	createResp, err := utils.Retry(ctx, meta, func(ctx context.Context) (interface{}, error) {
		return r.nClient.CreateNetwork(ctx, models.CreateNetworkRequest{
			Network:             createReq,
			ResourcePermissions: createReq.ResourcePermissions,
		})
	})
	if err != nil {
		return err
	}

	return tftags.Set(d, createResp.(models.CreateNetworkResponse).Network)
}

func (r *resNetwork) Update(ctx context.Context, d *utils.Data, meta interface{}) error {
	var networkReq models.CreateNetwork
	if err := tftags.Get(d, &networkReq); err != nil {
		return err
	}

	alignNetworkReq(&networkReq)
	resp, err := utils.Retry(ctx, meta, func(ctx context.Context) (interface{}, error) {
		return r.nClient.UpdateNetwork(ctx, networkReq.ID, models.CreateNetworkRequest{
			Network:             networkReq,
			ResourcePermissions: networkReq.ResourcePermissions,
		})
	})
	if err != nil {
		return err
	}

	updateResp := resp.(models.SuccessOrErrorMessage)
	if !updateResp.Success {
		return fmt.Errorf("failed to update network, got success = 'false' on update, error: %v", updateResp.Error)
	}
	d.SetID(networkReq.ID)

	return nil
}

func (r *resNetwork) Delete(ctx context.Context, d *utils.Data, meta interface{}) error {
	networkID := d.GetID()
	utils.Retry(ctx, meta, func(ctx context.Context) (interface{}, error) {
		return r.nClient.DeleteNetwork(ctx, networkID)
	})

	return nil
}

func alignNetworkReq(request *models.CreateNetwork) {
	request.Site.ID = request.GroupID
	if request.PoolID != 0 {
		request.Pool = &models.IDModel{ID: request.PoolID}
	}
	if request.NetworkDomainID != 0 {
		request.NetworkDomain = &models.IDModel{ID: request.NetworkDomainID}

	}
	if request.ProxyID != 0 {
		request.NetworkProxy = &models.IDModel{ID: request.ProxyID}
	}
}
