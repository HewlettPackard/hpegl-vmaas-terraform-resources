// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"
	"fmt"

	"github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/client"
	"github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/models"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/utils"
)

type networkProxy struct {
	nClient *client.NetworksAPIService
}

func newNetworkProxy(nClient *client.NetworksAPIService) *networkProxy {
	return &networkProxy{nClient: nClient}
}

func (n *networkProxy) Read(ctx context.Context, d *utils.Data, meta interface{}) error {
	name := d.GetString("name")
	resp, err := utils.Retry(ctx, meta, func(ctx context.Context) (interface{}, error) {
		return n.nClient.GetNetworkPool(ctx, nil)
	})
	if err != nil {
		return err
	}

	poolResp := resp.(models.GetNetworkPoolsResp)
	for _, p := range poolResp.NetworkPools {
		if p.DisplayName == name {
			d.SetString("display_name", p.DisplayName)
			d.SetID(p.ID)

			return nil
		}
	}

	return fmt.Errorf(errExactMatch, "Network Pool")
}
