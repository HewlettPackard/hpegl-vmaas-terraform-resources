// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"
	"fmt"

	"github.com/hpe-hcss/vmaas-cmp-go-sdk/pkg/client"
	"github.com/hpe-hcss/vmaas-cmp-go-sdk/pkg/models"
	"github.com/hpe-hcss/vmaas-terraform-resources/internal/logger"
	"github.com/hpe-hcss/vmaas-terraform-resources/internal/utils"
)

type plan struct {
	pClient *client.PlansAPIService
}

func newPlan(pClient *client.PlansAPIService) *plan {
	return &plan{pClient: pClient}
}

func (n *plan) Read(ctx context.Context, d *utils.Data) error {
	logger.Debug("Get plan")

	name := d.GetString("name")
	// Pre check
	if err := d.Error(); err != nil {
		return err
	}
	resp, err := utils.Retry(func() (interface{}, error) {
		return n.pClient.GetAllServicePlans(ctx, map[string]string{
			provisionTypeKey: vmware,
			nameKey:          name,
		})
	})
	if err != nil {
		return err
	}
	plans := resp.(models.ServicePlans)
	if len(plans.ServicePlansResponse) != 1 {
		return fmt.Errorf(errExactMatch, "plan")
	}
	d.SetID(plans.ServicePlansResponse[0].ID)

	return d.Error()
}
