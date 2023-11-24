// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/client"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/auth"
)

func setMeta(meta interface{}, apiClient client.APIClientHandler) {
	err := apiClient.SetMeta(meta, auth.SetScmClientToken)
	if err != nil {
		log.Printf("[ERROR] error while setting meta information for cmp-sdk, error: %v", err)
	}
}

func ParseVersion(version string) (int, error) {
	if version == "" {
		return 0, nil
	}

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

func GetCmpVersion(apiClient client.APIClientHandler) (int, error) {
	c := client.CmpStatus{
		Client: apiClient,
	}

	cmpVersion, err := c.GetCmpVersion(context.Background())
	if err != nil {
		return 0, err
	}

	return ParseVersion(cmpVersion.Appliance.BuildVersion)
}

func GetNsxTypeFromCMP(apiClient client.APIClientHandler) (string, error) {
	cmpVersion, err := GetCmpVersion(apiClient)
	if err != nil {
		return "", err
	}
	if v, _ := ParseVersion("6.2.4"); v <= cmpVersion {
		// from 6.2.4 onwards the display name of NSX-T has been change to NSX

		return nsx, nil
	}

	return nsxt, nil
}
