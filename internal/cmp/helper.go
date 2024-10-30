// (C) Copyright 2021-2024 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"
	"strconv"
	"strings"

	"github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/client"
)

func setMeta(meta interface{}, apiClient client.APIClientHandler) {
	// err := apiClient.SetMeta(meta, auth.SetScmClientToken)
	// if err != nil {
	// 	log.Printf("[ERROR] error while setting meta information for cmp-sdk, error: %v", err)
	// }
}

func ParseVersion(version string) (int, error) {
	if version == "" {
		return 0, nil
	}
	version = strings.Split(version, "-")[0]
	versionSplit := strings.Split(version, ".")

	mul := 10000
	sum := 0
	for _, v := range versionSplit {
		vInt, err := strconv.Atoi(v)
		if err != nil {
			return 0, err
		}
		sum += mul * vInt
		mul /= 100
	}

	return sum, nil
}

func GetCmpVersion(ctx context.Context, apiClient client.APIClientHandler) (int, error) {
	c := client.CmpStatus{
		Client: apiClient,
	}

	cmpVersion, err := c.GetCmpVersion(ctx)
	if err != nil {
		return 0, err
	}

	return ParseVersion(cmpVersion.Appliance.BuildVersion)
}

func GetNsxTypeFromCMP(ctx context.Context, apiClient client.APIClientHandler) (string, error) {
	cmpVersion, err := GetCmpVersion(ctx, apiClient)
	if err != nil {
		return "", err
	}
	if v, _ := ParseVersion("6.2.4"); v <= cmpVersion {
		// from 6.2.4 onwards the display name of NSX-T has been change to NSX

		return nsx, nil
	}

	return nsxt, nil
}
