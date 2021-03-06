//  (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package auth

import (
	"context"
	"log"

	"github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/client"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/constants"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/utils"
	"github.com/hewlettpackard/hpegl-provider-lib/pkg/token/common"
	"github.com/hewlettpackard/hpegl-provider-lib/pkg/token/retrieve"
)

// GetToken is a convenience function used by provider code to extract retrieve.TokenRetrieveFuncCtx from
// the meta argument passed-in by terraform and execute it with the context ctx
func GetToken(ctx context.Context, meta interface{}) (string, error) {
	trf := meta.(map[string]interface{})[common.TokenRetrieveFunctionKey].(retrieve.TokenRetrieveFuncCtx)

	return trf(ctx)
}

// SetScmClientToken fetches and sets the token  in context for scm client.
// Provided the client id and secret in provider
func SetScmClientToken(ctx *context.Context, meta interface{}) {
	if utils.GetEnvBool(constants.MockIAMKey) {
		return
	}

	token, err := GetToken(*ctx, meta)
	if err != nil {
		log.Printf("[WARN] Unable to fetch token for SCM client: %s", err)
	} else {
		*ctx = context.WithValue(*ctx, client.ContextAccessToken, token)
	}
}
