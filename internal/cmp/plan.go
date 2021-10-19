// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"
	"fmt"
	"log"

	"github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/client"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/utils"
)

type plan struct {
	pClient *client.PlansAPIService
}

func newPlan(pClient *client.PlansAPIService) *plan {
	return &plan{pClient: pClient}
}

func (n *plan) Read(ctx context.Context, d *utils.Data, meta interface{}) error {
	log.Printf("[DEBUG] Get plan")

	name := d.GetString("name")
	// Pre check
	if err := d.Error(); err != nil {
		return err
	}
	plans, err := n.pClient.GetAllServicePlans(ctx, map[string]string{
		provisionTypeKey: vmware,
		nameKey:          name,
	})
	if err != nil {
		return err
	}
	if len(plans.ServicePlansResponse) != 1 {
		return fmt.Errorf(errExactMatch, "plan")
	}
	d.SetID(plans.ServicePlansResponse[0].ID)

	return d.Error()
}
