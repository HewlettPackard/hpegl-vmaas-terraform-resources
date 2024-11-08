// (C) Copyright 2021-2024 Hewlett Packard Enterprise Development LP

package acceptancetest

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	api_client "github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/client"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/constants"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/utils"
	"github.com/hewlettpackard/hpegl-provider-lib/pkg/token/retrieve"
	"github.com/hewlettpackard/hpegl-provider-lib/pkg/token/serviceclient"
)

// func getOldAPIClient() (*api_client.APIClient, api_client.Configuration) {
// 	headers, queryParam, iamVersion := getHeadersAndQueryParamsAndIAMVersion()
// 	cfg := api_client.Configuration{
// 		Host:               os.Getenv("HPEGL_VMAAS_API_URL"),
// 		DefaultHeader:      headers,
// 		DefaultQueryParams: queryParam,
// 	}
// 	apiClient := api_client.NewAPIClient(&cfg)
// 	err := apiClientSetMeta(apiClient, iamVersion)
// 	if err != nil {
// 		log.Printf("[WARN] Error: %s", err)
// 	}

// 	return apiClient, cfg
// }

func getBrokerAPIClient() (*api_client.APIClient, api_client.Configuration) {
	headers, queryParam, iamVersion := getHeadersAndQueryParamsAndIAMVersion()

	cfg := api_client.Configuration{
		Host:               os.Getenv("HPEGL_VMAAS_BROKER_URL"),
		DefaultHeader:      headers,
		DefaultQueryParams: queryParam,
	}
	brokerAPIClient := api_client.NewAPIClient(&cfg)
	err := apiClientSetMeta(brokerAPIClient, iamVersion)
	if err != nil {
		log.Printf("[WARN] Error: %s", err)
	}

	return brokerAPIClient, cfg
}

func getAPIClient() (*api_client.APIClient, api_client.Configuration) {
	ctx := context.Background()
	brokerClient, _ := getBrokerAPIClient()
	cmpDetails, err := brokerClient.GetCMPDetails(ctx)
	if err != nil {
		log.Printf("[ERROR] Error getting cmp details: %s", err)
	}
	cfg := api_client.Configuration{
		Host:               cmpDetails.URL,
		DefaultHeader:      map[string]string{},
		DefaultQueryParams: map[string]string{},
	}
	cmpAPIClient := api_client.NewAPIClient(&cfg)
	cmpAPIClient.CMPToken = cmpDetails.AccessToken
	cmpAPIClient.TokenExpiry = cmpDetails.ValidTill
	cmpAPIClient.SetMetaFnAndVersion(nil, 0, func(ctx *context.Context, meta interface{}) {
		// Initialise token handler
		tokenExpiry := cmpAPIClient.TokenExpiry / 1000
		token := cmpAPIClient.CMPToken
		// Token is about to expire and get new
		if tokenExpiry < time.Now().Unix() {
			cmpDetails, err := brokerClient.GetCMPDetails(*ctx)
			if err != nil {
				log.Printf("[ERROR] Unable to fetch token for CMP client: %s", err)
				panic("Unable to renew token")
			} else {
				token = cmpDetails.AccessToken
				cmpAPIClient.CMPToken = cmpDetails.AccessToken
				cmpAPIClient.TokenExpiry = cmpDetails.ValidTill
			}

		}
		*ctx = context.WithValue(*ctx, api_client.ContextAccessToken, token)
	})
	ctx = context.WithValue(ctx, api_client.ContextAccessToken, cmpDetails.AccessToken)
	err = cmpAPIClient.SetCMPVersion(ctx)
	if err != nil {
		log.Printf("[ERROR] Unable to set CMP version client: %s", err)

	}
	return cmpAPIClient, cfg
}

func getHeadersAndQueryParamsAndIAMVersion() (map[string]string, map[string]string, string) {
	var headers map[string]string
	if utils.GetEnvBool("TF_ACC_MOCK_IAM") {
		headers = make(map[string]string)
		headers["Authorization"] = os.Getenv("HPEGL_IAM_TOKEN")
		headers["subject"] = os.Getenv(constants.CmpSubjectKey)
	}
	iamVersion := utils.GetEnv("HPEGL_IAM_VERSION", constants.IamGlcs)
	queryParam := map[string]string{
		constants.LocationKey: os.Getenv("HPEGL_VMAAS_LOCATION"),
	}
	if iamVersion == constants.IamGlp {
		queryParam[constants.WorkspaceKey] = os.Getenv("HPEGL_VMAAS_SPACE_NAME")
	} else {
		queryParam[constants.SpaceKey] = os.Getenv("HPEGL_VMAAS_SPACE_NAME")
	}

	return headers, queryParam, iamVersion
}

func apiClientSetMeta(apiClient *api_client.APIClient, iamVersion string) error {
	return apiClient.SetMeta(nil, func(ctx *context.Context, meta interface{}) {
		d := &utils.ResourceData{
			Data: map[string]interface{}{
				"iam_service_url":           os.Getenv("HPEGL_IAM_SERVICE_URL"),
				"tenant_id":                 os.Getenv("HPEGL_TENANT_ID"),
				"user_id":                   os.Getenv("HPEGL_USER_ID"),
				"user_secret":               os.Getenv("HPEGL_USER_SECRET"),
				"api_vended_service_client": true,
				"iam_token":                 os.Getenv("HPEGL_IAM_TOKEN"),
				"iam_version":               iamVersion,
			},
		}
		if utils.GetEnvBool(constants.MockIAMKey) {
			return
		}

		// Initialise token handler
		h, err := serviceclient.NewHandler(d)
		if err != nil {
			log.Printf("[WARN] Unable to fetch token for SCM client: %s", err)
		}

		// Get token retrieve func and put in c
		trf := retrieve.NewTokenRetrieveFunc(h)
		token, err := trf(*ctx)
		if err != nil {
			log.Printf("[WARN] Unable to fetch token for SCM client: %s", err)
		} else {
			*ctx = context.WithValue(*ctx, api_client.ContextAccessToken, token)
		}
	})
}

func toInt(s string) int {
	i, _ := strconv.Atoi(s)

	return i
}

func getAccContext() context.Context {
	return context.Background()
}
