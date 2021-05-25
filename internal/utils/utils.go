// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package utils

import (
	"encoding/json"
	"strconv"
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
