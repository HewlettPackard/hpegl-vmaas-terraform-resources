// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package acceptancetest

import (
	"context"
	"os"
	"strconv"

	api_client "github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/client"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/auth"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/constants"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/utils"
)

func getAPIClient() (*api_client.APIClient, api_client.Configuration) {
	var headers map[string]string
	if utils.GetEnvBool("TF_ACC_MOCK_IAM") {
		headers = make(map[string]string)
		headers["Authorization"] = os.Getenv("HPEGL_IAM_TOKEN")
		headers["subject"] = os.Getenv(constants.CmpSubjectKey)
	}

	cfg := api_client.Configuration{
		Host:          constants.AccServiceURL,
		DefaultHeader: headers,
		DefaultQueryParams: map[string]string{
			constants.SpaceKey:    os.Getenv("HPEGL_VMAAS_LOCATION"),
			constants.LocationKey: os.Getenv("HPEGL_VMAAS_SPACE_NAME"),
		},
	}
	apiClient := api_client.NewAPIClient(&cfg)
	apiClient.SetMeta(nil, auth.SetScmClientToken)

	return apiClient, cfg
}

func toInt(s string) int {
	i, _ := strconv.Atoi(s)

	return i
}

func getAccContext() context.Context {
	return context.Background()
}
