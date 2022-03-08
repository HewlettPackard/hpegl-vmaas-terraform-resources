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
