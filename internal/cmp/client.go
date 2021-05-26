// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package cmp

import apiClient "github.com/hpe-hcss/vmaas-cmp-go-sdk/pkg/client"

// Client is the cmp client which will implements all the
// functions in interface.go
type Client struct {
	Instance     Resource
	Network      DataSource
	Plan         DataSource
	Group        DataSource
	ResourcePool DataSource
}

// NewClient returns configured client
func NewClient(client *apiClient.APIClient, cfg apiClient.Configuration, sID string) *Client {
	return &Client{
		Instance:     newInstance(&apiClient.InstancesApiService{Client: client, Cfg: cfg}, sID),
		Network:      newNetwork(&apiClient.NetworksApiService{Client: client, Cfg: cfg}, sID),
		Plan:         newPlan(&apiClient.PlansApiService{Client: client, Cfg: cfg}, sID),
		Group:        newGroup(&apiClient.GroupsApiService{Client: client, Cfg: cfg}, sID),
		ResourcePool: newResourcePool(&apiClient.CloudsApiService{Client: client, Cfg: cfg}, sID),
	}
}
