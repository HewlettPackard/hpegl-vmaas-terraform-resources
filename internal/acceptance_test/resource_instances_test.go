// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package acceptancetest

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	api_client "github.com/hpe-hcss/vmaas-cmp-go-sdk/pkg/client"
)

func TestVmaasInstancePlan(t *testing.T) {
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
	if testing.Short() {
		t.Skip("Skipping instance resource creation in short mode")
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(testVmaasInstanceDestroy("hpegl_vmaas_instance.tf_acc_instance")),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceInstance(),
				Check: resource.ComposeTestCheckFunc(
					validateResource(
						"hpegl_vmaas_instance.tf_acc_instance",
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
			return fmt.Errorf("error while converting id into int, %v", err)
		}

		apiClient, cfg := getApiClient()
		iClient := api_client.InstancesApiService{
			Client: apiClient,
			Cfg:    cfg,
		}
		_, err = iClient.GetASpecificInstance(context.Background(), id)

		// Once error is wrapped and send as json string we can encode error here and check the status code
		// once sdk-api support that functionality, update here
		if err == nil {
			return fmt.Errorf("Expected %d error, but got nil", 404)
		}

		return nil
	}
}

func testAccResourceInstance() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	return fmt.Sprintf(`%s
		resource "hpegl_vmaas_instance" "tf_acc_instance" {
			name               = "tf_acc_%d"
			cloud_id           = 1
			group_id           = 1
			layout_id          = 113
			plan_id            = 407
			instance_type_code = "vmware"
			network {
			  id = 6
			}

			volume {
			  name         = "root_vol"
			  size         = %d
			  datastore_id = 13
			  root         = true
			}

			config {
			  resource_pool_id = 3
			  no_agent         = true
			  vm_folder        = "group-v284"
			  template_id	   = 580
			}
		}
	`, providerStanza, r.Int63n(999999), r.Intn(5)+5)
}
