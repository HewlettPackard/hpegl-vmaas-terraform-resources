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
	// var createReq models.CreateLBProfile
	// if err := tftags.Get(d, &createReq.CreateLBProfileReq); err != nil {
	// 	return err
	// }
	createReq := models.CreateLBProfile{
		CreateLBProfileReq: models.CreateLBProfileReq{
			Name:        d.GetString("name"),
			LbID:        d.GetInt("lb_id"),
			Description: d.GetString("description"),
			ServiceType: d.GetString("service_type"),
			ProfileConfig: models.LBProfile{
				ProfileType:              d.GetString("profile_type"),
				FastTCPIdleTimeout:       d.GetInt("fast_tcp_idle_timeout"),
				FastUDPIdleTimeout:       d.GetInt("fast_udp_idle_timeout"),
				HTTPIdleTimeout:          d.GetInt("http_idle_timeout"),
				ConnectionCloseTimeout:   d.GetInt("connection_close_timeout"),
				HaFlowMirroring:          d.GetBool("ha_flow_mirroring"),
				RequestHeaderSize:        d.GetInt("request_header_size"),
				ResponseHeaderSize:       d.GetInt("response_header_size"),
				HTTPsRedirect:            d.GetString("redirection"),
				XForwardedFor:            d.GetString("x_forwarded_for"),
				RequestBodySize:          d.GetString("request_body_size"),
				ResponseTimeout:          d.GetInt("response_timeout"),
				NtlmAuthentication:       d.GetBool("ntlm_authentication"),
				SharePersistence:         d.GetBool("share_persistence"),
				CookieName:               d.GetString("cookie_name"),
				CookieFallback:           d.GetBool("cookie_fallback"),
				CookieGarbling:           d.GetBool("cookie_garbling"),
				CookieMode:               d.GetString("cookie_mode"),
				CookieType:               d.GetString("cookie_type"),
				CookieDomain:             d.GetString("cookie_domain"),
				CookiePath:               d.GetString("cookie_path"),
				MaxIdleTime:              d.GetInt("max_idle_time"),
				MaxCookieAge:             d.GetInt("max_cookie_age"),
				HaPersistenceMirroring:   d.GetBool("ha_persistence_mirroring"),
				PersistenceEntryTimeout:  d.GetInt("persistence_entry_timeout"),
				PurgeEntries:             d.GetBool("purge_entries_when_full"),
				SSLSuite:                 d.GetString("ssl_suite"),
				SessionCache:             d.GetBool("session_cache"),
				SessionCacheEntryTimeout: d.GetInt("session_cache_timeout"),
				PreferServerCipher:       d.GetBool("prefer_server_cipher"),
				Tag: []models.Tags{
					{
						Tag:   d.GetString("tag"),
						Scope: d.GetString("scope"),
					},
				},
			},
		},
	}

	// if createReq.CreateLBProfileReq.ProfileConfig.ProfileType == "" {
	// 	if createReq.CreateLBProfileReq.ProfileConfig.ProfileType == ApplicationProfile &&
	// 		createReq.CreateLBProfileReq.ServiceType == "LBHttpProfile" {
	// 		createReq.CreateLBProfileReq.ProfileConfig.HTTPIdleTimeout = 15
	// 	} else if createReq.CreateLBProfileReq.ProfileConfig.ProfileType == ApplicationProfile &&
	// 		createReq.CreateLBProfileReq.ServiceType == "LBFastTcpProfile" {
	// 		createReq.CreateLBProfileReq.ProfileConfig.FastTCPIdleTimeout = 1800
	// 	} else if createReq.CreateLBProfileReq.ProfileConfig.ProfileType == ApplicationProfile &&
	// 		createReq.CreateLBProfileReq.ServiceType == "LBFastUdpProfile" {
	// 		createReq.CreateLBProfileReq.ProfileConfig.FastUDPIdleTimeout = 300
	// 	}
	// }

	if createReq.CreateLBProfileReq.ProfileConfig.ProfileType == "" {
		createReq.CreateLBProfileReq.ProfileConfig.ProfileType = ApplicationProfile
	}
	createReq.CreateLBProfileReq.ProfileConfig.ConnectionCloseTimeout = 10
	createReq.CreateLBProfileReq.ProfileConfig.RequestHeaderSize = RequestHeaderSize
	createReq.CreateLBProfileReq.ProfileConfig.ResponseHeaderSize = ResponseHeaderSize
	createReq.CreateLBProfileReq.ProfileConfig.ResponseTimeout = ResponseTimeout

	// createReq.CreateLBProfileReq.ProfileConfig.ConnectionCloseTimeout = 59
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
	// var updateReq models.CreateLBProfile
	// if err := tftags.Get(d, &updateReq.CreateLBProfileReq); err != nil {
	// 	return err
	// }

	updateReq := models.CreateLBProfile{
		CreateLBProfileReq: models.CreateLBProfileReq{
			Name:        d.GetString("name"),
			LbID:        d.GetInt("lb_id"),
			Description: d.GetString("description"),
			ServiceType: d.GetString("service_type"),
			ProfileConfig: models.LBProfile{
				ProfileType:              d.GetString("profile_type"),
				FastTCPIdleTimeout:       d.GetInt("fast_tcp_idle_timeout"),
				FastUDPIdleTimeout:       d.GetInt("fast_udp_idle_timeout"),
				HTTPIdleTimeout:          d.GetInt("http_idle_timeout"),
				ConnectionCloseTimeout:   d.GetInt("connection_close_timeout"),
				HaFlowMirroring:          d.GetBool("ha_flow_mirroring"),
				RequestHeaderSize:        d.GetInt("request_header_size"),
				ResponseHeaderSize:       d.GetInt("response_header_size"),
				HTTPsRedirect:            d.GetString("redirection"),
				XForwardedFor:            d.GetString("x_forwarded_for"),
				RequestBodySize:          d.GetString("request_body_size"),
				ResponseTimeout:          d.GetInt("response_timeout"),
				NtlmAuthentication:       d.GetBool("ntlm_authentication"),
				SharePersistence:         d.GetBool("share_persistence"),
				CookieName:               d.GetString("cookie_name"),
				CookieFallback:           d.GetBool("cookie_fallback"),
				CookieGarbling:           d.GetBool("cookie_garbling"),
				CookieMode:               d.GetString("cookie_mode"),
				CookieType:               d.GetString("cookie_type"),
				CookieDomain:             d.GetString("cookie_domain"),
				CookiePath:               d.GetString("cookie_path"),
				MaxIdleTime:              d.GetInt("max_idle_time"),
				MaxCookieAge:             d.GetInt("max_cookie_age"),
				HaPersistenceMirroring:   d.GetBool("ha_persistence_mirroring"),
				PersistenceEntryTimeout:  d.GetInt("persistence_entry_timeout"),
				PurgeEntries:             d.GetBool("purge_entries_when_full"),
				SSLSuite:                 d.GetString("ssl_suite"),
				SessionCache:             d.GetBool("session_cache"),
				SessionCacheEntryTimeout: d.GetInt("session_cache_timeout"),
				PreferServerCipher:       d.GetBool("prefer_server_cipher"),
				Tag: []models.Tags{
					{
						Tag:   d.GetString("tag"),
						Scope: d.GetString("scope"),
					},
				},
			},
		},
	}

	updateReq.CreateLBProfileReq.ProfileConfig.RequestHeaderSize = RequestHeaderSize
	updateReq.CreateLBProfileReq.ProfileConfig.ResponseHeaderSize = ResponseHeaderSize
	updateReq.CreateLBProfileReq.ProfileConfig.ResponseTimeout = ResponseTimeout
	if updateReq.CreateLBProfileReq.ProfileConfig.HTTPIdleTimeout == 0 {
		updateReq.CreateLBProfileReq.ProfileConfig.HTTPIdleTimeout = HTTPIdleTimeout

	}

	// updateReq.CreateLBProfileReq.ProfileConfig.ConnectionCloseTimeout = 59
	updateReq.CreateLBProfileReq.ProfileConfig.FastTCPIdleTimeout = 200
	updateReq.CreateLBProfileReq.ProfileConfig.FastUDPIdleTimeout = 300
	updateReq.CreateLBProfileReq.ProfileConfig.CookieName = "Cookie1"
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
