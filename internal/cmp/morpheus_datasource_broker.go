// (C) Copyright 2024 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"
	"github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/client"

	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/utils"
)

// morpheusBroker is used to read morpheus details using the Broker API
type morpheusBroker struct {
	bClient *client.BrokerAPIService
}

func newMorpheusBroker(bClient *client.BrokerAPIService) *morpheusBroker {
	return &morpheusBroker{bClient: bClient}
}

// Read reads the morpheus details using the Broker API
func (m *morpheusBroker) Read(ctx context.Context, d *utils.Data, meta interface{}) error {
	setMeta(meta, m.bClient.Client)

	// Get Morpheus Tokens and URL
	morpheusDetails, err := m.bClient.GetMorpheusDetails(ctx)
	if err != nil {
		return err
	}

	// Set access_token, refresh_token and morpheus_url
	if err = d.Set("access_token", morpheusDetails.AccessToken); err != nil {
		return err
	}

	if err = d.Set("refresh_token", morpheusDetails.RefreshToken); err != nil {
		return err
	}

	if err = d.Set("morpheus_url", morpheusDetails.RefreshToken); err != nil {
		return err
	}

	return nil
}
