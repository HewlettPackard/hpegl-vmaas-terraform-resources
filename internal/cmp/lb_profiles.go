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

type loadBalancer struct {
	lbClient *client.LoadBalancerAPIService
}

func newLoadBalancer(loadBalancerClient *client.LoadBalancerAPIService) *loadBalancer {
	return &loadBalancer{
		lbClient: loadBalancerClient,
	}
}

func (lb *loadBalancer) Read(ctx context.Context, d *utils.Data, meta interface{}) error {
	var lbProfileResp models.GetLBSpecificProfilesResp
	if err := tftags.Get(d, &lbProfileResp); err != nil {
		return err
	}
	getlbProfileResp, err := lb.lbClient.GetSpecificLBProfile(ctx, lbProfileResp.ID)
	if err != nil {
		return err
	}

	return tftags.Set(d, getlbProfileResp.GetLBSpecificProfilesResp)
}

func (lb *loadBalancer) Create(ctx context.Context, d *utils.Data, meta interface{}) error {
	createReq := models.CreateLBProfile{}
	if err := tftags.Get(d, &createReq.CreateLBProfileReq); err != nil {
		return err
	}

	lbProfileResp, err := lb.lbClient.CreateLBProfile(ctx, createReq)
	if err != nil {
		return err
	}
	if !lbProfileResp.Success {
		return fmt.Errorf(successErr, "creating loadBalancer Profile")
	}

	// wait until created
	retry := &utils.CustomRetry{
		RetryDelay:   1,
		InitialDelay: 1,
	}
	_, err = retry.Retry(ctx, meta, func(ctx context.Context) (interface{}, error) {
		return lb.lbClient.GetSpecificLBProfile(ctx, lbProfileResp.LBProfileResp.ID)
	})
	if err != nil {
		return err
	}

	return tftags.Set(d, createReq.CreateLBProfileReq)
}

func (lb *loadBalancer) Delete(ctx context.Context, d *utils.Data, meta interface{}) error {
	lbProfileID := d.GetID()
	_, err := lb.lbClient.DeleteLBProfile(ctx, lbProfileID)
	if err != nil {
		return err
	}

	return nil
}
