// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package utils

import (
	"encoding/json"

	api_client "github.com/hpe-hcss/vmaas-cmp-go-sdk/pkg/client"
)

func parseError(err error) api_client.CustomError {
	customErr := api_client.CustomError{}
	if err == nil {
		return customErr
	}
	jsonErr := json.Unmarshal([]byte(err.Error()), &customErr)
	if err != nil {
		customErr.Errors = jsonErr.Error()
	}

	return customErr
}

func GetStatusCode(err error) int {
	customErr := parseError(err)

	return customErr.StatusCode
}
