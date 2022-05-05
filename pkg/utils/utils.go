// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package utils

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"testing"

	api_client "github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/client"
	"github.com/spf13/viper"
)

type ResourceData struct {
	Data map[string]interface{}
}

func (r *ResourceData) Get(key string) interface{} {
	return r.Data[key]
}

func parseError(err error) api_client.CustomError {
	customErr := api_client.CustomError{}
	if err == nil {
		return customErr
	}
	jsonErr := json.Unmarshal([]byte(err.Error()), &customErr)
	if jsonErr != nil {
		customErr.Errors = jsonErr.Error()
	}

	return customErr
}

func GetStatusCode(err error) int {
	customErr := parseError(err)

	return customErr.StatusCode
}

var skipMap map[string]bool

func ReadSkip() {
	skipMap = make(map[string]bool)
	skipItems, ok := viper.Get("skip").([]interface{})
	if !ok {
		return
	}
	// loop over skip and store into skipMap
	for _, val := range skipItems {
		skipMap[val.(string)] = true
	}
}

func SkipAcc(t *testing.T, param string) {
	_, ok := skipMap[param]
	if ok {
		t.Skip("Acceptance test for " + param + " has been skipped...")
	}
}

func GetEnvBool(key string) bool {
	envVal := os.Getenv(key)
	if envVal == "" {
		return false
	}

	value, err := strconv.ParseBool(envVal)
	if err != nil {
		log.Printf("[ERROR] failed to parse env %s, got error %v", key, err)

		return false
	}

	return value
}
