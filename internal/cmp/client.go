// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package cmp

import apiClient "github.com/hpe-hcss/vmaas-cmp-go-sdk/pkg/client"

// Client is the cmp client which will implements all the
// functions in interface.go
type Client struct {
	// Instance resource
	Instance Resource
	// Network data source
	Network DataSource
	// Plan data source
	Plan DataSource
}

// NewClient returns configured client
func NewClient(client *apiClient.APIClient, cfg apiClient.Configuration) *Client {
	return &Client{
		Instance: newInstance(&apiClient.InstancesApiService{Client: client, Cfg: cfg}, ""),
		Network:  newNetwork(&apiClient.NetworksApiService{Client: client, Cfg: cfg}, ""),
		Plan:     newPlan(&apiClient.PlansApiService{Client: client, Cfg: cfg}, ""),
	}
}
