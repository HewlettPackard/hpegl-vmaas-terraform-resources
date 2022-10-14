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
	// align createReq and fill json related fields
	if err := r.networkRequest(&createReq); err != nil {
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

	// Create network
	createResp, err := r.nClient.CreateNetwork(ctx, models.CreateNetworkRequest{
		Network: createReq,
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

	// align createReq and fill json related fields
	if err := r.networkRequest(&networkReq); err != nil {
		return err
	}

	updateResp, err := r.nClient.UpdateNetwork(ctx, networkReq.ID, models.CreateNetworkRequest{
		Network: networkReq,
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

func (r *resNetwork) networkRequest(createReq *models.CreateNetwork) error {
	if createReq.TfDhcpNetwork != nil {
		createReq.Name = createReq.TfDhcpNetwork.Name
		createReq.Description = createReq.TfDhcpNetwork.Description
		createReq.DisplayName = createReq.TfDhcpNetwork.DisplayName
		createReq.Cidr = createReq.TfDhcpNetwork.Cidr
		createReq.DNSPrimary = createReq.TfDhcpNetwork.DNSPrimary
		createReq.DNSSecondary = createReq.TfDhcpNetwork.DNSSecondary
		createReq.Gateway = createReq.TfDhcpNetwork.Gateway
		createReq.NoProxy = createReq.TfDhcpNetwork.NoProxy
		createReq.ScopeID = createReq.TfDhcpNetwork.ScopeID
		createReq.SearchDomains = createReq.TfDhcpNetwork.SearchDomains
		createReq.Active = createReq.TfDhcpNetwork.Active
		createReq.AllowStaticOverride = createReq.TfDhcpNetwork.AllowStaticOverride
		createReq.DhcpServer = createReq.TfDhcpNetwork.DhcpServer
		createReq.AppURLProxyBypass = createReq.TfDhcpNetwork.AppURLProxyBypass
		createReq.ScanNetwork = createReq.TfDhcpNetwork.ScanNetwork
		if createReq.TfDhcpNetwork.ResourcePermissions != nil {
			createReq.ResourcePermissions.All = createReq.TfDhcpNetwork.ResourcePermissions.All
		}
		if createReq.TfDhcpNetwork.Site != nil {
			createReq.Site = &models.IDStringModel{createReq.TfDhcpNetwork.Site.ID}
		}
		if createReq.TfDhcpConfig != nil {
			createReq.DhcpConfig.ConnectedGateway = createReq.TfDhcpConfig.ConnectedGateway
			createReq.DhcpConfig.VlanIDs = createReq.TfDhcpConfig.VlanIDs
			createReq.DhcpConfig.SubnetDhcpLeaseTime = createReq.TfDhcpConfig.SubnetDhcpLeaseTime
			createReq.DhcpConfig.SubnetDhcpServerAddress = createReq.TfDhcpConfig.SubnetDhcpServerAddress
			createReq.DhcpConfig.SubnetIPManagementType = createReq.TfDhcpConfig.SubnetIPManagementType
			createReq.DhcpConfig.SubnetIPServerID = createReq.TfDhcpConfig.SubnetIPServerID
		}
	}

	if createReq.TfStaticNetwork != nil {
		createReq.Name = createReq.TfStaticNetwork.Name
		createReq.Description = createReq.TfStaticNetwork.Description
		createReq.DisplayName = createReq.TfStaticNetwork.DisplayName
		createReq.Cidr = createReq.TfStaticNetwork.Cidr
		createReq.DNSPrimary = createReq.TfStaticNetwork.DNSPrimary
		createReq.DNSSecondary = createReq.TfStaticNetwork.DNSSecondary
		createReq.Gateway = createReq.TfStaticNetwork.Gateway
		createReq.NoProxy = createReq.TfStaticNetwork.NoProxy
		createReq.ScopeID = createReq.TfStaticNetwork.ScopeID
		createReq.SearchDomains = createReq.TfStaticNetwork.SearchDomains
		createReq.Active = createReq.TfStaticNetwork.Active
		createReq.AllowStaticOverride = createReq.TfStaticNetwork.AllowStaticOverride
		createReq.AppURLProxyBypass = createReq.TfStaticNetwork.AppURLProxyBypass
		createReq.ScanNetwork = createReq.TfStaticNetwork.ScanNetwork
		createReq.PoolID = createReq.TfStaticNetwork.PoolID
		if createReq.TfStaticNetwork.ResourcePermissions != nil {
			createReq.ResourcePermissions.All = createReq.TfStaticNetwork.ResourcePermissions.All
		}
		if createReq.TfStaticConfig != nil {
			createReq.StaticConfig.ConnectedGateway = createReq.TfStaticConfig.ConnectedGateway
			createReq.StaticConfig.VlanIDs = createReq.TfStaticConfig.VlanIDs
		}
	}
	return nil
}
