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
	setMeta(meta, r.rClient.Client)
	// get network details
	var tfNetwork models.GetSpecificNetwork
	if err := tftags.Get(d, &tfNetwork); err != nil {
		return err
	}

	// Get network details with ID
	getNetwork, err := r.nClient.GetSpecificNetwork(ctx, tfNetwork.ID)
	if err != nil {
		return err
	}

	return tftags.Set(d, getNetwork.Network)
}

func (r *resNetwork) Create(ctx context.Context, d *utils.Data, meta interface{}) error {
	setMeta(meta, r.rClient.Client)
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
		return r.rClient.GetNetworkServices(ctx, nil)
	})
	typeResp, err := typeRetry.Wait()
	if err != nil {
		return err
	}
	networkTypeResp := typeResp.(models.GetNetworkTypesResponse)
	if len(networkTypeResp.NetworkTypes) != 1 {
		return fmt.Errorf(errExactMatch, "network type")
	}
	// Align request for Network Type
	createReq.Type.ID = networkTypeResp.NetworkTypes[0].ID

	serverResp, err := serverRetry.Wait()
	if err != nil {
		return err
	}

	// Align request for Network Server
	networkService := serverResp.(models.GetNetworkServicesResp)
	if len(networkService.NetworkServices) == 0 {
		return fmt.Errorf(errExactMatch, "network server")
	}
	for i, n := range networkService.NetworkServices {
		if n.TypeName == nsxt {
			createReq.NetworkServer.ID = networkService.NetworkServices[i].ID

			break
		}
	}

	alignNetworkReq(&createReq)

	// Create network
	createResp, err := r.nClient.CreateNetwork(ctx, models.CreateNetworkRequest{
		Network:             createReq,
		ResourcePermissions: createReq.ResourcePermissions,
	})
	if err != nil {
		return err
	}

	return tftags.Set(d, createResp.Network)
}

func (r *resNetwork) Update(ctx context.Context, d *utils.Data, meta interface{}) error {
	setMeta(meta, r.rClient.Client)
	var networkReq models.CreateNetwork
	if err := tftags.Get(d, &networkReq); err != nil {
		return err
	}

	alignNetworkReq(&networkReq)
	updateResp, err := r.nClient.UpdateNetwork(ctx, networkReq.ID, models.CreateNetworkRequest{
		Network:             networkReq,
		ResourcePermissions: networkReq.ResourcePermissions,
	})
	if err != nil {
		return err
	}

	if !updateResp.Success {
		return fmt.Errorf("failed to update network, got success = 'false' on update, error: %v", updateResp.Error)
	}
	d.SetID(networkReq.ID)

	return nil
}

func (r *resNetwork) Delete(ctx context.Context, d *utils.Data, meta interface{}) error {
	setMeta(meta, r.rClient.Client)
	networkID := d.GetID()
	resp, err := r.nClient.DeleteNetwork(ctx, networkID)
	if !resp.Success {
		return fmt.Errorf("got success as false on delete network, error: %w", err)
	}

	return err
}

func alignNetworkReq(request *models.CreateNetwork) {
	request.Site.ID = request.GroupID
	if request.NetworkDomainID != 0 {
		request.NetworkDomain = &models.IDModel{ID: request.NetworkDomainID}
	}
	if request.ProxyID != 0 {
		request.NetworkProxy = &models.IDModel{ID: request.ProxyID}
	}
}
