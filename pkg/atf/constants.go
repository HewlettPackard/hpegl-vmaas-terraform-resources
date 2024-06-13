package atf

import "fmt"

// const providerStanza = `
// 	provider hpegl {
// 		vmaas {
// 		}
// 	}
// `

var accTestPath = "../../acc-testcases"

const (
	accKey  = "acc"
	jsonKey = "json"
	tfKey   = "tf"

	randMaxLimit = 999999
)

func GetProviderStanza(alias string) string {
	return fmt.Sprintf(`
	provider hpegl {
	   alias = "%s"
		vmaas {
		}
	}
`, alias)
}
