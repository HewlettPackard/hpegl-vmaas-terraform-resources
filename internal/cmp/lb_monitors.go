// (C) Copyright 2022 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"
	"fmt"
	"time"

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
	setMeta(meta, lb.lbClient.Client)
	var lbMonitorResp models.CreateLBMonitorReq
	if err := tftags.Get(d, &lbMonitorResp); err != nil {
		return err
	}

	getMonitorLoadBalancer, err := lb.lbClient.GetSpecificLBMonitor(ctx, lbMonitorResp.LbID,
		lbMonitorResp.ID)
	if err != nil {
		return err
	}
	return tftags.Set(d, getMonitorLoadBalancer.GetSpecificLBMonitorResp)

}

func (lb *loadBalancerMonitor) Create(ctx context.Context, d *utils.Data, meta interface{}) error {
	setMeta(meta, lb.lbClient.Client)

	var createReq models.CreateLBMonitor
	if err := tftags.Get(d, &createReq.CreateLBMonitorReq); err != nil {
		return err
	}

	lbMonitorResp, err := lb.lbClient.CreateLBMonitor(ctx, createReq,
		createReq.CreateLBMonitorReq.LbID)
	if err != nil {
		return err
	}
	if !lbMonitorResp.Success {
		return fmt.Errorf(successErr, "creating loadBalancerMonitor Monitor")
	}

	createReq.CreateLBMonitorReq.ID = lbMonitorResp.LBMonitorResp.ID

	// wait until created
	retry := &utils.CustomRetry{
		InitialDelay: time.Second * 15,
		RetryDelay:   time.Second * 30,
	}
	_, err = retry.Retry(ctx, meta, func(ctx context.Context) (interface{}, error) {
		return lb.lbClient.GetSpecificLBMonitor(ctx, createReq.CreateLBMonitorReq.LbID,
			lbMonitorResp.LBMonitorResp.ID)
	})
	if err != nil {
		return err
	}

	return tftags.Set(d, createReq.CreateLBMonitorReq)
}

func (lb *loadBalancerMonitor) Delete(ctx context.Context, d *utils.Data, meta interface{}) error {
	setMeta(meta, lb.lbClient.Client)
	var tfLBMonitor models.CreateLBMonitorReq
	if err := tftags.Get(d, &tfLBMonitor); err != nil {
		return err
	}

	resp, err := lb.lbClient.DeleteLBMonitor(ctx, tfLBMonitor.LbID, tfLBMonitor.ID)
	if err != nil {
		return err
	}

	if !resp.Success {
		return fmt.Errorf("got success = 'false' while deleting LB-MONITOR")
	}

	return nil
}

func (lb *loadBalancerMonitor) Update(ctx context.Context, d *utils.Data, meta interface{}) error {
	id := d.GetID()

	var updateReq models.CreateLBMonitor
	if err := tftags.Get(d, &updateReq.CreateLBMonitorReq); err != nil {
		return err
	}

	retry := &utils.CustomRetry{
		InitialDelay: time.Second * 15,
		RetryDelay:   time.Second * 30,
	}
	_, err := retry.Retry(ctx, meta, func(ctx context.Context) (interface{}, error) {
		return lb.lbClient.UpdateLBMonitor(ctx, updateReq,
			updateReq.CreateLBMonitorReq.LbID, id)
	})
	if err != nil {
		return err
	}

	return tftags.Set(d, updateReq.CreateLBMonitorReq)
}
