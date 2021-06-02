// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package utils

const (
	StateUnknown      = "unknown"
	StateRunning      = "running"
	StateFailed       = "failed"
	StateProvisioning = "provisioning"
)

var powerStateMap = map[string]string{
	StateRunning:      "powerOn",
	StateProvisioning: "pending",
}

// return power state according to status
func GetPowerState(key string) string {
	if val, ok := powerStateMap[key]; ok {
		return val
	}
	return StateUnknown
}
