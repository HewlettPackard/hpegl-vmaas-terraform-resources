package atf

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

type validators func(*terraform.ResourceState) error

func validateResource(name string, v ...validators) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("[Validate Resource] resource %s not found", name)
		}

		id := rs.Primary.Attributes["id"]
		if id == "" {
			return fmt.Errorf("resource %s ID is not set", name)
		}
		for _, vs := range v {
			if err := vs(rs); err != nil {
				return fmt.Errorf("resource %s validation failed with error %w", name, err)
			}
		}

		return nil
	}
}

func getLocalName(res string) string {
	return res[len("hpegl_vmaas_"):]
}
