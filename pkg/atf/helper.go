package atf

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/tidwall/gjson"
)

// validateResource validates the resource exists in state file
func validateResource(name string, validations map[string]interface{}, getApi GetAPIFunc) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("[Validate Resource] resource %s not found", name)
		}

		id := rs.Primary.Attributes["id"]
		if id == "" {
			return fmt.Errorf("resource %s ID is not set", name)
		}

		resp, err := getApi(rs.Primary.Attributes)
		if err != nil {
			return err
		}

		jsonBody, err := json.Marshal(resp)
		if err != nil {
			return err
		}
		jsonStr := string(jsonBody)
		for requestKey, stateValue := range validations {
			result := gjson.Get(jsonStr, requestKey)
			if result.String() != fmt.Sprint(stateValue) {
				return fmt.Errorf("validation failed for %s. On API response, expected %s = %s, but got %s",
					name, requestKey, result.String(), stateValue)
			}
		}

		return nil
	}
}

// getLocalName truncates hpegl_vmaas_ and returns back remaining.
func getLocalName(res string) string {
	return res[len("hpegl_vmaas_"):]
}

func getTag(isResource bool) string {
	if isResource {
		return "resource"
	}
	return "data"
}
