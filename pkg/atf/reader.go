// package atf or acceptance-test-framework consists of helper files to parse,
// validate and run acceptance test with vmaas specified terraform acceptance
// test case format
package atf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/spf13/viper"
)

type accConfig struct {
	Config      string
	Validations map[string]interface{}
}

func getResourceConfig(resourceName string) []accConfig {
	tfResKey := getLocalName(resourceName)
	resKey := "vmaas.resource." + tfResKey
	testCases := viper.Get(resKey).([]interface{})
	configs := make([]accConfig, len(testCases))
	for i := range testCases {
		configs[i].Config = fmt.Sprintf(`
		%s
		resource "%s" "tf_%s" {
			%s
		}
		`,
			providerStanza, resourceName, tfResKey,
			viper.GetString(fmt.Sprintf(`%s.%d.config`, resKey, i)),
		)

		configs[i].Validations = make(map[string]interface{})
		validations, ok := viper.Get(fmt.Sprintf("%s.%d.validations", resKey, i)).(map[interface{}]interface{})
		if !ok {
			continue
		}
		for k, v := range validations {
			configs[i].Validations[k.(string)] = v
		}
	}

	return configs
}

func getTestCases(t *testing.T, resourceName string) []resource.TestStep {
	configs := getResourceConfig(resourceName)
	testSteps := make([]resource.TestStep, 0, len(configs))
	for _, c := range configs {
		testSteps = append(testSteps, resource.TestStep{
			Config: c.Config,
			Check: resource.ComposeTestCheckFunc(
				validateResource(
					fmt.Sprintf("%s.tf_%s", resourceName, getLocalName(resourceName)),
				),
			),
		})
	}

	return testSteps
}
