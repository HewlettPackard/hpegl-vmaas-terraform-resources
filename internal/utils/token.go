//  (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package utils

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	sdkConst "github.com/hpe-hcss/vmaas-cmp-go-sdk/pkg/client"
	"github.com/hpe-hcss/vmaas-terraform-resources/pkg/client"
)

// SetScmClientToken fetches and sets the token  in context for scm client.
// Provided the client id and secret in provider
func SetScmClientToken(ctx context.Context, meta interface{}) {
	token, err := client.GetToken(ctx, meta)
	if err != nil {
		log.Print(diag.Warning, "Error in getting token using SCM client: %s", err)
	} else {
		ctx = context.WithValue(ctx, sdkConst.ContextAccessToken, token)
	}
}
