// package atf or acceptance-test-framework consists of helper files to parse,
// validate and run acceptance test with vmaas specified terraform acceptance
// test case format
package atf

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/spf13/viper"
)

type accConfig struct {
	config      string
	validations []validation
}

type validation struct {
	isJson bool
	key    string
	value  interface{}
}

func getViperConfig(name string, isResource bool) *viper.Viper {
	if path := os.Getenv(constants.AccTestPathKey); path != "" {
		accTestPath = path
	}

	v := viper.New()
	v.SetConfigFile(fmt.Sprintf("%s/%s/%s.yaml", accTestPath, getTag(isResource), name))
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Sprint("error while reading config, ", err))
	}

	return v
}

func parseValidations(vls map[interface{}]interface{}) []validation {
	m := make([]validation, 0, len(vls))
	for k, v := range vls {
		kStr := k.(string)
		kSplit := strings.Split(kStr, ".")
		if len(kSplit) > 1 && (kSplit[0] == jsonKey || kSplit[0] == tfKey) {
			isJson := false
			if kSplit[0] == jsonKey {
				isJson = true
			}

			m = append(m, validation{
				isJson: isJson,
				key:    kStr[len(kSplit[0])+1:],
				value:  v,
			})
		} else {
			panic("invalid validation format. validation format should be '[json|tf].key1.key2....keyn: value'")
		}
	}

	return m
}

// parseConfig populates terraform configuration and parse to accConfig
func parseConfig(name string, isResource bool) []accConfig {
	tfKey := getLocalName(name)
	v := getViperConfig(tfKey, isResource)

	testCases := v.Get(accKey).([]interface{})
	configs := make([]accConfig, len(testCases))
	for i := range testCases {
		configs[i].config = fmt.Sprintf(`
		%s
		%s "%s" "tf_%s" {
			%s
		}
		`,
			providerStanza, getType(isResource), name, tfKey,
			v.GetString(fmt.Sprintf(`%s.%d.config`, accKey, i)),
		)

		validations, ok := v.Get(fmt.Sprintf("%s.%d.validations", accKey, i)).(map[interface{}]interface{})
		if !ok {
			continue
		}
		configs[i].validations = parseValidations(validations)
	}

	return configs
}

// getTestCases populate TestSteps
func getTestCases(t *testing.T, name string, getAPI GetAPIFunc, isResource bool) []resource.TestStep {
	configs := parseConfig(name, isResource)
	testSteps := make([]resource.TestStep, 0, len(configs))

	tag := ""
	if !isResource {
		tag = "data."
	}

	for _, c := range configs {
		testSteps = append(testSteps, resource.TestStep{
			Config: c.config,
			Check: resource.ComposeTestCheckFunc(
				validateResource(
					fmt.Sprintf("%s%s.tf_%s", tag, name, getLocalName(name)),
					c.validations,
					getAPI,
				),
			),
		})
	}

	return testSteps
}
