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

	getMonitorLoadBalancer, err := lb.lbClient.GetSpecificLBMonitor(ctx, lbMonitorResp.LbID, lbMonitorResp.ID)
	if err != nil {
		return err
	}
	return tftags.Set(d, getMonitorLoadBalancer.GetSpecificLBMonitorResp)

}

func (lb *loadBalancerMonitor) Create(ctx context.Context, d *utils.Data, meta interface{}) error {
	setMeta(meta, lb.lbClient.Client)

	// var createReq models.CreateLBMonitor
	// if err := tftags.Get(d, &createReq.CreateLBMonitorReq); err != nil {
	// 	return err
	// }

	createReq := models.CreateLBMonitor{
		CreateLBMonitorReq: models.CreateLBMonitorReq{
			Name:                d.GetString("name"),
			LbID:                d.GetInt("lb_id"),
			Description:         d.GetString("description"),
			Type:                d.GetString("type"),
			Timeout:             d.GetInt("timeout"),
			Interval:            d.GetInt("interval"),
			RequestVersion:      d.GetString("request_version"),
			RequestMethod:       d.GetString("request_method"),
			ResponseStatusCodes: d.GetString("response_status_codes"),
			MaxFail:             d.GetInt("max_fail"),
			ResponseData:        d.GetString("response_data"),
			RequestURL:          d.GetString("request_url"),
			RequestBody:         d.GetString("request_body"),
			AliasPort:           d.GetInt("alias_port"),
			RiseCount:           d.GetInt("rise_count"),
			FallCount:           d.GetInt("fall_count"),
			DataLength:          d.GetInt("data_length"),
		},
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

	// updateReq := models.CreateLBMonitor{
	// 	CreateLBMonitorReq: models.CreateLBMonitorReq{
	// 		Name:               d.GetString("name"),
	// 		LbID:               d.GetInt("lb_id"),
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
