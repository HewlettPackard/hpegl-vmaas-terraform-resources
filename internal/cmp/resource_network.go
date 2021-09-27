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
	rClient *client.NetworksAPIService
}

func newResNetwork(client *client.NetworksAPIService) *resNetwork {
	return &resNetwork{
		rClient: client,
	}
}

func (r *resNetwork) Read(ctx context.Context, d *utils.Data, meta interface{}) error {
	// get network details
	var tfNetwork models.GetSpecificNetworkBody
	if err := tftags.Get(d, &tfNetwork); err != nil {
		return err
	}
	tfNetwork.Network.ID = d.GetID()

	// Get network details with ID
	response, err := utils.Retry(ctx, meta, func(ctx context.Context) (interface{}, error) {
		return r.rClient.GetSpecificNetwork(ctx, tfNetwork.Network.ID)
	})
	if err != nil {
		return err
	}
	getNetwork := response.(models.GetSpecificNetworkBody)

	if err := tftags.Set(d, getNetwork.Network); err != nil {
		return err
	}

	d.SetID(tfNetwork.Network.ID)

	return nil
}

func (r *resNetwork) Create(ctx context.Context, d *utils.Data, meta interface{}) error {
	var createReq models.CreateNetwork
	if err := tftags.Get(d, &createReq); err != nil {
		return err
	}

	createReq.Zone.ID = createReq.CloudID
	createReq.Site.ID = createReq.GroupID
	// Get network type ID for NSX T Segment
	typeResp, err := utils.Retry(ctx, meta, func(ctx context.Context) (interface{}, error) {
		return r.rClient.GetNetworkType(ctx, map[string]string{
			nameKey: nsxtSegment,
		})
	})
	if err != nil {
		return fmt.Errorf("failed to retrieve %s, got error %v", nsxtSegment, err)
	}

	networkTypes := typeResp.(models.GetNetworkTypesResponse)
	if len(networkTypes.NetworkTypes) != 1 {
		return fmt.Errorf("couldn't find NSX-T integration type")
	}
	createReq.Type.ID = networkTypes.NetworkTypes[0].ID

	// Create network
	resp, err := utils.Retry(ctx, meta, func(ctx context.Context) (interface{}, error) {
		return r.rClient.CreateNetwork(ctx, models.CreateNetworkRequest{
			Network: createReq,
		})
	})

	if err != nil {
		return err
	}

	d.SetID(resp.(models.CreateNetworkResponse).Network.ID)

	return nil
}

func (r *resNetwork) Update(ctx context.Context, d *utils.Data, meta interface{}) error {
	return nil
}

func (r *resNetwork) Delete(ctx context.Context, d *utils.Data, meta interface{}) error {
	networkID := d.GetID()
	utils.Retry(ctx, meta, func(ctx context.Context) (interface{}, error) {
		return r.rClient.DeleteNetwork(ctx, networkID)
	})

	return nil
}
