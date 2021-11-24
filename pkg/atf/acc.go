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

type Acc struct {
	PreCheck     func(t *testing.T)
	Providers    map[string]*schema.Provider
	GetApi       func(attr map[string]string) (interface{}, error)
	ResourceName string
}

func (a *Acc) RunResourcePlanTest(t *testing.T) {
	a.runPlanTest(t, true)
}

func (a *Acc) RunDataSourcePlanTest(t *testing.T) {
	a.runPlanTest(t, false)
}
func (a *Acc) RunCreateTests(t *testing.T) {
	testSteps := getTestCases(t, a.ResourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { a.PreCheck(t) },
		Providers: a.Providers,
		CheckDestroy: resource.ComposeTestCheckFunc(
			a.checkResourceDestroy,
		),
		Steps: testSteps,
	})
}

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

func (a *Acc) runPlanTest(t *testing.T, isResource bool) {
	pkgutils.SkipAcc(t, "vmaas.resource."+getLocalName(a.ResourceName))
	if testing.Short() {
		t.Skip("Skipping ", a.ResourceName, " resource creation in short mode")
	}

	tag := "resource"
	if !isResource {
		tag = "datasource"
	}
	pkgutils.SkipAcc(t, fmt.Sprintf("vmaas.%s.%s", tag, getLocalName(a.ResourceName)))
	configs := getResourceConfig(a.ResourceName)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { a.PreCheck(t) },
		Providers: a.Providers,
		Steps: []resource.TestStep{
			{
				Config:             configs[0].Config,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
