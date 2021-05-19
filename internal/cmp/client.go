// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package cmp

import apiClient "github.com/hpe-hcss/vmaas-cmp-go-sdk/pkg/client"

// Client is the cmp client which will implements all the
// functions in interface.go
type Client struct {
	Instance
}

// NewClient returns configured client
func NewClient(apiClient *apiClient.APIClient) *Client {
	return &Client{
		Instance: &instance{
			client: apiClient,
		},
	}
}
