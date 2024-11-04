// (C) Copyright 2024 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"
	"fmt"
	"strconv"
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
	setMetaHpegl(meta, m.bClient.Client)

	// Get Morpheus Tokens and URL
	morpheusDetails, err := m.bClient.GetMorpheusDetails(ctx)
	if err != nil {
		return err
	}

	// Convert the Unix timestamp to Duration in seconds expressed as a string
	validDuration := time.Until(time.UnixMilli(morpheusDetails.ValidTill))
	// We do the following since we cannot get a string representation of a Duration in seconds
	validSeconds := validDuration.Round(time.Second).Seconds() // Round to the nearest second, in float64
	validSecondsString := fmt.Sprintf("%ss", strconv.FormatFloat(validSeconds, 'f', -1, 64))

	// Set all of the details
	d.SetId(morpheusDetails.ID)

	if err = d.Set("access_token", morpheusDetails.AccessToken); err != nil {
		return err
	}

	if err = d.Set("valid_till", validSecondsString); err != nil {
		return err
	}

	if err = d.Set("url", morpheusDetails.URL); err != nil {
		return err
	}

	return nil
}
