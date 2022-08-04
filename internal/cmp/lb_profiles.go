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
	createReq := models.CreateLBProfile{}
	if err := tftags.Get(d, &createReq.CreateLBProfileReq); err != nil {
		return err
	}

	// align createReq and fill json related fields
	if err := lb.profileAlignprofileTypeRequest(ctx, meta, &createReq.CreateLBProfileReq); err != nil {
		return err
	}

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
	updateReq := models.CreateLBProfile{}
	if err := tftags.Get(d, &updateReq.CreateLBProfileReq); err != nil {
		return err
	}

	// align createReq and fill json related fields
	if err := lb.profileAlignprofileTypeRequest(ctx, meta, &updateReq.CreateLBProfileReq); err != nil {
		return err
	}

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

func (lb *loadBalancerProfile) profileAlignprofileTypeRequest(ctx context.Context, meta interface{}, profileReq *models.CreateLBProfileReq) error {
	if profileReq.TfHTTPConfig != nil {
		profileReq.ProfileConfig.HTTPIdleTimeout = profileReq.TfHTTPConfig.HTTPIdleTimeout
		profileReq.ProfileConfig.HTTPSRedirect = profileReq.TfHTTPConfig.HTTPSRedirect
		profileReq.ProfileConfig.NtlmAuthentication = profileReq.TfHTTPConfig.NtlmAuthentication
		profileReq.ProfileConfig.ProfileType = profileReq.ProfileType
		profileReq.ServiceType = profileReq.TfHTTPConfig.ServiceType
		profileReq.ProfileConfig.RequestBodySize = profileReq.TfHTTPConfig.RequestBodySize
		profileReq.ProfileConfig.RequestHeaderSize = profileReq.TfHTTPConfig.RequestHeaderSize
		profileReq.ProfileConfig.ResponseHeaderSize = profileReq.TfHTTPConfig.ResponseHeaderSize
		profileReq.ProfileConfig.ResponseTimeout = profileReq.TfHTTPConfig.ResponseTimeout
		profileReq.ProfileConfig.XForwardedFor = profileReq.TfHTTPConfig.XForwardedFor
	} else if profileReq.TfTCPConfig != nil {
		profileReq.ProfileConfig.ConnectionCloseTimeout = profileReq.TfTCPConfig.ConnectionCloseTimeout
		profileReq.ProfileConfig.FastTCPIdleTimeout = profileReq.TfTCPConfig.FastTCPIdleTimeout
		profileReq.ProfileConfig.HaFlowMirroring = profileReq.TfTCPConfig.HaFlowMirroring
		profileReq.ProfileConfig.ProfileType = profileReq.ProfileType
		profileReq.ServiceType = profileReq.TfTCPConfig.ServiceType

	} else if profileReq.TfUDPConfig != nil {
		profileReq.ProfileConfig.FastUDPIdleTimeout = profileReq.TfUDPConfig.FastUDPIdleTimeout
		profileReq.ProfileConfig.HaFlowMirroring = profileReq.TfUDPConfig.HaFlowMirroring
		profileReq.ServiceType = profileReq.TfUDPConfig.ServiceType
		profileReq.ProfileConfig.ProfileType = profileReq.ProfileType

	} else if profileReq.TfCookieConfig != nil {
		profileReq.ProfileConfig.CookieDomain = profileReq.TfCookieConfig.CookieDomain
		profileReq.ProfileConfig.CookieFallback = profileReq.TfCookieConfig.CookieFallback
		profileReq.ProfileConfig.CookieGarbling = profileReq.TfCookieConfig.CookieGarbling
		profileReq.ProfileConfig.CookieMode = profileReq.TfCookieConfig.CookieMode
		profileReq.ProfileConfig.CookieName = profileReq.TfCookieConfig.CookieName
		profileReq.ProfileConfig.CookiePath = profileReq.TfCookieConfig.CookiePath
		profileReq.ProfileConfig.CookieType = profileReq.TfCookieConfig.CookieType
		profileReq.ProfileConfig.MaxCookieAge = profileReq.TfCookieConfig.MaxCookieAge
		profileReq.ProfileConfig.MaxIdleTime = profileReq.TfCookieConfig.MaxIdleTime
		profileReq.ServiceType = profileReq.TfCookieConfig.ServiceType
		profileReq.ProfileConfig.ProfileType = profileReq.ProfileType
		profileReq.ProfileConfig.SharePersistence = profileReq.TfCookieConfig.SharePersistence

	} else if profileReq.TfGenericConfig != nil {
		profileReq.ProfileConfig.HaPersistenceMirroring = profileReq.TfGenericConfig.HaPersistenceMirroring
		profileReq.ProfileConfig.PersistenceEntryTimeout = profileReq.TfGenericConfig.PersistenceEntryTimeout
		profileReq.ProfileConfig.ProfileType = profileReq.ProfileType
		profileReq.ServiceType = profileReq.TfGenericConfig.ServiceType
		profileReq.ProfileConfig.SharePersistence = profileReq.TfGenericConfig.SharePersistence

	} else if profileReq.TfSourceConfig != nil {
		profileReq.ProfileConfig.HaPersistenceMirroring = profileReq.TfSourceConfig.HaPersistenceMirroring
		profileReq.ProfileConfig.PersistenceEntryTimeout = profileReq.TfSourceConfig.PersistenceEntryTimeout
		profileReq.ProfileConfig.ProfileType = profileReq.ProfileType
		profileReq.ServiceType = profileReq.TfSourceConfig.ServiceType
		profileReq.ProfileConfig.PurgeEntries = profileReq.TfSourceConfig.PurgeEntries
		profileReq.ProfileConfig.SharePersistence = profileReq.TfSourceConfig.SharePersistence

	} else if profileReq.TfServerConfig != nil {
		profileReq.ProfileConfig.ProfileType = profileReq.ProfileType
		profileReq.ProfileConfig.SSLSuite = profileReq.TfServerConfig.SSLSuite
		profileReq.ServiceType = profileReq.TfServerConfig.ServiceType
		profileReq.ProfileConfig.SessionCache = profileReq.TfServerConfig.SessionCache

	} else if profileReq.TfClientConfig != nil {
		profileReq.ProfileConfig.PreferServerCipher = profileReq.TfClientConfig.PreferServerCipher
		profileReq.ProfileConfig.ProfileType = profileReq.ProfileType
		profileReq.ProfileConfig.SSLSuite = profileReq.TfClientConfig.SSLSuite
		profileReq.ServiceType = profileReq.TfClientConfig.ServiceType
		profileReq.ProfileConfig.SessionCache = profileReq.TfClientConfig.SessionCache
		profileReq.ProfileConfig.SessionCacheEntryTimeout = profileReq.TfClientConfig.SessionCacheEntryTimeout
	}
	return nil
}
