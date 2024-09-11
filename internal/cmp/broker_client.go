// (C) Copyright 2024 Hewlett Packard Enterprise Development LP

package cmp

import apiClient "github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/client"

// BrokerClient - struct to hold the broker client details
type BrokerClient struct {
	DSMorpheusDetails DataSource
}

// NewBrokerClient - function to create a new broker client
func NewBrokerClient(client *apiClient.APIClient, cfg apiClient.Configuration) *BrokerClient {
	return &BrokerClient{
		DSMorpheusDetails: newMorpheusBroker(
			&apiClient.BrokerAPIService{Client: client, Cfg: cfg},
		),
	}
}
