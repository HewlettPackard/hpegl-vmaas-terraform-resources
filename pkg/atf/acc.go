package atf

import (
	"fmt"
	"net/http"
	"testing"

	pkgutils "github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// GetAPIFunc accepts terraform states attribures as params and
// expects response and error as return values
type GetAPIFunc func(attr map[string]string) (interface{}, error)

type Acc struct {
	PreCheck     func(t *testing.T)
	Providers    map[string]*schema.Provider
	GetApi       GetAPIFunc
	ResourceName string
}

// RunResourcePlanTest to run resource plan only test case. This will take first
// config from specific resource
func (a *Acc) RunResourcePlanTest(t *testing.T) {
	a.runPlanTest(t, true)
}

// RunDataSourcePlanTest to run data source plan only test case. This will take first
// config from specific data source
func (a *Acc) RunDataSourcePlanTest(t *testing.T) {
	a.runPlanTest(t, false)
}

// RunTests creates test cases and run tests which includes create/update/delete/read
func (a *Acc) RunTests(t *testing.T) {
	// populate test cases
	testSteps := getTestCases(t, a.ResourceName, a.GetApi, true)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { a.PreCheck(t) },
		Providers: a.Providers,
		CheckDestroy: resource.ComposeTestCheckFunc(
			a.checkResourceDestroy,
		),
		Steps: testSteps,
	})
}

// checkResourceDestroy checks resource destroy conditions. This function will parse error
// and check status code is 404 or not
func (a *Acc) checkResourceDestroy(s *terraform.State) error {
	rs, ok := s.RootModule().Resources[fmt.Sprintf("%s.tf_%s", a.ResourceName, getLocalName(a.ResourceName))]
	if !ok {
		return fmt.Errorf("[Check Destroy] resource %s not found", a.ResourceName)
	}
	_, err := a.GetApi(rs.Primary.Attributes)
	statusCode := pkgutils.GetStatusCode(err)
	if statusCode != http.StatusNotFound {
		return fmt.Errorf("expected %d statuscode, but got %d", 404, statusCode)
	}

	return nil
}

// runs plan test for resource or data source. only first config from test case
// will considered on plan test
func (a *Acc) runPlanTest(t *testing.T, isResource bool) {
	pkgutils.SkipAcc(t, fmt.Sprintf("vmaas.%s.%s", getTag(isResource), getLocalName(a.ResourceName)))

	testSteps := getTestCases(t, a.ResourceName, a.GetApi, true)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { a.PreCheck(t) },
		Providers: a.Providers,
		Steps: []resource.TestStep{
			{
				Config:             testSteps[0].Config,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
				Check:              testSteps[0].Check,
			},
		},
	})
}
