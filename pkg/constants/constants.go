// (C) Copyright 2021-2024 Hewlett Packard Enterprise Development LP

// Package constants - constants that are used in pkg/client and pkg/resources
package constants

const (
	// ServiceName - the service mnemonic
	ServiceName        = "vmaas"
	ServiceURL         = "https://iac-vmaas.us1.greenlake-hpe.com"
	BrokerURL          = "https://vmaas-broker.us1.greenlake-hpe.com"
	IamGlp      string = "glp"
	IamGlcs     string = "glcs"
	TenantID    string = "tenant_id"

	LOCATION       = "location"
	SPACENAME      = "space_name"
	APIURL         = "api_url"
	BROKERRURL     = "broker_url"
	INSECURE       = "allow_insecure"
	MORPHEUS_URL   = "morpheus_url"
	MORPHEUS_TOKEN = "morpheus_token"
	SpaceKey       = "space"
	TenantIDKey    = "tenantID"
	LocationKey    = "location"

	MockIAMKey     = "TF_ACC_MOCK_IAM"
	CmpSubjectKey  = "TF_ACC_CMP_SUBJECT"
	AccTestPathKey = "TF_ACC_TEST_PATH"
)
