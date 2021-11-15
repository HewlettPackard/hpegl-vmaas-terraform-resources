// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package acceptancetest

import (
	"os"
	"testing"

	testutils "github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/test-utils"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/constants"
	pkgutils "github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	libUtils "github.com/hewlettpackard/hpegl-provider-lib/pkg/utils"
)

var (
	testAccProviders map[string]*schema.Provider
	testAccProvider  *schema.Provider
)

func init() {
	testAccProvider = testutils.ProviderFunc()()
	testAccProviders = map[string]*schema.Provider{
		"hpegl": testAccProvider,
	}
}

func testAccPreCheck(t *testing.T) {
	// validate all required envs are present, if not then throws error
	if pkgutils.GetEnvBool(constants.MockIAMKey) {
		requiredenvs := []string{constants.CmpSubjectKey}
		for _, r := range requiredenvs {
			if os.Getenv(r) == "" {
				panic(r + " env is required, but not found")
			}
		}
	}

	t.Helper()
}

func TestProvider(t *testing.T) {
	if err := testutils.ProviderFunc()().InternalValidate(); err != nil {
		t.Fatalf("%s\n", err)
	}
	testAccPreCheck(t)
}

// TestMain is configure so that individual test cases can be ran
// nolint
func TestMain(m *testing.M) {
	// TF_ACC_CONFIG_PATH set in make acceptance
	libUtils.ReadAccConfig(os.Getenv("TF_ACC_CONFIG_PATH"))
	// Read skip
	pkgutils.ReadSkip()
	m.Run()
	// os.Exit(0)
}
