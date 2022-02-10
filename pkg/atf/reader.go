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
	vars        map[string]interface{}
}

func newReader(t *testing.T, isResource bool, name string) *reader {
	return &reader{
		t:          t,
		isResource: isResource,
		name:       name,
		vars:       make(map[string]interface{}),
	}
}

func (r *reader) fatalf(format string, v ...interface{}) {
	r.t.Fatalf("[acc-test] test case for resource "+r.name+" failed. "+format, v...)
}

func (r *reader) skipf(format string, v ...interface{}) {
	r.t.Skipf("[acc-test] test case for resource "+r.name+" is skipped. "+format, v...)
}

func (r *reader) readVars(v *viper.Viper) {
	vars, ok := v.Get("vars").(map[string]interface{})
	if !ok {
		return
	}

	for key, val := range vars {
		r.vars[key] = parseMeta(fmt.Sprint(val))
	}
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
		r.skipf("error while reading config, %v", err)
	}

	return v
}

func replaceVar(vars map[string]interface{}, config string) string {
	exp := `\$\([a-zA-Z_0-9]+\)`
	reg := regexp.MustCompile(exp)

	matches := reg.FindAllString(config, -1)
	for _, m := range matches {
		config = strings.Replace(config, m, fmt.Sprint(vars[m[2:len(m)-1]]), 1)
	}

	return config
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
		str := k.(string)
		split := strings.Split(str, ".")
		if len(split) > 1 && (split[0] == jsonKey || split[0] == tfKey) {
			isJSON := false
			if split[0] == jsonKey {
				isJSON = true
			}

			m = append(m, validation{
				isJSON: isJSON,
				key:    str[len(split[0])+1:],
				value:  v,
			})
		} else {
			r.fatalf("invalid validation format. validation format should be '[json|tf].key1.key2....keyn: value'")
		}
	}

	return m
}

func (r *reader) parseRegex(v *viper.Viper, i int) {
	if expectErrStr := v.GetString(path(accKey, i, "expect_error")); expectErrStr != "" {
		var err error
		r.expectError, err = regexp.Compile(expectErrStr)
		if err != nil {
			r.fatalf("error while compiling regex %s, got error %v", expectErrStr, err)
		}
	}
}

// parseConfig populates terraform configuration and parse to accConfig
func (r *reader) parseConfig(v *viper.Viper) []accConfig {
	tfKey := getLocalName(r.name)

	testCases := v.Get(accKey).([]interface{})
	configs := make([]accConfig, len(testCases))
	for i := range testCases {
		tfConfig := replaceVar(r.vars, v.GetString(path(accKey, i, "config")))
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

	r.readVars(v)
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
