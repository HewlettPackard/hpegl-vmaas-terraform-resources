// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package cmp

import apiClient "github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/client"

// Client is the cmp client which will implements all the
// functions in interface.go
type Client struct {
	Instance         Resource
	InstanceClone    Resource
	Router           Resource
	ResNetwork       Resource
	Network          DataSource
	NetworkType      DataSource
	NetworkPool      DataSource
	Plan             DataSource
	Group            DataSource
	Cloud            DataSource
	ResourcePool     DataSource
	Layout           DataSource
	Datastore        DataSource
	PowerSchedule    DataSource
	Template         DataSource
	Environment      DataSource
	NetworkInterface DataSource
	CloudFolder      DataSource
}

// NewClient returns configured client
func NewClient(client *apiClient.APIClient, cfg apiClient.Configuration) *Client {
	return &Client{
		Instance: newInstance(
			&apiClient.InstancesAPIService{Client: client, Cfg: cfg},
			&apiClient.ServersAPIService{Client: client, Cfg: cfg},
		),
		InstanceClone: newInstanceClone(
			&apiClient.InstancesAPIService{Client: client, Cfg: cfg},
			&apiClient.ServersAPIService{Client: client, Cfg: cfg},
		),
		ResNetwork:    newResNetwork(&apiClient.NetworksAPIService{Client: client, Cfg: cfg}),
		Router:        newRouter(&apiClient.RouterAPIService{Client: client, Cfg: cfg}),
		Network:       newNetwork(&apiClient.NetworksAPIService{Client: client, Cfg: cfg}),
		NetworkType:   newNetworkType(&apiClient.NetworksAPIService{Client: client, Cfg: cfg}),
		NetworkPool:   newNetworkPool(&apiClient.NetworksAPIService{Client: client, Cfg: cfg}),
		Plan:          newPlan(&apiClient.PlansAPIService{Client: client, Cfg: cfg}),
		Group:         newGroup(&apiClient.GroupsAPIService{Client: client, Cfg: cfg}),
		Layout:        newLayout(&apiClient.LibraryAPIService{Client: client, Cfg: cfg}),
		Cloud:         newCloud(&apiClient.CloudsAPIService{Client: client, Cfg: cfg}),
		ResourcePool:  newResourcePool(&apiClient.CloudsAPIService{Client: client, Cfg: cfg}),
		Datastore:     newDatastore(&apiClient.CloudsAPIService{Client: client, Cfg: cfg}),
		PowerSchedule: newPowerSchedule(&apiClient.PowerSchedulesAPIService{Client: client, Cfg: cfg}),
		Template:      newTemplate(&apiClient.VirtualImagesAPIService{Client: client, Cfg: cfg}),
		Environment:   newEnvironment(&apiClient.EnvironmentAPIService{Client: client, Cfg: cfg}),
		NetworkInterface: newNetworkInterface(&apiClient.CloudsAPIService{Client: client, Cfg: cfg},
			&apiClient.ProvisioningAPIService{Client: client, Cfg: cfg}),
		CloudFolder: newCloudFolder(&apiClient.CloudsAPIService{Client: client, Cfg: cfg}),
	}
}
