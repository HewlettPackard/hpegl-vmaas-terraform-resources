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
	log               logger.Logger
	serviceInstanceId string
}

func newNetwork(nClient *client.NetworksApiService, serviceInstanceId string) *network {
	return &network{nClient: nClient, serviceInstanceId: serviceInstanceId}
}
func (n *network) Read(ctx context.Context, d *utils.Data) error {
	n.log.Debug("Get Network")

	name := d.GetString("name")
	networks, err := n.nClient.GetAllNetworks(ctx, n.serviceInstanceId, map[string]string{
		"name": name,
	})
	if err != nil {
		return err
	}
	n.log.Info(networks.NetworkCount)
	n.log.Info(networks.Networks)
	if len(networks.Networks) != 1 {
		return errors.New("error coudn't find exact network, please check the name")
	}
	d.SetID(strconv.Itoa(int(networks.Networks[0].Id)))
	if d.HaveError() {
		return err
	}

	return nil
}
