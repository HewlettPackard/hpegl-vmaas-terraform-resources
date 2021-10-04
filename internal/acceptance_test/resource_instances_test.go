// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package acceptancetest

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"testing"
	"time"

	api_client "github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/client"
	pkgutils "github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/spf13/viper"
)

func TestVmaasInstancePlan(t *testing.T) {
	pkgutils.SkipAcc(t, "vmaas.resource.instance")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:             testAccResourceInstance(),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccResourceInstanceCreate(t *testing.T) {
	pkgutils.SkipAcc(t, "vmaas.resource.instance")
	if testing.Short() {
		t.Skip("Skipping instance resource creation in short mode")
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(testVmaasInstanceDestroy("hpegl_vmaas_instance.tf_instance")),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceInstance(),
				Check: resource.ComposeTestCheckFunc(
					validateResource(
						"hpegl_vmaas_instance.tf_instance",
						validateVmaasInstanceStatus,
					),
				),
			},
		},
	})
}

func validateVmaasInstanceStatus(rs *terraform.ResourceState) error {
	if rs.Primary.Attributes["status"] != "running" {
		return fmt.Errorf("expected %s but got %s", "running", rs.Primary.Attributes["status"])
	}

	return nil
}

func testVmaasInstanceDestroy(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("resource %s not found", name)
		}
		id, err := strconv.Atoi(rs.Primary.Attributes["id"])
		if err != nil {
			return fmt.Errorf("error while converting id into int, %w", err)
		}

		apiClient, cfg := getAPIClient()
		iClient := api_client.InstancesAPIService{
			Client: apiClient,
			Cfg:    cfg,
		}
		_, err = iClient.GetASpecificInstance(context.Background(), id)

		statusCode := pkgutils.GetStatusCode(err)
		if statusCode != http.StatusNotFound {
			return fmt.Errorf("Expected %d statuscode, but got %d", 404, statusCode)
		}

		return nil
	}
}

func testAccResourceInstance() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	instancePrimitiveStanza := fmt.Sprintf(`
			name               = "%s_%d"
			cloud_id           = %d
			group_id           = %d
			layout_id          = %d
			plan_id            = %d
			instance_type_code = "%s"
			scale = %d
`,
		viper.GetString("vmaas.resource.instance.name"),
		r.Int63n(999999),
		viper.GetInt("vmaas.resource.instance.cloud_id"),
		viper.GetInt("vmaas.resource.instance.group_id"),
		viper.GetInt("vmaas.resource.instance.layout_id"),
		viper.GetInt("vmaas.resource.instance.plan_id"),
		viper.GetString("vmaas.resource.instance.instance_type_code"),
		viper.GetInt("vmaas.resource.instance.scale"))

	configStanza := fmt.Sprintf(`
			config {
				resource_pool_id = %d
				no_agent         = %t
				template_id	   = %d
				}
			`,
		viper.GetInt("vmaas.resource.instance.config.resource_pool_id"),
		viper.GetBool("vmaas.resource.instance.config.no_agent"),
		viper.GetInt("vmaas.resource.instance.config.template_id"))

	return providerStanza + fmt.Sprintf(`
		resource "hpegl_vmaas_instance" "%s" {
			%s
			%s
			%s
			%s
		}
	`, "tf_instance",
		instancePrimitiveStanza,
		getNetworkStanza(),
		getVolumeStanza(),
		configStanza,
	)
}
