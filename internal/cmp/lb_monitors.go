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
	var lbMonitorResp models.GetSpecificLBMonitorResp
	if err := tftags.Get(d, &lbMonitorResp); err != nil {
		return err
	}

	lbDetails, err := lb.lbClient.GetLoadBalancers(ctx)
	if err != nil {
		return err
	}

	getlbMonitorResp, err := lb.lbClient.GetSpecificLBMonitor(ctx, lbDetails.GetNetworkLoadBalancerResp[0].ID, lbMonitorResp.ID)
	if err != nil {
		return err
	}

	return tftags.Set(d, getlbMonitorResp.GetSpecificLBMonitorResp)
}

func (lb *loadBalancerMonitor) Create(ctx context.Context, d *utils.Data, meta interface{}) error {

	setMeta(meta, lb.lbClient.Client)

	createReq := models.CreateLBMonitor{
		CreateLBMonitorReq: models.CreateLBMonitorReq{
			Name:               d.GetString("name"),
			Description:        d.GetString("description"),
			MonitorType:        d.GetString("monitor_type"),
			MonitorTimeout:     d.GetInt("monitor_timeout"),
			MonitorInterval:    d.GetInt("monitor_interval"),
			SendVersion:        d.GetString("send_version"),
			SendType:           d.GetString("send_type"),
			MonitorDestination: d.GetString("monitor_destination"),
			MonitorReverse:     d.GetBool("monitor_reverse"),
			MonitorTransparent: d.GetBool("monitor_transparent"),
			MonitorAdaptive:    d.GetBool("monitor_adaptive"),
			FallCount:          d.GetInt("fall_count"),
			RiseCount:          d.GetInt("rise_count"),
			AliasPort:          d.GetInt("alias_port"),
		},
	}
	if err := tftags.Get(d, &createReq.CreateLBMonitorReq); err != nil {
		return err
	}

	lbDetails, err := lb.lbClient.GetLoadBalancers(ctx)
	if err != nil {
		return err
	}

	lbMonitorResp, err := lb.lbClient.CreateLBMonitor(ctx, createReq, lbDetails.GetNetworkLoadBalancerResp[0].ID)
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
		return lb.lbClient.GetSpecificLBMonitor(ctx, lbDetails.GetNetworkLoadBalancerResp[0].ID, lbMonitorResp.LBMonitorResp.ID)
	})
	if err != nil {
		return err
	}

	return tftags.Set(d, createReq.CreateLBMonitorReq)
}

func (lb *loadBalancerMonitor) Delete(ctx context.Context, d *utils.Data, meta interface{}) error {
	lbMonitorID := d.GetID()
	lbDetails, err := lb.lbClient.GetLoadBalancers(ctx)
	if err != nil {
		return err
	}
	_, err = lb.lbClient.DeleteLBMonitor(ctx, lbDetails.GetNetworkLoadBalancerResp[0].ID, lbMonitorID)
	if err != nil {
		return err
	}

	return nil
}

func (lb *loadBalancerMonitor) Update(ctx context.Context, d *utils.Data, meta interface{}) error {
	return nil
}
