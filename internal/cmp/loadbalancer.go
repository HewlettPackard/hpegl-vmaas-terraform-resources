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

type loadBalancer struct {
	lbClient *client.LoadBalancerAPIService
	rClient  *client.RouterAPIService
}

func newLoadBalancer(loadBalancerClient *client.LoadBalancerAPIService, routerClient *client.RouterAPIService) *loadBalancer {
	return &loadBalancer{
		lbClient: loadBalancerClient,
		rClient:  routerClient,
	}
}

func (lb *loadBalancer) Read(ctx context.Context, d *utils.Data, meta interface{}) error {
	var loadBalancerResp models.GetSpecificNetworkLoadBalancerResp
	if err := tftags.Get(d, &loadBalancerResp); err != nil {
		return err
	}
	getResLoadBalancer, err := lb.lbClient.GetSpecificLoadBalancers(ctx, loadBalancerResp.ID)
	if err != nil {
		return err
	}

	return tftags.Set(d, getResLoadBalancer.GetSpecificNetworkLoadBalancerResp)
}

func (lb *loadBalancer) Update(ctx context.Context, d *utils.Data, meta interface{}) error {
	id := d.GetID()
	var updateReq models.CreateLoadBalancerRequest
	if err := tftags.Get(d, &updateReq.NetworkLoadBalancer); err != nil {
		return err
	}

	retry := &utils.CustomRetry{
		InitialDelay: time.Second * 15,
		RetryDelay:   time.Second * 30,
	}
	_, err := retry.Retry(ctx, meta, func(ctx context.Context) (interface{}, error) {
		return lb.lbClient.UpdateLoadBalancer(ctx, id, updateReq)
	})
	if err != nil {
		return err
	}

	return tftags.Set(d, updateReq.NetworkLoadBalancer)
}

func (lb *loadBalancer) loadBalancerAlignRequest(ctx context.Context, meta interface{},
	createReq *models.CreateLoadBalancerRequest) error {

	allTypes, _ := lb.lbClient.GetLoadBalancerTypes(ctx, map[string]string{
		nameKey: nsxt,
	})

	// Get network service ID
	setMeta(meta, lb.rClient.Client)
	nsRetry := utils.CustomRetry{}
	nsRetry.RetryParallel(ctx, meta, func(ctx context.Context) (interface{}, error) {
		return lb.rClient.GetNetworkServices(ctx, nil)
	})

	// Align Network Server
	nsResp, err := nsRetry.Wait()
	if err != nil {
		return err
	}
	networkService := nsResp.(models.GetNetworkServicesResp)

	for i, n := range networkService.NetworkServices {
		if n.TypeName == nsxt {
			createReq.NetworkLoadBalancer.NetworkServerID = networkService.NetworkServices[i].ID
			break
		}
	}
	for i, n := range allTypes.LoadBalancerTypes {
		if n.Name == nsxt {
			createReq.NetworkLoadBalancer.Type = allTypes.LoadBalancerTypes[i].Name
			break
		}
	}
	return nil
}

func (lb *loadBalancer) Create(ctx context.Context, d *utils.Data, meta interface{}) error {
	setMeta(meta, lb.lbClient.Client)
	var createReq models.CreateLoadBalancerRequest
	if err := tftags.Get(d, &createReq.NetworkLoadBalancer); err != nil {
		return err
	}

	// align createReq and fill json related fields
	if err := lb.loadBalancerAlignRequest(ctx, meta, &createReq); err != nil {
		return err
	}

	lbResp, err := lb.lbClient.CreateLoadBalancer(ctx, createReq)
	if err != nil {
		return err
	}
	if !lbResp.Success {
		return fmt.Errorf(successErr, "creating LB")
	}
	createReq.NetworkLoadBalancer.ID = lbResp.NetworkLoadBalancerResp.ID

	// wait until created
	retry := &utils.CustomRetry{
		InitialDelay: time.Second * 15,
		RetryDelay:   time.Second * 30,
	}
	_, err = retry.Retry(ctx, meta, func(ctx context.Context) (interface{}, error) {
		return lb.lbClient.GetSpecificLoadBalancers(ctx, lbResp.NetworkLoadBalancerResp.ID)
	})
	if err != nil {
		return err
	}

	return tftags.Set(d, createReq.NetworkLoadBalancer)
}

func (lb *loadBalancer) Delete(ctx context.Context, d *utils.Data, meta interface{}) error {
	lbID := d.GetID()
	_, err := lb.lbClient.DeleteLoadBalancer(ctx, lbID)
	if err != nil {
		return err
	}

	return nil
}
