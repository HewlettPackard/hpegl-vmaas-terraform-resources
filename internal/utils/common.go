// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package utils

import (
	"encoding/json"
	"strconv"
)

func JSONNumber(in interface{}) json.Number {
	if in == nil {
		return json.Number("")
	}
	if a, ok := in.(int); ok {
		return json.Number(strconv.Itoa(a))
	}

	return json.Number(in.(string))
}

func ParseInt(str string) (int64, error) {
	return strconv.ParseInt(str, 10, 64)
}

func ParsePowerState(state string) string {
	switch state {
	case StateRunning:
		return PowerOn
	case StateStopped:
		return PowerOff
	case StateSuspended:
		return Suspend
	}

	return state
}
