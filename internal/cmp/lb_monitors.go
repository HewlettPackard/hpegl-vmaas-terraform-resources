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

	_, err := lb.lbClient.GetSpecificLBMonitor(ctx, lbMonitorResp.LbID, lbMonitorResp.ID)
	if err != nil {
		return err
	}
	return tftags.Set(d, lbMonitorResp)
}

func (lb *loadBalancerMonitor) Create(ctx context.Context, d *utils.Data, meta interface{}) error {
	setMeta(meta, lb.lbClient.Client)
	var createReq models.CreateLBMonitor
	err := tftags.Get(d, &createReq.CreateLBMonitorReq)
	if err != nil {
		return err
	}

	// createReq := models.CreateLBMonitor{
	// 	CreateLBMonitorReq: models.CreateLBMonitorReq{
	// 		Name:               d.GetString("name"),
	// 		Description:        d.GetString("description"),
	// 		MonitorType:        d.GetString("monitor_type"),
	// 		MonitorTimeout:     d.GetInt("monitor_timeout"),
	// 		MonitorInterval:    d.GetInt("monitor_interval"),
	// 		SendVersion:        d.GetString("send_version"),
	// 		SendType:           d.GetString("send_type"),
	// 		MonitorDestination: d.GetString("monitor_destination"),
	// 		MonitorReverse:     d.GetBool("monitor_reverse"),
	// 		MonitorTransparent: d.GetBool("monitor_transparent"),
	// 		MonitorAdaptive:    d.GetBool("monitor_adaptive"),
	// 		FallCount:          d.GetInt("fall_count"),
	// 		RiseCount:          d.GetInt("rise_count"),
	// 		AliasPort:          d.GetInt("alias_port"),
	// 	},
	// }
	if err := tftags.Get(d, &createReq.CreateLBMonitorReq); err != nil {
		return err
	}

	lbMonitorResp, err := lb.lbClient.CreateLBMonitor(ctx, createReq, createReq.CreateLBMonitorReq.LbID)
	if err != nil {
		return err
	}
	if !lbMonitorResp.Success {
		return fmt.Errorf(successErr, "creating loadBalancerMonitor Monitor")
	}

	createReq.CreateLBMonitorReq.ID = lbMonitorResp.LBMonitorResp.ID

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
	id := d.GetID()

	updateReq := models.CreateLBMonitor{
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

	if err := d.Error(); err != nil {
		return err
	}

	lbDetails, err := lb.lbClient.GetLoadBalancers(ctx)
	if err != nil {
		return err
	}
	retry := &utils.CustomRetry{
		InitialDelay: time.Second * 15,
		RetryDelay:   time.Second * 30,
	}
	_, err = retry.Retry(ctx, meta, func(ctx context.Context) (interface{}, error) {
		return lb.lbClient.UpdateLBMonitor(ctx, updateReq, lbDetails.GetNetworkLoadBalancerResp[0].ID, id)
	})
	if err != nil {
		return err
	}

	//return d.Error()

	return nil
}
