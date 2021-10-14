// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"
	"fmt"

	"github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/client"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/utils"
)

type networkPool struct {
	nClient *client.NetworksAPIService
}

func newNetworkPool(nClient *client.NetworksAPIService) *networkPool {
	return &networkPool{nClient: nClient}
}

func (n *networkPool) Read(ctx context.Context, d *utils.Data, meta interface{}) error {
	name := d.GetString("name")
	poolResp, err := n.nClient.GetNetworkPool(ctx, nil)
	if err != nil {
		return err
	}

	for _, p := range poolResp.NetworkPools {
		if p.DisplayName == name {
			d.SetString("display_name", p.DisplayName)
			d.SetID(p.ID)

			return nil
		}
	}

	return fmt.Errorf(errExactMatch, "Network Pool")
}
