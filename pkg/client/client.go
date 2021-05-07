// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package client

import (
	"github.com/hpe-hcss/hpegl-provider-lib/pkg/client"
	"github.com/hpe-hcss/hpegl-provider-lib/pkg/provider"
)

// keyForGLClientMap is the key in the map[string]interface{} that is passed down by hpegl used to store *Client
// This must be unique, hpegl will error-out if it isn't
const keyForGLClientMap = "vmaasClient"

// Assert that InitialiseClient satisfies the client.Initialisation interface
var _ client.Initialisation = (*InitialiseClient)(nil)

// Client is the client struct that is used by the provider code
type Client struct {
	IAMToken   string
	VMaaSAPIUrl   string
	VMaaSToken string
}

// InitialiseClient is imported by hpegl from each service repo
type InitialiseClient struct{}

// NewClient takes an argument of all of the provider.ConfigData, and returns an interface{} and error
// If there is no error interface{} will contain *Client.
// The hpegl provider will put *Client at the value of keyForGLClientMap (returned by ServiceName) in
// the map of clients that it creates and passes down to provider code.  hpegl executes NewClient for each service.
func (i InitialiseClient) NewClient(config provider.ConfigData) (interface{}, error) {

	client := new(Client)
	client.IAMToken = config.IAMToken
	client.VMaaSAPIUrl = config.VMaaSAPIUrl // To be replaced with VMaaSAPIUrl

	// Agena API Call
	//


	// Call agena-api pass the IAM Token to get the VMaaSAPIToken
	// client.VMaaSToken = getVMaaSToken(config.IAMToken)
	client.VMaaSToken = ""

	return client, nil
}

// ServiceName is used to return the value of keyForGLClientMap, for use by hpegl
func (i InitialiseClient) ServiceName() string {
	return keyForGLClientMap
}

// GetClientFromMetaMap is a convenience function used by provider code to extract *Client from the
// meta argument passed-in by terraform
func GetClientFromMetaMap(meta interface{}) *Client {
	return meta.(map[string]interface{})[keyForGLClientMap].(*Client)
}
