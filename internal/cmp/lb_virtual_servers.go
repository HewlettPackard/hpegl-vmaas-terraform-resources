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

type loadBalancerVirtualServer struct {
	lbClient *client.LoadBalancerAPIService
}

func newLoadBalancerVirtualServer(loadBalancerClient *client.LoadBalancerAPIService) *loadBalancerVirtualServer {
	return &loadBalancerVirtualServer{
		lbClient: loadBalancerClient,
	}
}

func (lb *loadBalancerVirtualServer) Read(ctx context.Context, d *utils.Data, meta interface{}) error {
	var lbVirtualServerResp models.GetSpecificLBVirtualServersResp
	if err := tftags.Get(d, &lbVirtualServerResp); err != nil {
		return err
	}
	getlbVirtualServerResp, err := lb.lbClient.GetSpecificLBVirtualServer(ctx, 1, lbVirtualServerResp.ID)
	if err != nil {
		return err
	}

	return tftags.Set(d, getlbVirtualServerResp.GetSpecificLBVirtualServersResp)
}

func (lb *loadBalancerVirtualServer) Update(ctx context.Context, d *utils.Data, meta interface{}) error {
	return nil
}

func (lb *loadBalancerVirtualServer) Create(ctx context.Context, d *utils.Data, meta interface{}) error {
	createReq := models.CreateLBVirtualServers{}
	if err := tftags.Get(d, &createReq.CreateLBVirtualServersReq); err != nil {
		return err
	}

	lbVirtualServersResp, err := lb.lbClient.CreateLBVirtualServers(ctx, createReq, 1)
	if err != nil {
		return err
	}
	if !lbVirtualServersResp.Success {
		return fmt.Errorf(successErr, "creating loadBalancerVirtualServer Virtual Servers")
	}

	// wait until created
	retry := &utils.CustomRetry{
		RetryDelay:   1,
		InitialDelay: 1,
	}
	_, err = retry.Retry(ctx, meta, func(ctx context.Context) (interface{}, error) {
		return lb.lbClient.GetSpecificLBVirtualServer(ctx, 1, lbVirtualServersResp.CreateLBVirtualServersResp.ID)
	})
	if err != nil {
		return err
	}

	return tftags.Set(d, createReq.CreateLBVirtualServersReq)
}

func (lb *loadBalancerVirtualServer) Delete(ctx context.Context, d *utils.Data, meta interface{}) error {
	lbVirtualServerID := d.GetID()
	_, err := lb.lbClient.DeleteLBVirtualServers(ctx, 1, lbVirtualServerID)
	if err != nil {
		return err
	}

	return nil
}
