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

type transportZone struct {
	tClient *client.RouterAPIService
}

func newTransportZone(tClient *client.RouterAPIService) *transportZone {
	return &transportZone{
		tClient: tClient,
	}
}

func (r *transportZone) Read(ctx context.Context, d *utils.Data, meta interface{}) error {
	setMeta(meta, r.tClient.Client)
	var tfScope models.NetworkScope
	if err := tftags.Get(d, &tfScope); err != nil {
		return err
	}

	// Get network server ID for nsx-t
	serverResp, err := r.tClient.GetNetworkServices(ctx, nil)
	if err != nil {
		return err
	}

	var serverID int
	for i, n := range serverResp.NetworkServices {
		if n.TypeName == nsxt {
			serverID = serverResp.NetworkServices[i].ID

			break
		}
	}

	if serverID == 0 {
		return fmt.Errorf(errExactMatch, "network server")
	}

	resp, err := r.tClient.GetTransportZones(ctx, serverID, tfScope.Name)
	if err != nil {
		return err
	}
	return tftags.Set(d, resp)
}
