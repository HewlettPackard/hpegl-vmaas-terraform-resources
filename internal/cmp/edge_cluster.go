// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"
	"fmt"
	"log"

	"github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/client"
	"github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/models"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/utils"
	"github.com/tshihad/tftags"
)

type edgeCluster struct {
	tClient *client.RouterAPIService
}

func newEdgeCluster(tClient *client.RouterAPIService) *edgeCluster {
	return &edgeCluster{
		tClient: tClient,
	}
}

func (r *edgeCluster) Read(ctx context.Context, d *utils.Data, meta interface{}) error {
	setMeta(meta, r.tClient.Client)
	log.Printf("[INFO] Get Edge Cluster")
	var tfEdgeCluster models.NetworkEdgeClusters
	if err := tftags.Get(d, &tfEdgeCluster); err != nil {
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
	resp, err := r.tClient.GetEdgeCluster(ctx, serverID, tfEdgeCluster.Name)
	if err != nil {
		return err
	}

	return tftags.Set(d, resp)
}
