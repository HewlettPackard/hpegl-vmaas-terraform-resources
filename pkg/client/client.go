// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package client

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hewlettpackard/hpegl-provider-lib/pkg/client"
	"github.com/hewlettpackard/hpegl-provider-lib/pkg/token/common"
	"github.com/hewlettpackard/hpegl-provider-lib/pkg/token/retrieve"

	api_client "github.com/hpe-hcss/vmaas-cmp-go-sdk/pkg/client"
	cmp_client "github.com/hpe-hcss/vmaas-terraform-resources/internal/cmp"
	"github.com/hpe-hcss/vmaas-terraform-resources/pkg/constants"
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
func getHeaders(token, location, spaceName string) map[string]string {
	header := make(map[string]string)
	serviceURL = constants.ServiceURL
	if strings.ToLower(os.Getenv("TF_ACC")) == "true" {
		serviceURL = constants.AccServiceURL
		header["subject"] = os.Getenv("CMP_SUBJECT")
	}
	if strings.ToLower(os.Getenv("SERVICE_ACCOUNT")) == "intg" {
		serviceURL = constants.IntgServiceURL
	}
	header["Authorization"] = token
	header["location"] = location
	header["space"] = spaceName

	return header
}

// InitialiseClient is imported by hpegl from each service repo
type InitialiseClient struct{}

// NewClient takes an argument of all of the provider.ConfigData, and returns an interface{} and error
// If there is no error interface{} will contain *Client.
// The hpegl provider will put *Client at the value of keyForGLClientMap (returned by ServiceName) in
// the map of clients that it creates and passes down to provider code.  hpegl executes NewClient for each service.
func (i InitialiseClient) NewClient(r *schema.ResourceData) (interface{}, error) {
	token := r.Get("iam_token").(string)
	vmaasProviderSettings, err := client.GetServiceSettingsMap(constants.ServiceName, r)
	if err != nil {
		return nil, err
	}

	// Read the value supplied in the tf file
	location := vmaasProviderSettings[constants.LOCATION].(string)
	spaceName := vmaasProviderSettings[constants.SPACENAME].(string)
	allowInsecure := vmaasProviderSettings[constants.INSECURE].(bool)

	// Create VMaas Client
	client := new(Client)

	cfg := api_client.Configuration{
		Host:          serviceURL,
		DefaultHeader: getHeaders(token, location, spaceName),
	}
	apiClient := api_client.NewAPIClient(&cfg, !allowInsecure)
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

// GetToken is a convenience function used by provider code to extract retrieve.TokenRetrieveFuncCtx from
// the meta argument passed-in by terraform and execute it with the context ctx
func GetToken(ctx context.Context, meta interface{}) (string, error) {
	trf := meta.(map[string]interface{})[common.TokenRetrieveFunctionKey].(retrieve.TokenRetrieveFuncCtx)

	return trf(ctx)
}

// SetScmClientToken fetches and sets the token  in context for scm client.
// Provided the client id and secret in provider
func SetScmClientToken(ctx *context.Context, meta interface{}) {
	token, err := GetToken(*ctx, meta)
	if err != nil {
		log.Printf("[WARN] Unable to fetch token for SCM client: %s", err)
	} else {
		*ctx = context.WithValue(*ctx, api_client.ContextAccessToken, token)
	}
}
