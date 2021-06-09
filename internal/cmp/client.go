// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package cmp

import apiClient "github.com/hpe-hcss/vmaas-cmp-go-sdk/pkg/client"

// Client is the cmp client which will implements all the
// functions in interface.go
type Client struct {
	Instance      Resource
	Snapshot      Resource
	Network       DataSource
	Plan          DataSource
	Group         DataSource
	Cloud         DataSource
	ResourcePool  DataSource
	Layout        DataSource
	Datastore     DataSource
	PowerSchedule DataSource
	Template      DataSource
}

// NewClient returns configured client
func NewClient(client *apiClient.APIClient, cfg apiClient.Configuration) *Client {
	return &Client{
		Instance:      newInstance(&apiClient.InstancesApiService{Client: client, Cfg: cfg}),
		Snapshot:      newSnapshot(&apiClient.InstancesApiService{Client: client, Cfg: cfg}),
		Network:       newNetwork(&apiClient.NetworksApiService{Client: client, Cfg: cfg}),
		Plan:          newPlan(&apiClient.PlansApiService{Client: client, Cfg: cfg}),
		Group:         newGroup(&apiClient.GroupsApiService{Client: client, Cfg: cfg}),
		Layout:        newLayout(&apiClient.LibraryApiService{Client: client, Cfg: cfg}),
		Cloud:         newCloud(&apiClient.CloudsApiService{Client: client, Cfg: cfg}),
		ResourcePool:  newResourcePool(&apiClient.CloudsApiService{Client: client, Cfg: cfg}),
		Datastore:     newDatastore(&apiClient.CloudsApiService{Client: client, Cfg: cfg}),
		PowerSchedule: newPowerSchedule(&apiClient.PowerSchedulesApiService{Client: client, Cfg: cfg}),
		Template:      newTemplate(&apiClient.VirtualImagesApiService{Client: client, Cfg: cfg}),
	}
}
