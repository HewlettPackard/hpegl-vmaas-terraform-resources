// (C) Copyright 2023 Hewlett Packard Enterprise Development LP

package utils

import (
	"context"
	"strconv"
	"strings"

	"github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/client"
)

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
