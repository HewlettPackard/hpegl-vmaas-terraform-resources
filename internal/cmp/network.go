// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"
	"fmt"

	"github.com/hpe-hcss/vmaas-cmp-go-sdk/pkg/client"
	"github.com/hpe-hcss/vmaas-cmp-go-sdk/pkg/models"
	"github.com/hpe-hcss/vmaas-terraform-resources/internal/logger"
	"github.com/hpe-hcss/vmaas-terraform-resources/internal/utils"
	"github.com/hpe-hcss/vmaas-terraform-resources/pkg/auth"
)

type network struct {
	nClient *client.NetworksAPIService
}

func newNetwork(nClient *client.NetworksAPIService) *network {
	return &network{nClient: nClient}
}

func (n *network) Read(ctx context.Context, d *utils.Data, meta interface{}) error {
	logger.Debug("Get Network")

	name := d.GetString("name")
	// Pre check
	if err := d.Error(); err != nil {
		return err
	}
	resp, err := utils.Retry(func() (interface{}, error) {
		auth.SetScmClientToken(&ctx, meta)
		return n.nClient.GetAllNetworks(ctx, nil)
	})
	if err != nil {
		return err
	}

	isMatch := false
	networks := resp.(models.ListNetworksBody)
	for i, n := range networks.Networks {
		if n.DisplayName == name {
			isMatch = true
			d.SetID(networks.Networks[i].ID)

			break
		}
	}
	if !isMatch {
		return fmt.Errorf(errExactMatch, "Network")
	}

	// post check
	return d.Error()
}
