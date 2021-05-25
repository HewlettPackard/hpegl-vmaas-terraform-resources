// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"

	"github.com/hpe-hcss/vmaas-cmp-go-sdk/pkg/client"
	"github.com/hpe-hcss/vmaas-terraform-resources/internal/logger"
	"github.com/hpe-hcss/vmaas-terraform-resources/internal/utils"
)

type plan struct {
	pClient           *client.PlansApiService
	serviceInstanceID string
}

func newPlan(pClient *client.PlansApiService, serviceInstanceID string) *plan {
	return &plan{pClient: pClient, serviceInstanceID: serviceInstanceID}
}

func (n *plan) Read(ctx context.Context, d *utils.Data) error {
	logger.Debug("Get plan")

	// name := d.GetString("name")
	_, err := n.pClient.GetAllServicePlans(ctx, n.serviceInstanceID)
	// plans, err := n.pClient.GetAllServicePlans(ctx, n.serviceInstanceID, map[string]string{"name": name})
	if err != nil {
		return err
	}
	// if len(plans) != 1 {
	// 	return errors.New("Coudn't find exact plan, please check the name")
	// }
	// d.SetID(strconv.Itoa(int(plans.plans[0].Id)))
	// if d.HaveError() {
	// 	return err
	// }

	return nil
}
