// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hpe-hcss/vmaas-cmp-go-sdk/pkg/client"
	"github.com/hpe-hcss/vmaas-cmp-go-sdk/pkg/models"
	"github.com/hpe-hcss/vmaas-terraform-resources/internal/logger"
	"github.com/hpe-hcss/vmaas-terraform-resources/internal/utils"
)

type network struct {
	nClient           *client.NetworksApiService
	serviceInstanceID string
}

func newNetwork(nClient *client.NetworksApiService, serviceInstanceID string) *network {
	return &network{nClient: nClient, serviceInstanceID: serviceInstanceID}
}

func (n *network) Read(ctx context.Context, d *utils.Data) error {
	logger.Debug("Get Network")

	name := d.GetString("name")
	// Pre check
	if err := d.Error(); err != nil {
		return err
	}
	resp, err := utils.Retry(func() (interface{}, error) {
		return n.nClient.GetAllNetworks(ctx, n.serviceInstanceID, map[string]string{
			nameKey: name,
		})
	})
	if err != nil {
		return err
	}

	networks := resp.(models.ListNetworksBody)
	if len(networks.Networks) != 1 {
		return fmt.Errorf(errExactMatch, "Network")
	}
	d.SetID(strconv.Itoa(networks.Networks[0].Id))

	// post check
	return d.Error()
}
