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

	createReq := models.CreateLBProfile{
		CreateLBProfileReq: models.CreateLBProfileReq{
			Name:        d.GetString("name"),
			LbID:        d.GetInt("lb_id"),
			Description: d.GetString("description"),
			ServiceType: d.GetString("service_type"),
			ProfileConfig: models.LBProfile{
				ProfileType:            d.GetString("profile_type"),
				RequestHeaderSize:      d.GetInt("request_header_size"),
				ResponseHeaderSize:     d.GetInt("response_header_size"),
				ResponseTimeout:        d.GetInt("response_timeout"),
				HTTPIdleTimeoutName:    d.GetInt("http_idle_timeout"),
				FastTCPIdleTimeout:     d.GetInt("fast_tcp_idle_timeout"),
				ConnectionCloseTimeout: d.GetInt("connection_close_timeout"),
				HaFlowMirroring:        d.GetBool("ha_flow_mirroring"),
				CookieMode:             d.GetString("cookie_mode"),
				CookieName:             d.GetString("cookie_name"),
				CookieType:             d.GetString("cookie_type"),
				CookieFallback:         d.GetBool("cookie_fallback"),
				CookieGarbling:         d.GetBool("cookie_garbling"),
				SSLSuite:               d.GetString("ssl_suite"),
			},
		},
	}

	createReq.CreateLBProfileReq.ProfileConfig.ProfileType = "application-profile"
	createReq.CreateLBProfileReq.ProfileConfig.ConnectionCloseTimeout = 15
	createReq.CreateLBProfileReq.ProfileConfig.FastTCPIdleTimeout = 15
	createReq.CreateLBProfileReq.ProfileConfig.RequestHeaderSize = 30
	createReq.CreateLBProfileReq.ProfileConfig.ResponseHeaderSize = 40
	createReq.CreateLBProfileReq.ProfileConfig.ResponseTimeout = 40
	createReq.CreateLBProfileReq.ProfileConfig.HTTPIdleTimeoutName = 50

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

	updateReq := models.CreateLBProfile{
		CreateLBProfileReq: models.CreateLBProfileReq{
			Name:        d.GetString("name"),
			LbID:        d.GetInt("lb_id"),
			Description: d.GetString("description"),
			ServiceType: d.GetString("service_type"),
			ProfileConfig: models.LBProfile{
				ProfileType:            d.GetString("profile_type"),
				RequestHeaderSize:      d.GetInt("request_header_size"),
				ResponseHeaderSize:     d.GetInt("response_header_size"),
				ResponseTimeout:        d.GetInt("response_timeout"),
				HTTPIdleTimeoutName:    d.GetInt("http_idle_timeout"),
				FastTCPIdleTimeout:     d.GetInt("fast_tcp_idle_timeout"),
				ConnectionCloseTimeout: d.GetInt("connection_close_timeout"),
				HaFlowMirroring:        d.GetBool("ha_flow_mirroring"),
				CookieMode:             d.GetString("cookie_mode"),
				CookieName:             d.GetString("cookie_name"),
				CookieType:             d.GetString("cookie_type"),
				CookieFallback:         d.GetBool("cookie_fallback"),
				CookieGarbling:         d.GetBool("cookie_garbling"),
				SSLSuite:               d.GetString("ssl_suite"),
			},
		},
	}

	updateReq.CreateLBProfileReq.ProfileConfig.ProfileType = "application-profile"
	updateReq.CreateLBProfileReq.ProfileConfig.ConnectionCloseTimeout = 15
	updateReq.CreateLBProfileReq.ProfileConfig.FastTCPIdleTimeout = 15
	updateReq.CreateLBProfileReq.ProfileConfig.RequestHeaderSize = 30
	updateReq.CreateLBProfileReq.ProfileConfig.ResponseHeaderSize = 40
	updateReq.CreateLBProfileReq.ProfileConfig.ResponseTimeout = 40
	updateReq.CreateLBProfileReq.ProfileConfig.HTTPIdleTimeoutName = 50

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
