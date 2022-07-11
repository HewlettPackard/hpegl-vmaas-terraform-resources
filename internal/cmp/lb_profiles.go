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

type loadBalancerProfile struct {
	lbClient *client.LoadBalancerAPIService
}

func newLoadBalancerProfile(loadBalancerClient *client.LoadBalancerAPIService) *loadBalancerProfile {
	return &loadBalancerProfile{
		lbClient: loadBalancerClient,
	}
}

func (lb *loadBalancerProfile) Read(ctx context.Context, d *utils.Data, meta interface{}) error {
	setMeta(meta, lb.lbClient.Client)
	var lbProfileResp models.CreateLBProfileReq

	if err := tftags.Get(d, &lbProfileResp); err != nil {
		return err
	}

	getProfileLoadBalancer, err := lb.lbClient.GetSpecificLBProfile(ctx, lbProfileResp.LbID, lbProfileResp.ID)
	if err != nil {
		return err
	}
	return tftags.Set(d, getProfileLoadBalancer.GetLBSpecificProfilesResp)
}

func (lb *loadBalancerProfile) Create(ctx context.Context, d *utils.Data, meta interface{}) error {
	setMeta(meta, lb.lbClient.Client)
	var createReq models.CreateLBProfile
	if err := tftags.Get(d, &createReq.CreateLBProfileReq); err != nil {
		return err
	}

	if createReq.CreateLBProfileReq.ProfileConfig.ProfileType == "" {
		createReq.CreateLBProfileReq.ProfileConfig.ProfileType = ApplicationProfile
	}
	createReq.CreateLBProfileReq.ProfileConfig.ConnectionCloseTimeout = 10
	createReq.CreateLBProfileReq.ProfileConfig.RequestHeaderSize = RequestHeaderSize
	createReq.CreateLBProfileReq.ProfileConfig.ResponseHeaderSize = ResponseHeaderSize
	createReq.CreateLBProfileReq.ProfileConfig.ResponseTimeout = ResponseTimeout

	createReq.CreateLBProfileReq.ProfileConfig.FastTCPIdleTimeout = 200
	createReq.CreateLBProfileReq.ProfileConfig.FastUDPIdleTimeout = 300
	createReq.CreateLBProfileReq.ProfileConfig.CookieName = "Cookie1"
	createReq.CreateLBProfileReq.ProfileConfig.CookieMode = "INSERT"
	createReq.CreateLBProfileReq.ProfileConfig.CookieType = "LBPersistenceCookieTime"
	createReq.CreateLBProfileReq.ProfileConfig.PersistenceEntryTimeout = 300
	createReq.CreateLBProfileReq.ProfileConfig.SessionCacheEntryTimeout = 300
	createReq.CreateLBProfileReq.ProfileConfig.SSLSuite = "BALANCED"

	lbProfileResp, err := lb.lbClient.CreateLBProfile(ctx, createReq, createReq.CreateLBProfileReq.LbID)
	if err != nil {
		return err
	}

	if !lbProfileResp.Success {
		return fmt.Errorf(successErr, "creating loadBalancerProfile Profile")
	}
	createReq.CreateLBProfileReq.ID = lbProfileResp.LBProfileResp.ID

	// wait until created
	retry := &utils.CustomRetry{
		InitialDelay: time.Second * 15,
		RetryDelay:   time.Second * 30,
	}
	_, err = retry.Retry(ctx, meta, func(ctx context.Context) (interface{}, error) {
		return lb.lbClient.GetSpecificLBProfile(ctx, createReq.CreateLBProfileReq.LbID,
			lbProfileResp.LBProfileResp.ID)
	})
	if err != nil {
		return err
	}

	return tftags.Set(d, createReq.CreateLBProfileReq)
}

func (lb *loadBalancerProfile) Delete(ctx context.Context, d *utils.Data, meta interface{}) error {
	setMeta(meta, lb.lbClient.Client)
	var tfLBProfile models.CreateLBProfileReq
	if err := tftags.Get(d, &tfLBProfile); err != nil {
		return err
	}

	resp, err := lb.lbClient.DeleteLBProfile(ctx, tfLBProfile.LbID, tfLBProfile.ID)
	if err != nil {
		return err
	}

	if !resp.Success {
		return fmt.Errorf("got success = 'false' while deleting LB-PROFILE")
	}

	return nil
}

func (lb *loadBalancerProfile) Update(ctx context.Context, d *utils.Data, meta interface{}) error {
	id := d.GetID()
	setMeta(meta, lb.lbClient.Client)
	var updateReq models.CreateLBProfile
	if err := tftags.Get(d, &updateReq.CreateLBProfileReq); err != nil {
		return err
	}

	updateReq.CreateLBProfileReq.ProfileConfig.RequestHeaderSize = RequestHeaderSize
	updateReq.CreateLBProfileReq.ProfileConfig.ResponseHeaderSize = ResponseHeaderSize
	updateReq.CreateLBProfileReq.ProfileConfig.ResponseTimeout = ResponseTimeout
	if updateReq.CreateLBProfileReq.ProfileConfig.HTTPIdleTimeout == 0 {
		updateReq.CreateLBProfileReq.ProfileConfig.HTTPIdleTimeout = HTTPIdleTimeout

	}

	updateReq.CreateLBProfileReq.ProfileConfig.FastTCPIdleTimeout = 200
	updateReq.CreateLBProfileReq.ProfileConfig.FastUDPIdleTimeout = 300
	updateReq.CreateLBProfileReq.ProfileConfig.CookieName = "Cookie"
	updateReq.CreateLBProfileReq.ProfileConfig.CookieMode = "INSERT"
	updateReq.CreateLBProfileReq.ProfileConfig.CookieType = "LBPersistenceCookieTime"
	updateReq.CreateLBProfileReq.ProfileConfig.PersistenceEntryTimeout = 300
	updateReq.CreateLBProfileReq.ProfileConfig.SessionCacheEntryTimeout = 300
	updateReq.CreateLBProfileReq.ProfileConfig.SSLSuite = "BALANCED"

	retry := &utils.CustomRetry{
		InitialDelay: time.Second * 15,
		RetryDelay:   time.Second * 30,
	}
	_, err := retry.Retry(ctx, meta, func(ctx context.Context) (interface{}, error) {
		return lb.lbClient.UpdateLBProfile(ctx, updateReq,
			updateReq.CreateLBProfileReq.LbID, id)
	})
	if err != nil {
		return err
	}

	return tftags.Set(d, updateReq.CreateLBProfileReq)
}
