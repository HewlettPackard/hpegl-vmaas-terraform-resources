// package atf or acceptance-test-framework consists of helper files to parse,
// validate and run acceptance test with vmaas specified terraform acceptance
// test case format
package atf

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
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
	isJSON bool
	key    string
	value  interface{}
}

type reader struct {
	t           *testing.T
	isResource  bool
	name        string
	expectError *regexp.Regexp
}

func newReader(t *testing.T, isResource bool, name string) *reader {
	return &reader{
		t:          t,
		isResource: isResource,
		name:       name,
	}
}

func (r *reader) fatal(format string, v ...interface{}) {
	r.t.Fatalf("[acc-test] test case for resource "+r.name+" failed. "+format, v...)
}

func (r *reader) getViperConfig(version string) *viper.Viper {
	tfName := getLocalName(r.name)
	if path := os.Getenv(constants.AccTestPathKey); path != "" {
		accTestPath = path
	}
	var postfix string
	if version != "" {
		postfix = fmt.Sprintf("-%s", version)
	}
	v := viper.New()
	v.SetConfigFile(fmt.Sprintf("%s/%s/%s%s.yaml", accTestPath, getTag(r.isResource), tfName, postfix))
	err := v.ReadInConfig()
	if err != nil {
		r.fatal("error while reading config, %v", err)
	}

	return v
}

func parseMeta(data string) string {
	exp := `%(rand_int)(\{[0-9]+,[0-9]+\})?`
	reg := regexp.MustCompile(exp)

	matches := reg.FindAllString(data, -1)
	var randInt int
	r := newRand()
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

func (r *reader) parseValidations(vip *viper.Viper, i int) []validation {
	vls, ok := vip.Get(fmt.Sprintf("%s.%d.validations", accKey, i)).(map[interface{}]interface{})
	if !ok {
		return nil
	}
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
			r.fatal("invalid validation format. validation format should be '[json|tf].key1.key2....keyn: value'")
		}
	}

	return m
}

func (r *reader) parseRegex(v *viper.Viper, i int) {
	if expectErrStr := v.GetString(fmt.Sprintf("%s.%d.expect_error", accKey, i)); expectErrStr != "" {
		var err error
		r.expectError, err = regexp.Compile(expectErrStr)
		if err != nil {
			r.fatal("error while compiling regex %s, got error %v", expectErrStr, err)
		}
	}
}

// parseConfig populates terraform configuration and parse to accConfig
func (r *reader) parseConfig(v *viper.Viper) []accConfig {
	tfKey := getLocalName(r.name)

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
			providerStanza, getType(r.isResource), r.name, tfKey, tfConfig,
		)
		r.parseRegex(v, i)

		configs[i].validations = r.parseValidations(v, i)
	}

	return configs
}

// getTestCases populate TestSteps
func (r *reader) getTestCases(version string, getAPI GetAPIFunc) []resource.TestStep {
	v := r.getViperConfig(version)
	if v.GetBool("ignore") {
		r.t.Skip("ignoring tests for resource ", r.name)
	}

	configs := r.parseConfig(v)
	testSteps := make([]resource.TestStep, 0, len(configs))

	tag := ""
	if !r.isResource {
		tag = "data."
	}

	for _, c := range configs {
		testSteps = append(testSteps, resource.TestStep{
			Config: c.config,
			Check: resource.ComposeTestCheckFunc(
				validateResource(
					fmt.Sprintf("%s%s.tf_%s", tag, r.name, getLocalName(r.name)),
					c.validations,
					getAPI,
				),
			),
			ExpectError: r.expectError,
		})
	}

	return testSteps
}
