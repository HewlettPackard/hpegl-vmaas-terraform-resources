// Copyright (c) 2016-2020 Hewlett Packard Enterprise Development LP.

package acceptancetest

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	testutils "github.com/hpe-hcss/poc-caas-terraform-resources/internal/test-utils"
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
	t.Helper()
	// this fails c is a nil interface....
	// c := testAccProvider.Meta().(*Config)
	// if c.member.GetHosterID() == "" {
	// 	t.Fatalf("Acceptance tests must be run with hoster-scope %+v", c.member)
	// }
}

func TestProvider(t *testing.T) {
	if err := testutils.ProviderFunc()().InternalValidate(); err != nil {
		t.Fatalf("%s\n", err)
	}
	testAccPreCheck(t)
}
