// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package acceptancetest

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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
func TestAccResourceInstance(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
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
			  datastore_id = 6
			  root         = true
			}

			config {
			  resource_pool_id = 3
			  no_agent         = true
			  vm_folder        = "ComputeFolder"
			}
		}
	`, providerStanza, r.Int63n(999999), r.Intn(5)+5)
}
