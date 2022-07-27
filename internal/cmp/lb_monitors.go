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
	if err := lb.monitorAlignMonitorTypeRequest(ctx, meta, &createReq); err != nil {
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
	if err := lb.monitorAlignMonitorTypeRequest(ctx, meta, &updateReq); err != nil {
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

func (lb *loadBalancerMonitor) monitorAlignMonitorTypeRequest(ctx context.Context, meta interface{}, monitorReq *models.CreateLBMonitor) error {
	if monitorReq.CreateLBMonitorReq.TfHttpConfig != nil {
		monitorReq.CreateLBMonitorReq.RequestBody = monitorReq.CreateLBMonitorReq.TfHttpConfig.RequestBody
		monitorReq.CreateLBMonitorReq.AliasPort = monitorReq.CreateLBMonitorReq.TfHttpConfig.AliasPort
		monitorReq.CreateLBMonitorReq.FallCount = monitorReq.CreateLBMonitorReq.TfHttpConfig.FallCount
		monitorReq.CreateLBMonitorReq.Interval = monitorReq.CreateLBMonitorReq.TfHttpConfig.Interval
		monitorReq.CreateLBMonitorReq.RequestMethod = monitorReq.CreateLBMonitorReq.TfHttpConfig.RequestMethod
		monitorReq.CreateLBMonitorReq.RequestURL = monitorReq.CreateLBMonitorReq.TfHttpConfig.RequestURL
		monitorReq.CreateLBMonitorReq.RequestVersion = monitorReq.CreateLBMonitorReq.TfHttpConfig.RequestVersion
		monitorReq.CreateLBMonitorReq.ResponseData = monitorReq.CreateLBMonitorReq.TfHttpConfig.ResponseData
		monitorReq.CreateLBMonitorReq.ResponseStatusCodes = monitorReq.CreateLBMonitorReq.TfHttpConfig.ResponseStatusCodes
		monitorReq.CreateLBMonitorReq.RiseCount = monitorReq.CreateLBMonitorReq.TfHttpConfig.RiseCount
		monitorReq.CreateLBMonitorReq.Timeout = monitorReq.CreateLBMonitorReq.TfHttpConfig.Timeout
	} else if monitorReq.CreateLBMonitorReq.TfHttpsConfig != nil {
		monitorReq.CreateLBMonitorReq.RequestBody = monitorReq.CreateLBMonitorReq.TfHttpsConfig.RequestBody
		monitorReq.CreateLBMonitorReq.AliasPort = monitorReq.CreateLBMonitorReq.TfHttpsConfig.AliasPort
		monitorReq.CreateLBMonitorReq.FallCount = monitorReq.CreateLBMonitorReq.TfHttpsConfig.FallCount
		monitorReq.CreateLBMonitorReq.Interval = monitorReq.CreateLBMonitorReq.TfHttpsConfig.Interval
		monitorReq.CreateLBMonitorReq.RequestMethod = monitorReq.CreateLBMonitorReq.TfHttpsConfig.RequestMethod
		monitorReq.CreateLBMonitorReq.RequestURL = monitorReq.CreateLBMonitorReq.TfHttpsConfig.RequestURL
		monitorReq.CreateLBMonitorReq.RequestVersion = monitorReq.CreateLBMonitorReq.TfHttpsConfig.RequestVersion
		monitorReq.CreateLBMonitorReq.ResponseData = monitorReq.CreateLBMonitorReq.TfHttpsConfig.ResponseData
		monitorReq.CreateLBMonitorReq.ResponseStatusCodes = monitorReq.CreateLBMonitorReq.TfHttpsConfig.ResponseStatusCodes
		monitorReq.CreateLBMonitorReq.RiseCount = monitorReq.CreateLBMonitorReq.TfHttpsConfig.RiseCount
		monitorReq.CreateLBMonitorReq.Timeout = monitorReq.CreateLBMonitorReq.TfHttpsConfig.Timeout
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
	} else if monitorReq.CreateLBMonitorReq.TfTcpConfig != nil {
		monitorReq.CreateLBMonitorReq.AliasPort = monitorReq.CreateLBMonitorReq.TfTcpConfig.AliasPort
		monitorReq.CreateLBMonitorReq.FallCount = monitorReq.CreateLBMonitorReq.TfTcpConfig.FallCount
		monitorReq.CreateLBMonitorReq.Interval = monitorReq.CreateLBMonitorReq.TfTcpConfig.Interval
		monitorReq.CreateLBMonitorReq.RequestBody = monitorReq.CreateLBMonitorReq.TfTcpConfig.RequestBody
		monitorReq.CreateLBMonitorReq.ResponseData = monitorReq.CreateLBMonitorReq.TfTcpConfig.ResponseData
		monitorReq.CreateLBMonitorReq.RiseCount = monitorReq.CreateLBMonitorReq.TfTcpConfig.RiseCount
		monitorReq.CreateLBMonitorReq.Timeout = monitorReq.CreateLBMonitorReq.TfTcpConfig.Timeout
	} else if monitorReq.CreateLBMonitorReq.TfUdpConfig != nil {
		monitorReq.CreateLBMonitorReq.AliasPort = monitorReq.CreateLBMonitorReq.TfUdpConfig.AliasPort
		monitorReq.CreateLBMonitorReq.FallCount = monitorReq.CreateLBMonitorReq.TfUdpConfig.FallCount
		monitorReq.CreateLBMonitorReq.Interval = monitorReq.CreateLBMonitorReq.TfUdpConfig.Interval
		monitorReq.CreateLBMonitorReq.RequestBody = monitorReq.CreateLBMonitorReq.TfUdpConfig.RequestBody
		monitorReq.CreateLBMonitorReq.ResponseData = monitorReq.CreateLBMonitorReq.TfUdpConfig.ResponseData
		monitorReq.CreateLBMonitorReq.RiseCount = monitorReq.CreateLBMonitorReq.TfUdpConfig.RiseCount
		monitorReq.CreateLBMonitorReq.Timeout = monitorReq.CreateLBMonitorReq.TfUdpConfig.Timeout
	}
	return nil
}
