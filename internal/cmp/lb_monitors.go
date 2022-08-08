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
	createReq := models.CreateLBMonitor{}
	if err := tftags.Get(d, &createReq.CreateLBMonitorReq); err != nil {
		return err
	}
	// align createReq and fill json related fields
	if err := lb.monitorAlignMonitorTypeRequest(&createReq); err != nil {
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

	// align createReq and fill json related fields
	if err := lb.monitorAlignMonitorTypeRequest(&updateReq); err != nil {
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

func (lb *loadBalancerMonitor) monitorAlignMonitorTypeRequest(monitorReq *models.CreateLBMonitor) error {
	if monitorReq.CreateLBMonitorReq.TfHTTPConfig != nil {
		monitorReq.CreateLBMonitorReq.RequestBody = monitorReq.CreateLBMonitorReq.TfHTTPConfig.RequestBody
		monitorReq.CreateLBMonitorReq.AliasPort = monitorReq.CreateLBMonitorReq.TfHTTPConfig.AliasPort
		monitorReq.CreateLBMonitorReq.FallCount = monitorReq.CreateLBMonitorReq.TfHTTPConfig.FallCount
		monitorReq.CreateLBMonitorReq.Interval = monitorReq.CreateLBMonitorReq.TfHTTPConfig.Interval
		monitorReq.CreateLBMonitorReq.RequestMethod = monitorReq.CreateLBMonitorReq.TfHTTPConfig.RequestMethod
		monitorReq.CreateLBMonitorReq.RequestURL = monitorReq.CreateLBMonitorReq.TfHTTPConfig.RequestURL
		monitorReq.CreateLBMonitorReq.RequestVersion = monitorReq.CreateLBMonitorReq.TfHTTPConfig.RequestVersion
		monitorReq.CreateLBMonitorReq.ResponseData = monitorReq.CreateLBMonitorReq.TfHTTPConfig.ResponseData
		monitorReq.CreateLBMonitorReq.ResponseStatusCodes = monitorReq.CreateLBMonitorReq.TfHTTPConfig.ResponseStatusCodes
		monitorReq.CreateLBMonitorReq.RiseCount = monitorReq.CreateLBMonitorReq.TfHTTPConfig.RiseCount
		monitorReq.CreateLBMonitorReq.Timeout = monitorReq.CreateLBMonitorReq.TfHTTPConfig.Timeout
	} else if monitorReq.CreateLBMonitorReq.TfHTTPSConfig != nil {
		monitorReq.CreateLBMonitorReq.RequestBody = monitorReq.CreateLBMonitorReq.TfHTTPSConfig.RequestBody
		monitorReq.CreateLBMonitorReq.AliasPort = monitorReq.CreateLBMonitorReq.TfHTTPSConfig.AliasPort
		monitorReq.CreateLBMonitorReq.FallCount = monitorReq.CreateLBMonitorReq.TfHTTPSConfig.FallCount
		monitorReq.CreateLBMonitorReq.Interval = monitorReq.CreateLBMonitorReq.TfHTTPSConfig.Interval
		monitorReq.CreateLBMonitorReq.RequestMethod = monitorReq.CreateLBMonitorReq.TfHTTPSConfig.RequestMethod
		monitorReq.CreateLBMonitorReq.RequestURL = monitorReq.CreateLBMonitorReq.TfHTTPSConfig.RequestURL
		monitorReq.CreateLBMonitorReq.RequestVersion = monitorReq.CreateLBMonitorReq.TfHTTPSConfig.RequestVersion
		monitorReq.CreateLBMonitorReq.ResponseData = monitorReq.CreateLBMonitorReq.TfHTTPSConfig.ResponseData
		monitorReq.CreateLBMonitorReq.ResponseStatusCodes = monitorReq.CreateLBMonitorReq.TfHTTPSConfig.ResponseStatusCodes
		monitorReq.CreateLBMonitorReq.RiseCount = monitorReq.CreateLBMonitorReq.TfHTTPSConfig.RiseCount
		monitorReq.CreateLBMonitorReq.Timeout = monitorReq.CreateLBMonitorReq.TfHTTPSConfig.Timeout
	} else if monitorReq.CreateLBMonitorReq.TfIcmpConfig != nil {
		monitorReq.CreateLBMonitorReq.AliasPort = monitorReq.CreateLBMonitorReq.TfIcmpConfig.AliasPort
		monitorReq.CreateLBMonitorReq.DataLength = monitorReq.CreateLBMonitorReq.TfIcmpConfig.DataLength
		monitorReq.CreateLBMonitorReq.FallCount = monitorReq.CreateLBMonitorReq.TfIcmpConfig.FallCount
		monitorReq.CreateLBMonitorReq.Interval = monitorReq.CreateLBMonitorReq.TfIcmpConfig.Interval
		monitorReq.CreateLBMonitorReq.RiseCount = monitorReq.CreateLBMonitorReq.TfIcmpConfig.RiseCount
		monitorReq.CreateLBMonitorReq.Timeout = monitorReq.CreateLBMonitorReq.TfIcmpConfig.Timeout
	} else if monitorReq.CreateLBMonitorReq.TfPassiveConfig != nil {
		monitorReq.CreateLBMonitorReq.MaxFail = monitorReq.CreateLBMonitorReq.TfPassiveConfig.MaxFail
		monitorReq.CreateLBMonitorReq.Timeout = monitorReq.CreateLBMonitorReq.TfPassiveConfig.Timeout
	} else if monitorReq.CreateLBMonitorReq.TfTCPConfig != nil {
		monitorReq.CreateLBMonitorReq.AliasPort = monitorReq.CreateLBMonitorReq.TfTCPConfig.AliasPort
		monitorReq.CreateLBMonitorReq.FallCount = monitorReq.CreateLBMonitorReq.TfTCPConfig.FallCount
		monitorReq.CreateLBMonitorReq.Interval = monitorReq.CreateLBMonitorReq.TfTCPConfig.Interval
		monitorReq.CreateLBMonitorReq.RequestBody = monitorReq.CreateLBMonitorReq.TfTCPConfig.RequestBody
		monitorReq.CreateLBMonitorReq.ResponseData = monitorReq.CreateLBMonitorReq.TfTCPConfig.ResponseData
		monitorReq.CreateLBMonitorReq.RiseCount = monitorReq.CreateLBMonitorReq.TfTCPConfig.RiseCount
		monitorReq.CreateLBMonitorReq.Timeout = monitorReq.CreateLBMonitorReq.TfTCPConfig.Timeout
	} else if monitorReq.CreateLBMonitorReq.TfUDPConfig != nil {
		monitorReq.CreateLBMonitorReq.AliasPort = monitorReq.CreateLBMonitorReq.TfUDPConfig.AliasPort
		monitorReq.CreateLBMonitorReq.FallCount = monitorReq.CreateLBMonitorReq.TfUDPConfig.FallCount
		monitorReq.CreateLBMonitorReq.Interval = monitorReq.CreateLBMonitorReq.TfUDPConfig.Interval
		monitorReq.CreateLBMonitorReq.RequestBody = monitorReq.CreateLBMonitorReq.TfUDPConfig.RequestBody
		monitorReq.CreateLBMonitorReq.ResponseData = monitorReq.CreateLBMonitorReq.TfUDPConfig.ResponseData
		monitorReq.CreateLBMonitorReq.RiseCount = monitorReq.CreateLBMonitorReq.TfUDPConfig.RiseCount
		monitorReq.CreateLBMonitorReq.Timeout = monitorReq.CreateLBMonitorReq.TfUDPConfig.Timeout
	}

	return nil
}
