// package atf or acceptance-test-framework consists of helper files to parse,
// validate and run acceptance test with vmaas specified terraform acceptance
// test case format
package atf

import (
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/spf13/viper"
)

type accConfig struct {
	config      string
	validations []validation
}

type validation struct {
	isJSON bool
	key    string
	value  interface{}
}

func getViperConfig(name, version string, isResource bool) *viper.Viper {
	if path := os.Getenv(constants.AccTestPathKey); path != "" {
		accTestPath = path
	}
	var postfix string
	if version != "" {
		postfix = fmt.Sprintf("-%s", version)
	}
	v := viper.New()
	v.SetConfigFile(fmt.Sprintf("%s/%s/%s%s.yaml", accTestPath, getTag(isResource), name, postfix))
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Sprint("error while reading config, ", err))
	}

	return v
}

func parseMeta(data string) string {
	exp := `%(rand_int)(\{[0-9]+,[0-9]+\})?`
	reg := regexp.MustCompile(exp)

	matches := reg.FindAllString(data, -1)
	src := rand.NewSource(time.Now().Unix())
	r := rand.New(src)
	var randInt int
	for _, m := range matches {
		offReg := regexp.MustCompile(`[0-9]+,[0-9]`)
		numStr := offReg.FindString(m)
		if numStr != "" {
			intSplit := strings.Split(numStr, ",")
			n1 := toInt(intSplit[0])
			n2 := toInt(intSplit[1])
			randInt = r.Intn(n2-n1) + n1
		} else {
			randInt = r.Intn(randMaxLimit)
		}
		data = strings.Replace(data, m, strconv.Itoa(randInt), 1)
	}

	return data
}

func parseValidations(vls map[interface{}]interface{}) []validation {
	m := make([]validation, 0, len(vls))
	for k, v := range vls {
		kStr := k.(string)
		kSplit := strings.Split(kStr, ".")
		if len(kSplit) > 1 && (kSplit[0] == jsonKey || kSplit[0] == tfKey) {
			isJSON := false
			if kSplit[0] == jsonKey {
				isJSON = true
			}

			m = append(m, validation{
				isJSON: isJSON,
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
func parseConfig(v *viper.Viper, name string, isResource bool) []accConfig {
	tfKey := getLocalName(name)

	testCases := v.Get(accKey).([]interface{})
	configs := make([]accConfig, len(testCases))
	for i := range testCases {
		tfConfig := parseMeta(v.GetString(fmt.Sprintf(`%s.%d.config`, accKey, i)))
		configs[i].config = fmt.Sprintf(`
		%s
		%s "%s" "tf_%s" {
			%s
		}
		`,
			providerStanza, getType(isResource), name, tfKey, tfConfig,
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
func getTestCases(t *testing.T, name, version string, getAPI GetAPIFunc, isResource bool) []resource.TestStep {
	v := getViperConfig(getLocalName(name), version, isResource)
	if v.GetBool("ignore") {
		t.Skip("ignoring tests for resource ", name)
	}

	configs := parseConfig(v, name, isResource)
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
