// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"
	"fmt"
	"strconv"

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

	name := d.GetString("name")
	plans, err := n.pClient.GetAllServicePlans(ctx, n.serviceInstanceID, map[string]string{
		provisionTypeKey: vmware,
		nameKey:          name,
	})
	if err != nil {
		return err
	}
	if len(plans.ServicePlansResponse) != 1 {
		return fmt.Errorf(errExactMatch, "plan")
	}
	d.SetID(strconv.Itoa(plans.ServicePlansResponse[0].ID))

	return d.Error()
}
