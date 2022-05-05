// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package client

import (
	"fmt"
	"os"

	api_client "github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/client"
	cmp_client "github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/cmp"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/constants"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hewlettpackard/hpegl-provider-lib/pkg/client"
)

// keyForGLClientMap is the key in the map[string]interface{} that is passed down by hpegl used to store *Client
// This must be unique, hpegl will error-out if it isn't
const keyForGLClientMap = "vmaasClient"

var serviceURL string

// Assert that InitialiseClient satisfies the client.Initialisation interface
var _ client.Initialisation = (*InitialiseClient)(nil)

// Client is the client struct that is used by the provider code
type Client struct {
	CmpClient *cmp_client.Client
}

// Get env configurations for VmaaS services
func getHeaders() map[string]string {
	token := os.Getenv("HPEGL_IAM_TOKEN")
	header := make(map[string]string)
	if utils.GetEnvBool(constants.MockIAMKey) {
		header["subject"] = os.Getenv(constants.CmpSubjectKey)
		header["Authorization"] = token
	}

	return header
}

// InitialiseClient is imported by hpegl from each service repo
type InitialiseClient struct{}

// NewClient takes an argument of all of the provider.ConfigData, and returns an interface{} and error
// If there is no error interface{} will contain *Client.
// The hpegl provider will put *Client at the value of keyForGLClientMap (returned by ServiceName) in
// the map of clients that it creates and passes down to provider code.  hpegl executes NewClient for each service.
func (i InitialiseClient) NewClient(r *schema.ResourceData) (interface{}, error) {
	vmaasProviderSettings, err := client.GetServiceSettingsMap(constants.ServiceName, r)
	if err != nil {
		return nil, nil
	}

	// Create VMaas Client
	client := new(Client)

	cfg := api_client.Configuration{
		Host:          vmaasProviderSettings[constants.APIURL].(string),
		DefaultHeader: getHeaders(),
		DefaultQueryParams: map[string]string{
			constants.SpaceKey:    vmaasProviderSettings[constants.SPACENAME].(string),
			constants.LocationKey: vmaasProviderSettings[constants.LOCATION].(string),
		},
	}
	apiClient := api_client.NewAPIClient(&cfg)
	utils.SetMeta(apiClient, r)
	client.CmpClient = cmp_client.NewClient(apiClient, cfg)

	return client, nil
}

// ServiceName is used to return the value of keyForGLClientMap, for use by hpegl
func (i InitialiseClient) ServiceName() string {
	return keyForGLClientMap
}

// GetClientFromMetaMap is a convenience function used by provider code to extract *Client from the
// meta argument passed-in by terraform
func GetClientFromMetaMap(meta interface{}) (*Client, error) {
	cli := meta.(map[string]interface{})[keyForGLClientMap]
	if cli == nil {
		return nil, fmt.Errorf("client is not initialised, make sure that vmaas block is defined in hpegl stanza")
	}

	return cli.(*Client), nil
}
