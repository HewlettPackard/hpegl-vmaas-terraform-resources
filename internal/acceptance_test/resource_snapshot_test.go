// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package acceptancetest

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestVmaasSnapshotPlan(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:             testAccResourceSnapshot(),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccResourceSnapshot() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	return fmt.Sprintf(`%s
		resource "hpegl_vmaas_snapshot" "tf_acc_instance" {
			name        = "tf_acc_snapshot_%d"
			instance_id = 30
		}
	`, providerStanza, r.Int63n(999999))
}
