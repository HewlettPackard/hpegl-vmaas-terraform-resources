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

	createReq := models.CreateLBPool{
		CreateLBPoolReq: models.CreateLBPoolReq{
			Name:        d.GetString("name"),
			LbID:        d.GetInt("lb_id"),
			Description: d.GetString("description"),
			VipBalance:  d.GetString("vip_balance"),
			MinActive:   d.GetInt("min_active"),
			PoolConfig: models.PoolConfig{
				SnatTranslationType:   d.GetString("snat_translation_type"),
				PassiveMonitorPath:    d.GetInt("passive_monitor_path"),
				ActiveMonitorPaths:    d.GetInt("active_monitor_paths"),
				TCPMultiplexing:       d.GetBool("tcp_multiplexing"),
				TCPMultiplexingNumber: d.GetInt("tcp_multiplexing_number"),
				SnatIPAddress:         d.GetString("snat_ip_address"),
				MemberGroup: models.MemberGroup{
					Name:             d.GetString("name"),
					Path:             d.GetString("path"),
					IPRevisionFilter: d.GetString("ip_revision_filter"),
					Port:             d.GetInt("port"),
				},
			},
		},
	}

	createReq.CreateLBPoolReq.PoolConfig.SnatTranslationType = "LBSnatAutoMap"
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
		RetryDelay:   1,
		InitialDelay: 1,
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

	updateReq := models.CreateLBPool{
		CreateLBPoolReq: models.CreateLBPoolReq{
			Name:        d.GetString("name"),
			LbID:        d.GetInt("lb_id"),
			Description: d.GetString("description"),
			VipBalance:  d.GetString("vip_balance"),
			MinActive:   d.GetInt("min_active"),
			PoolConfig: models.PoolConfig{
				SnatTranslationType:   d.GetString("snat_translation_type"),
				PassiveMonitorPath:    d.GetInt("passive_monitor_path"),
				ActiveMonitorPaths:    d.GetInt("active_monitor_paths"),
				TCPMultiplexingNumber: d.GetInt("tcp_multiplexing_number"),
				SnatIPAddress:         d.GetString("snat_ip_address"),
				MemberGroup: models.MemberGroup{
					Name:             d.GetString("name"),
					Path:             d.GetString("path"),
					IPRevisionFilter: d.GetString("ip_revision_filter"),
					Port:             d.GetInt("port"),
				},
			},
		},
	}

	updateReq.CreateLBPoolReq.PoolConfig.SnatTranslationType = "LBSnatAutoMap"

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
