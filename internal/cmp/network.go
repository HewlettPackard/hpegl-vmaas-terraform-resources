// (C) Copyright 2021-2024 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"
	"fmt"
	"log"

	"github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/client"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/utils"
)

type network struct {
	nClient *client.NetworksAPIService
}

func newNetwork(nClient *client.NetworksAPIService) *network {
	return &network{nClient: nClient}
}

func (n *network) Read(ctx context.Context, d *utils.Data, meta interface{}) error {
	setMeta(meta, n.nClient.Client)
	log.Printf("[DEBUG] Get Network")

	name := d.GetString("name")
	// Pre check
	if err := d.Error(); err != nil {
		return err
	}
	param := map[string]string{maxKey: "-1"}
	networks, err := n.nClient.GetAllNetworks(ctx, param)
	if err != nil {
		return err
	}

	for i, n := range networks.Networks {
		if n.Name == name {
			d.SetID(networks.Networks[i].ID)

			return nil
		}
	}

	return fmt.Errorf(errExactMatch, "Network")
}
