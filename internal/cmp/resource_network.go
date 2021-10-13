// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"

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
	var tfNetwork models.GetSpecificNetwork
	if err := tftags.Get(d, &tfNetwork); err != nil {
		return err
	}

	// Get network details with ID
	response, err := utils.Retry(ctx, meta, func(ctx context.Context) (interface{}, error) {
		return r.rClient.GetSpecificNetwork(ctx, tfNetwork.ID)
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

	createReq.Zone.ID = createReq.CloudID
	createReq.Site.ID = createReq.GroupID
	createReq.Type.ID = createReq.TypeID

	// Create network
	resp, err := utils.Retry(ctx, meta, func(ctx context.Context) (interface{}, error) {
		return r.rClient.CreateNetwork(ctx, models.CreateNetworkRequest{
			Network: createReq,
		})
	})
	if err != nil {
		return err
	}

	return tftags.Set(d, resp.(models.CreateNetworkResponse).Network)
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
