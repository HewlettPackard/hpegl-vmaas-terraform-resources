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
	var lbProfileResp models.GetLBSpecificProfilesResp
	if err := tftags.Get(d, &lbProfileResp); err != nil {
		return err
	}

	lbDetails, err := lb.lbClient.GetLoadBalancers(ctx)
	if err != nil {
		return err
	}

	getlbProfileResp, err := lb.lbClient.GetSpecificLBProfile(ctx, lbDetails.GetNetworkLoadBalancerResp[0].ID, lbProfileResp.ID)
	if err != nil {
		return err
	}

	return tftags.Set(d, getlbProfileResp.GetLBSpecificProfilesResp)
}

func (lb *loadBalancerProfile) Update(ctx context.Context, d *utils.Data, meta interface{}) error {
	return nil
}

func (lb *loadBalancerProfile) Create(ctx context.Context, d *utils.Data, meta interface{}) error {
	setMeta(meta, lb.lbClient.Client)
	createReq := models.CreateLBProfile{
		CreateLBProfileReq: models.CreateLBProfileReq{
			Name:        d.GetString("name"),
			Description: d.GetString("description"),
			ServiceType: d.GetString("service_type"),
			ProfileConfig: models.LBProfile{
				ProfileType:            d.GetString("profile_type"),
				RequestHeaderSize:      d.GetInt("request_header_size"),
				ResponseHeaderSize:     d.GetInt("response_header_size"),
				ResponseTimeout:        d.GetInt("response_timeout"),
				HTTPIdleTimeout:        d.GetInt("http_idle_timeout"),
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

	if err := tftags.Get(d, &createReq.CreateLBProfileReq); err != nil {
		return err
	}

	lbDetails, err := lb.lbClient.GetLoadBalancers(ctx)
	if err != nil {
		return err
	}

	createReq.CreateLBProfileReq.ProfileConfig.ProfileType = "application-profile"
	createReq.CreateLBProfileReq.ProfileConfig.ConnectionCloseTimeout = 15
	createReq.CreateLBProfileReq.ProfileConfig.FastTCPIdleTimeout = 15
	createReq.CreateLBProfileReq.ProfileConfig.RequestHeaderSize = 30
	createReq.CreateLBProfileReq.ProfileConfig.ResponseHeaderSize = 40
	createReq.CreateLBProfileReq.ProfileConfig.ResponseTimeout = 40
	createReq.CreateLBProfileReq.ProfileConfig.HTTPIdleTimeout = 50

	lbProfileResp, err := lb.lbClient.CreateLBProfile(ctx, createReq, lbDetails.GetNetworkLoadBalancerResp[0].ID)
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
		return lb.lbClient.GetSpecificLBProfile(ctx, lbDetails.GetNetworkLoadBalancerResp[0].ID, lbProfileResp.LBProfileResp.ID)
	})
	if err != nil {
		return err
	}

	return tftags.Set(d, createReq.CreateLBProfileReq)
}

func (lb *loadBalancerProfile) Delete(ctx context.Context, d *utils.Data, meta interface{}) error {
	lbProfileID := d.GetID()
	lbDetails, err := lb.lbClient.GetLoadBalancers(ctx)
	if err != nil {
		return err
	}
	_, err = lb.lbClient.DeleteLBProfile(ctx, lbDetails.GetNetworkLoadBalancerResp[0].ID, lbProfileID)
	if err != nil {
		return err
	}

	return nil
}
