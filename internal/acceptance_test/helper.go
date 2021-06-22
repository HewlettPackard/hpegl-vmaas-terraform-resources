// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package acceptancetest

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func validateDataSourceID(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("data source not found %s", name)
		}

		id := rs.Primary.Attributes["id"]
		if id == "" {
			return fmt.Errorf("data source %s is not set", name)
		}

		return nil
	}
}

const providerStanza = `
	provider hpegl {
		vmaas {
			allow_insecure = true
			space_name = "tf_acceptance"
			location = "tf_acc_location"
		}
	}

`
