// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package cmp

import (
	"log"

	"github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/client"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/auth"
)

func setMeta(meta interface{}, apiClient client.APIClientHandler) {
	err := apiClient.SetMeta(meta, auth.SetScmClientToken)
	if err != nil {
		log.Printf("[ERROR] error while setting meta information for cmp-sdk, error: %v", err)
	}
}
