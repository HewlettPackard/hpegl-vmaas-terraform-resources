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

// parseConfig populates terraform configuration and parse to accConfig
func parseConfig(name string, isResource bool) []accConfig {
	tfKey := getLocalName(name)

	resKey := fmt.Sprintf("vmaas.%s.%s", getTag(isResource), tfKey)
	testCases := viper.Get(resKey).([]interface{})
	configs := make([]accConfig, len(testCases))
	for i := range testCases {
		configs[i].Config = fmt.Sprintf(`
		%s
		%s "%s" "tf_%s" {
			%s
		}
		`,
			providerStanza, getTag(isResource), name, tfKey,
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

// getTestCases populate TestSteps
func getTestCases(t *testing.T, name string, getAPI GetAPIFunc, isResource bool) []resource.TestStep {
	configs := parseConfig(name, isResource)
	testSteps := make([]resource.TestStep, 0, len(configs))
	for _, c := range configs {
		testSteps = append(testSteps, resource.TestStep{
			Config: c.Config,
			Check: resource.ComposeTestCheckFunc(
				validateResource(
					fmt.Sprintf("%s.tf_%s", name, getLocalName(name)),
					c.Validations,
					getAPI,
				),
			),
		})
	}

	return testSteps
}
