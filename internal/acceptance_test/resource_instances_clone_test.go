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

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	api_client "github.com/hpe-hcss/vmaas-cmp-go-sdk/pkg/client"
	"github.com/hpe-hcss/vmaas-terraform-resources/internal/utils"
)

func TestVmaasInstanceClonePlan(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:             testAccResourceInstanceClone(),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccResourceInstanceCloneCreate(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping instance clone resource creation in short mode")
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testVmaasInstanceCloneDestroy("hpegl_vmaas_instance_clone.tf_acc_instance_clone"),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceInstanceClone(),
				Check: resource.ComposeTestCheckFunc(
					validateResource(
						"hpegl_vmaas_instance_clone.tf_acc_instance_clone",
						validateVmaasInstanceCloneStatus,
					),
				),
			},
		},
	})
}

func validateVmaasInstanceCloneStatus(rs *terraform.ResourceState) error {
	if rs.Primary.Attributes["status"] != "running" {
		return fmt.Errorf("expected %s but got %s", "running", rs.Primary.Attributes["status"])
	}

	return nil
}

func testVmaasInstanceCloneDestroy(name string) resource.TestCheckFunc {
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

		statusCode := utils.GetStatusCode(err)
		if statusCode != http.StatusNotFound {
			return fmt.Errorf("Expected %d statuscode, but got %d", 404, statusCode)
		}

		return nil
	}
}

func testAccResourceInstanceClone() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	return fmt.Sprintf(`%s
		resource "hpegl_vmaas_instance_clone" "tf_acc_instance_clone" {
			name               = "tf_acc_clone_%d"
			source_instance_id = 150
			network {
			  id = 6
			  interface_id = 9
			}
			network {
			  id = 6
			  interface_id = 4
			}
		}
	`, providerStanza, r.Int63n(999999))
}