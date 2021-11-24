package atf

import "github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/constants"

const providerStanza = `
	provider hpegl {
		vmaas {
			space_name = "` + constants.AccSpace + `"
			location = "` + constants.AccLocation + `"
		}
	}

`
