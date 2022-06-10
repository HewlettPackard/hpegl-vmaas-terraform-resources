// (C) Copyright 2022 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"
	"fmt"

	"github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/client"
	"github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/models"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/utils"
	"github.com/tshihad/tftags"
)

type loadBalancerMonitor struct {
	lbClient *client.LoadBalancerAPIService
}

func newLoadBalancerMonitor(loadBalancerClient *client.LoadBalancerAPIService) *loadBalancerMonitor {
	return &loadBalancerMonitor{
		lbClient: loadBalancerClient,
	}
}

func (lb *loadBalancerMonitor) Read(ctx context.Context, d *utils.Data, meta interface{}) error {
	var lbMonitorResp models.GetSpecificLBMonitorResp
	if err := tftags.Get(d, &lbMonitorResp); err != nil {
		return err
	}
	getlbMonitorResp, err := lb.lbClient.GetSpecificLBMonitor(ctx, 1, lbMonitorResp.ID)
	if err != nil {
		return err
	}

	return tftags.Set(d, getlbMonitorResp.GetSpecificLBMonitorResp)
}

func (lb *loadBalancerMonitor) Create(ctx context.Context, d *utils.Data, meta interface{}) error {
	createReq := models.CreateLBMonitor{}
	if err := tftags.Get(d, &createReq.CreateLBMonitorReq); err != nil {
		return err
	}

	lbMonitorResp, err := lb.lbClient.CreateLBMonitor(ctx, createReq, 1)
	if err != nil {
		return err
	}
	if !lbMonitorResp.Success {
		return fmt.Errorf(successErr, "creating  loadBalancerMonitor Monitor")
	}

	// wait until created
	retry := &utils.CustomRetry{
		RetryDelay:   1,
		InitialDelay: 1,
	}
	_, err = retry.Retry(ctx, meta, func(ctx context.Context) (interface{}, error) {
		return lb.lbClient.GetSpecificLBMonitor(ctx, 1, lbMonitorResp.LBMonitorResp.ID)
	})
	if err != nil {
		return err
	}

	return tftags.Set(d, createReq.CreateLBMonitorReq)
}

func (lb *loadBalancerMonitor) Delete(ctx context.Context, d *utils.Data, meta interface{}) error {
	lbMonitorID := d.GetID()
	_, err := lb.lbClient.DeleteLBMonitor(ctx, 1, lbMonitorID)
	if err != nil {
		return err
	}

	return nil
}

func (lb *loadBalancerMonitor) Update(ctx context.Context, d *utils.Data, meta interface{}) error {
	return nil
}
