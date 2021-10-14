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

type networkProxy struct {
	nClient *client.NetworksAPIService
}

func newNetworkProxy(nClient *client.NetworksAPIService) *networkProxy {
	return &networkProxy{nClient: nClient}
}

func (n *networkProxy) Read(ctx context.Context, d *utils.Data, meta interface{}) error {
	tfProxy := models.GetNetworkProxy{}
	tftags.Get(d, &tfProxy)

	resp, err := utils.Retry(ctx, meta, func(ctx context.Context) (interface{}, error) {
		return n.nClient.GetNetworkProxy(ctx, map[string]string{
			nameKey: tfProxy.Name,
		})
	})
	if err != nil {
		return err
	}
	proxyResp := resp.(models.GetAllNetworkProxies)
	if len(proxyResp.GetNetworkProxies) != 1 {
		return fmt.Errorf(errExactMatch, "network proxy")
	}

	return tftags.Set(d, proxyResp.GetNetworkProxies[0])
}
