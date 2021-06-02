// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package utils

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/hpe-hcss/vmaas-terraform-resources/internal/logger"
)

const (
	defaultTimeout    = time.Second * 5
	defaultRetryCount = 3
)

func JSONNumber(in interface{}) json.Number {
	if a, ok := in.(int); ok {
		return json.Number(strconv.Itoa(a))
	}

	return json.Number(in.(string))
}

func ParseInt(str string) (int64, error) {
	return strconv.ParseInt(str, 10, 64)
}

func Retry(fn func() (interface{}, error)) (interface{}, error) {
	var err error
	var resp interface{}
	for i := 0; i < defaultRetryCount; i++ {
		resp, err = fn()
		if err == nil {
			break
		}
		logger.Error("error inside retry, ", err, ". Response: ", resp)
		time.Sleep(defaultTimeout)
	}

	return resp, err
}
