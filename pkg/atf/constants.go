package atf

const providerStanza = `
	provider hpegl {
		vmaas {
		}
	}
`

var accTestPath = "../../acc-testcases"

const (
	accKey  = "acc"
	jsonKey = "json"
	tfKey   = "tf"

	randMaxLimit = 999999
)
