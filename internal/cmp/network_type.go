// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"
	"fmt"

	"github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/client"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/utils"
)

type networkType struct {
	nClient *client.NetworksAPIService
}

func newNetworkType(nClient *client.NetworksAPIService) *networkType {
	return &networkType{nClient: nClient}
}

func (n *networkType) Read(ctx context.Context, d *utils.Data, meta interface{}) error {
	name := d.GetString("name")
	networkResp, err := n.nClient.GetNetworkType(ctx, map[string]string{
		nameKey: name,
	})
	if err != nil {
		return err
	}
	if len(networkResp.NetworkTypes) != 1 {
		return fmt.Errorf(errExactMatch, "network-type")
	}

	d.SetID(networkResp.NetworkTypes[0].ID)

	return nil
}
