// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

// Package constants - constants that are used in pkg/client and pkg/resources
package constants

const (
	// ServiceName - the service mnemonic
	ServiceName = "vmaas"
	ServiceURL  = "https://iac-vmaas.us1.greenlake-hpe.com"

	LOCATION    = "location"
	SPACENAME   = "space_name"
	APIURL      = "api_url"
	INSECURE    = "allow_insecure"
	SpaceKey    = "space"
	LocationKey = "location"

	MockIAMKey     = "TF_ACC_MOCK_IAM"
	CmpSubjectKey  = "TF_ACC_CMP_SUBJECT"
	AccTestPathKey = "TF_ACC_TEST_PATH"
)
