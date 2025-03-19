// (C) Copyright 2021-2025 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"
	"fmt"
	"strings"
	"time"

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
	nsxType, err := GetNsxTypeFromCMP(ctx, r.rClient.Client)
	if err != nil {
		return err
	}
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
			nameKey: fmt.Sprintf("%s %s", nsxType, nsxSegment),
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
		if n.TypeName == nsxType {
			createReq.NetworkServer.ID = networkService.NetworkServices[i].ID

			break
		}
	}

	// Create network
	createResp, err := r.nClient.CreateNetwork(ctx, models.CreateNetworkRequest{
		Network:             createReq,
		ResourcePermissions: createReq.ResourcePermissions,
	})
	if err != nil {
		return err
	}
	cmpVersion := r.rClient.Client.GetSCMVersion()
	// from 8.0.2  onwards the network obj is fixed
	if v, _ := ParseVersion("8.0.3"); v > cmpVersion {

		// Refresh NSX integration
		serverRefreshResp, err := r.rClient.RefreshNetworkServices(ctx, createReq.NetworkServer.ID, nil)
		if !serverRefreshResp.Success {
			return fmt.Errorf("failed refresh NSX integration post NSX object creation")
		}
		if err != nil {
			return err
		}
		errCount := 0

		cRetry := utils.CustomRetry{
			Timeout:      time.Minute * 15,
			RetryDelay:   time.Second * 10,
			InitialDelay: time.Second * 20,
			Cond: func(response interface{}, err error) (bool, error) {
				if err != nil {
					errCount++
					// return false as condition if same error returns 3 times.
					if errCount == 3 {
						return false, err
					}

					return false, nil
				}

				networkResponse, ok := response.(models.GetSpecificNetworkBody)
				if !ok {
					errCount++
					if errCount == 3 {
						return false, fmt.Errorf("%s", "error while getting Network")
					}

					return false, nil
				}
				errCount = 0

				if strings.Contains(networkResponse.Network.ExternalID, utils.PortGroupPrefix) {
					return true, nil
				}

				return false, nil
			},
		}

		_, err = cRetry.Retry(ctx, meta, func(ctx context.Context) (interface{}, error) {
			return r.nClient.GetSpecificNetwork(ctx, createResp.Network.ID)
		})
		if err != nil {
			return err
		}
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
	// wait until deleted
	retry := &utils.CustomRetry{
		InitialDelay: time.Second * 10,
		RetryDelay:   time.Second * 10,
		Timeout:      time.Minute * 10,
	}
	resp, err := retry.Retry(ctx, meta, func(ctx context.Context) (interface{}, error) {
		return r.nClient.DeleteNetwork(ctx, networkID)
	})
	response := resp.(models.SuccessOrErrorMessage)
	if !response.Success {
		return fmt.Errorf("failed to delete the network due to following error: %s", response.Msg)
	}
	if err != nil {
		return err
	}

	return err
}

func (r *resNetwork) networkRequest(createReq *models.CreateNetwork) error {
	createReq.Site.ID = createReq.GroupID
	if createReq.NetworkDomainID != 0 {
		createReq.NetworkDomain = &models.IDModel{ID: createReq.NetworkDomainID}
	}
	if createReq.ProxyID != 0 {
		createReq.NetworkProxy = &models.IDModel{ID: createReq.ProxyID}
	}

	if createReq.TfDhcpNetwork != nil {
		createReq.Config.ConnectedGateway = createReq.ConnectedGateway
		createReq.Config.VlanIDs = createReq.VlanIDs
		createReq.Config.SubnetDhcpLeaseTime = createReq.TfDhcpNetwork.SubnetDhcpLeaseTime
		createReq.Config.DhcpRange = createReq.TfDhcpNetwork.DhcpRange
		createReq.Config.SubnetDhcpServerAddress = createReq.TfDhcpNetwork.SubnetDhcpServerAddress
		createReq.Config.SubnetIPManagementType = createReq.TfDhcpNetwork.SubnetIPManagementType
		createReq.Config.SubnetIPServerID = createReq.TfDhcpNetwork.SubnetIPServerID
	}
	if createReq.TfStaticNetwork != nil {
		createReq.Config.ConnectedGateway = createReq.ConnectedGateway
		createReq.Config.VlanIDs = createReq.VlanIDs
		createReq.PoolID = createReq.TfStaticNetwork.PoolID
	}

	return nil
}
