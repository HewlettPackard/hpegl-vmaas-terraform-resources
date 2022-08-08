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

type loadBalancerPool struct {
	lbClient *client.LoadBalancerAPIService
}

func newLoadBalancerPool(loadBalancerClient *client.LoadBalancerAPIService) *loadBalancerPool {
	return &loadBalancerPool{
		lbClient: loadBalancerClient,
	}
}

func (lb *loadBalancerPool) Read(ctx context.Context, d *utils.Data, meta interface{}) error {
	setMeta(meta, lb.lbClient.Client)
	var lbPoolResp models.CreateLBPoolReq
	if err := tftags.Get(d, &lbPoolResp); err != nil {
		return err
	}

	getPoolLoadBalancer, err := lb.lbClient.GetSpecificLBPool(ctx, lbPoolResp.LbID, lbPoolResp.ID)
	if err != nil {
		return err
	}

	return tftags.Set(d, getPoolLoadBalancer.GetSpecificLBPoolResp)
}

func (lb *loadBalancerPool) Create(ctx context.Context, d *utils.Data, meta interface{}) error {
	setMeta(meta, lb.lbClient.Client)

	var createReq models.CreateLBPool
	if err := tftags.Get(d, &createReq.CreateLBPoolReq); err != nil {
		return err
	}

	lbPoolResp, err := lb.lbClient.CreateLBPool(ctx, createReq, createReq.CreateLBPoolReq.LbID)
	if err != nil {
		return err
	}
	if !lbPoolResp.Success {
		return fmt.Errorf(successErr, "creating loadBalancer Pool")
	}
	createReq.CreateLBPoolReq.ID = lbPoolResp.LBPoolResp.ID
	// wait until created
	retry := &utils.CustomRetry{
		InitialDelay: time.Second * 15,
		RetryDelay:   time.Second * 30,
	}
	_, err = retry.Retry(ctx, meta, func(ctx context.Context) (interface{}, error) {
		return lb.lbClient.GetSpecificLBPool(ctx, createReq.CreateLBPoolReq.LbID,
			lbPoolResp.LBPoolResp.ID)
	})
	if err != nil {
		return err
	}

	return tftags.Set(d, createReq.CreateLBPoolReq)
}

func (lb *loadBalancerPool) Update(ctx context.Context, d *utils.Data, meta interface{}) error {
	id := d.GetID()

	var updateReq models.CreateLBPool
	if err := tftags.Get(d, &updateReq.CreateLBPoolReq); err != nil {
		return err
	}

	retry := &utils.CustomRetry{
		InitialDelay: time.Second * 15,
		RetryDelay:   time.Second * 30,
	}
	_, err := retry.Retry(ctx, meta, func(ctx context.Context) (interface{}, error) {
		return lb.lbClient.UpdateLBPool(ctx, updateReq,
			updateReq.CreateLBPoolReq.LbID, id)
	})
	if err != nil {
		return err
	}

	return tftags.Set(d, updateReq.CreateLBPoolReq)
}

func (lb *loadBalancerPool) Delete(ctx context.Context, d *utils.Data, meta interface{}) error {
	setMeta(meta, lb.lbClient.Client)
	var tfLBPool models.CreateLBPoolReq
	if err := tftags.Get(d, &tfLBPool); err != nil {
		return err
	}

	resp, err := lb.lbClient.DeleteLBPool(ctx, tfLBPool.LbID, tfLBPool.ID)
	if err != nil {
		return err
	}

	if !resp.Success {
		return fmt.Errorf("got success = 'false' while deleting LB-POOL")
	}

	return nil
}
