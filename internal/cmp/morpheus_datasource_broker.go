// (C) Copyright 2024 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"
	"time"

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

	// Convert the Unix timestamp to Duration in seconds
	validSeconds := time.Until(time.Unix(morpheusDetails.ValidTill, 0)) / time.Second

	// Set all of the details
	d.SetId(morpheusDetails.ID)

	if err = d.Set("access_token", morpheusDetails.AccessToken); err != nil {
		return err
	}

	if err = d.Set("valid_till", validSeconds); err != nil {
		return err
	}

	if err = d.Set("url", morpheusDetails.URL); err != nil {
		return err
	}

	return nil
}
