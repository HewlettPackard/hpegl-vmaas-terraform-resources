package utils

import (
	"context"
	"log"

	"github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hewlettpackard/hpegl-provider-lib/pkg/token/retrieve"
	"github.com/hewlettpackard/hpegl-provider-lib/pkg/token/serviceclient"
)

func SetMeta(apiClient *client.APIClient, r *schema.ResourceData) {
	err := apiClient.SetMeta(nil, func(ctx *context.Context, meta interface{}) {
		// Initialise token handler
		h, err := serviceclient.NewHandler(r)
		if err != nil {
			log.Printf("[WARN] Unable to fetch token for SCM client: %s", err)
		}

		// Get token retrieve func and put in c
		trf := retrieve.NewTokenRetrieveFunc(h)
		token, err := trf(*ctx)
		if err != nil {
			log.Printf("[WARN] Unable to fetch token for SCM client: %s", err)
		} else {
			*ctx = context.WithValue(*ctx, client.ContextAccessToken, token)
		}
	})
	if err != nil {
		log.Printf("[WARN] Error: %s", err)
	}
}

// SetMetaFnAndVersion sets the token-generation function and version for the Broker API client
func SetMetaFnAndVersion(apiClient *client.APIClient, r *schema.ResourceData, version int) {
	apiClient.SetMetaFnAndVersion(nil, version, func(ctx *context.Context, meta interface{}) {
		// Initialise token handler
		h, err := serviceclient.NewHandler(r)
		if err != nil {
			log.Printf("[WARN] Unable to fetch token for SCM client: %s", err)
		}

		// Get token retrieve func and put in c
		trf := retrieve.NewTokenRetrieveFunc(h)
		token, err := trf(*ctx)
		if err != nil {
			log.Printf("[WARN] Unable to fetch token for SCM client: %s", err)
		} else {
			*ctx = context.WithValue(*ctx, client.ContextAccessToken, token)
		}
	})
}
func SetCMPVars(apiClient, brokerClient *client.APIClient) (cfg client.Configuration, err error) {
	cmpDetails, err := brokerClient.GetCMPVars(context.Background())
	apiClient.SetHostandToken(cmpDetails.URL, cmpDetails.AccessToken)
	apiClient.SetCMPVersion(context.Background())
	cfg.Host = cmpDetails.URL
	cfg.DefaultQueryParams = map[string]string{}
	cfg.DefaultHeader = map[string]string{"Authorization": "Bearer " + cmpDetails.AccessToken}
	apiClient.SetCMPMeta(nil, brokerClient, func(ctx *context.Context, meta interface{}) {
		// Initialise token handler
		// client.GetClientFromMetaMap(meta)
		cmpDetails, err := brokerClient.GetCMPVars(*ctx)

		if err != nil {
			log.Printf("[ERROR] Unable to fetch token for CMP client: %s", err)

		} else {
			apiClient.SetHostandToken(cmpDetails.URL, cmpDetails.AccessToken)
			// *ctx = context.WithValue(*ctx, client.ContextAccessToken, cmpDetails.AccessToken)
		}
	})
	// if err != nil {
	// 	log.Printf("[WARN] Error: %s", err)
	// }
	return
	// apiClient.SetMetaFnAndVersion(nil, version, func(ctx *context.Context, meta interface{}) {
	// 	// Initialise token handler
	// 	h, err := serviceclient.NewHandler(r)
	// 	if err != nil {
	// 		log.Printf("[WARN] Unable to fetch token for SCM client: %s", err)
	// 	}

	// 	// Get token retrieve func and put in c
	// 	trf := retrieve.NewTokenRetrieveFunc(h)
	// 	token, err := trf(*ctx)
	// 	if err != nil {
	// 		log.Printf("[WARN] Unable to fetch token for SCM client: %s", err)
	// 	} else {
	// 		*ctx = context.WithValue(*ctx, client.ContextAccessToken, token)
	// 	}
	// })
}
