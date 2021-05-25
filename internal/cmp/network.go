// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"
	"errors"
	"strconv"

	"github.com/hpe-hcss/vmaas-cmp-go-sdk/pkg/client"
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
	networks, err := n.nClient.GetAllNetworks(ctx, n.serviceInstanceID, map[string]string{
		"name": name,
	})
	if err != nil {
		return err
	}

	if len(networks.Networks) != 1 {
		return errors.New("error coudn't find exact network, please check the name")
	}
	d.SetID(strconv.Itoa(networks.Networks[0].Id))

	// post check
	if err := d.Error(); err != nil {
		return err
	}

	return nil
}
