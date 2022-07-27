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
	var lbPoolResp models.GetSpecificLBPoolResp
	if err := tftags.Get(d, &lbPoolResp); err != nil {
		return err
	}

	lbDetails, err := lb.lbClient.GetLoadBalancers(ctx)
	if err != nil {
		return err
	}

	getlbPoolResp, err := lb.lbClient.GetSpecificLBPool(ctx, lbDetails.GetNetworkLoadBalancerResp[0].ID, lbPoolResp.ID)
	if err != nil {
		return err
	}

	return tftags.Set(d, getlbPoolResp.GetSpecificLBPoolResp)
}

func (lb *loadBalancerPool) Create(ctx context.Context, d *utils.Data, meta interface{}) error {
	setMeta(meta, lb.lbClient.Client)
	createReq := models.CreateLBPool{
		CreateLBPoolReq: models.CreateLBPoolReq{
			Name:        d.GetString("name"),
			Description: d.GetString("description"),
			VipBalance:  d.GetString("vip_balance"),
			MinActive:   d.GetInt("min_active"),
			// PoolConfig: models.PoolConfig{
			// 	SnatTranslationType:   d.GetString("snat_translation_type"),
			// 	PassiveMonitorPath:    d.GetInt("passive_monitor_path"),
			// 	ActiveMonitorPaths:    d.GetInt("active_monitor_paths"),
			// 	TCPMultiplexing:       d.GetBool("tcp_multiplexing"),
			// 	TCPMultiplexingNumber: d.GetInt("tcp_multiplexing_number"),
			// 	SnatIPAddress:         d.GetString("snat_ip_address"),
			// 	MemberGroup: models.MemberGroup{
			// 		Name:             d.GetString("name"),
			// 		Path:             d.GetString("path"),
			// 		IPRevisionFilter: d.GetString("ip_revision_filter"),
			// 		Port:             d.GetInt("port"),
			// 	},
			// },
		},
	}

	if err := tftags.Get(d, &createReq.CreateLBPoolReq); err != nil {
		return err
	}

	lbDetails, err := lb.lbClient.GetLoadBalancers(ctx)
	if err != nil {
		return err
	}

	createReq.CreateLBPoolReq.PoolConfig.SnatTranslationType = "LBSnatAutoMap"
	lbPoolResp, err := lb.lbClient.CreateLBPool(ctx, createReq, lbDetails.GetNetworkLoadBalancerResp[0].ID)
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
		return lb.lbClient.GetSpecificLBPool(ctx, lbDetails.GetNetworkLoadBalancerResp[0].ID,
			lbPoolResp.LBPoolResp.ID)
	})
	if err != nil {
		return err
	}

	return tftags.Set(d, createReq.CreateLBPoolReq)
}

func (lb *loadBalancerPool) Update(ctx context.Context, d *utils.Data, meta interface{}) error {
	return nil
}

func (lb *loadBalancerPool) Delete(ctx context.Context, d *utils.Data, meta interface{}) error {
	lbPoolID := d.GetID()
	lbDetails, err := lb.lbClient.GetLoadBalancers(ctx)
	if err != nil {
		return err
	}

	_, err = lb.lbClient.DeleteLBPool(ctx, lbDetails.GetNetworkLoadBalancerResp[0].ID, lbPoolID)
	if err != nil {
		return err
	}

	return nil
}
