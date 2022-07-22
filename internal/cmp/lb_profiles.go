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
	if profileReq.ProfileConfig.TfHttpConfig != nil && profileReq.ServiceType == httpProfile {
		profileReq.ProfileConfig.HTTPIdleTimeout = profileReq.ProfileConfig.TfHttpConfig.HTTPIdleTimeout
		profileReq.ProfileConfig.HTTPsRedirect = profileReq.ProfileConfig.TfHttpConfig.HTTPsRedirect
		profileReq.ProfileConfig.NtlmAuthentication = profileReq.ProfileConfig.TfHttpConfig.NtlmAuthentication
		profileReq.ProfileConfig.ProfileType = profileReq.ProfileConfig.TfHttpConfig.ProfileType
		profileReq.ProfileConfig.RequestBodySize = profileReq.ProfileConfig.TfHttpConfig.RequestBodySize
		profileReq.ProfileConfig.RequestHeaderSize = profileReq.ProfileConfig.TfHttpConfig.RequestHeaderSize
		profileReq.ProfileConfig.ResponseHeaderSize = profileReq.ProfileConfig.TfHttpConfig.ResponseHeaderSize
		profileReq.ProfileConfig.ResponseTimeout = profileReq.ProfileConfig.TfHttpConfig.ResponseTimeout
		profileReq.ProfileConfig.XForwardedFor = profileReq.ProfileConfig.TfHttpConfig.XForwardedFor
	} else if profileReq.ProfileConfig.TfTcpConfig != nil && profileReq.ServiceType == tcpProfile {
		profileReq.ProfileConfig.ConnectionCloseTimeout = profileReq.ProfileConfig.TfTcpConfig.ConnectionCloseTimeout
		profileReq.ProfileConfig.FastTCPIdleTimeout = profileReq.ProfileConfig.TfTcpConfig.FastTCPIdleTimeout
		profileReq.ProfileConfig.HaFlowMirroring = profileReq.ProfileConfig.TfTcpConfig.HaFlowMirroring
		profileReq.ProfileConfig.ProfileType = profileReq.ProfileConfig.TfTcpConfig.ProfileType

	} else if profileReq.ProfileConfig.TfUdpConfig != nil && profileReq.ServiceType == udpProfile {
		profileReq.ProfileConfig.FastUDPIdleTimeout = profileReq.ProfileConfig.TfUdpConfig.FastUDPIdleTimeout
		profileReq.ProfileConfig.HaFlowMirroring = profileReq.ProfileConfig.TfUdpConfig.HaFlowMirroring
		profileReq.ProfileConfig.ProfileType = profileReq.ProfileConfig.TfUdpConfig.ProfileType

	} else if profileReq.ProfileConfig.TfCookieConfig != nil && profileReq.ServiceType == cookieProfile {
		profileReq.ProfileConfig.CookieDomain = profileReq.ProfileConfig.TfCookieConfig.CookieDomain
		profileReq.ProfileConfig.CookieFallback = profileReq.ProfileConfig.TfCookieConfig.CookieFallback
		profileReq.ProfileConfig.CookieGarbling = profileReq.ProfileConfig.TfCookieConfig.CookieGarbling
		profileReq.ProfileConfig.CookieMode = profileReq.ProfileConfig.TfCookieConfig.CookieMode
		profileReq.ProfileConfig.CookieName = profileReq.ProfileConfig.TfCookieConfig.CookieName
		profileReq.ProfileConfig.CookiePath = profileReq.ProfileConfig.TfCookieConfig.CookiePath
		profileReq.ProfileConfig.CookieType = profileReq.ProfileConfig.TfCookieConfig.CookieType
		profileReq.ProfileConfig.MaxCookieAge = profileReq.ProfileConfig.TfCookieConfig.MaxCookieAge
		profileReq.ProfileConfig.MaxCookieLife = profileReq.ProfileConfig.TfCookieConfig.MaxCookieLife
		profileReq.ProfileConfig.MaxIdleTime = profileReq.ProfileConfig.TfCookieConfig.MaxIdleTime
		profileReq.ProfileConfig.ProfileType = profileReq.ProfileConfig.TfCookieConfig.ProfileType
		profileReq.ProfileConfig.SharePersistence = profileReq.ProfileConfig.TfCookieConfig.SharePersistence

	} else if profileReq.ProfileConfig.TfGenericConfig != nil && profileReq.ServiceType == genericProfile {
		profileReq.ProfileConfig.HaPersistenceMirroring = profileReq.ProfileConfig.TfGenericConfig.HaPersistenceMirroring
		profileReq.ProfileConfig.PersistenceEntryTimeout = profileReq.ProfileConfig.TfGenericConfig.PersistenceEntryTimeout
		profileReq.ProfileConfig.ProfileType = profileReq.ProfileConfig.TfGenericConfig.ProfileType
		profileReq.ProfileConfig.SharePersistence = profileReq.ProfileConfig.TfGenericConfig.SharePersistence

	} else if profileReq.ProfileConfig.TfSourceConfig != nil && profileReq.ServiceType == sourceProfile {
		profileReq.ProfileConfig.HaPersistenceMirroring = profileReq.ProfileConfig.TfSourceConfig.HaPersistenceMirroring
		profileReq.ProfileConfig.PersistenceEntryTimeout = profileReq.ProfileConfig.TfSourceConfig.PersistenceEntryTimeout
		profileReq.ProfileConfig.ProfileType = profileReq.ProfileConfig.TfSourceConfig.ProfileType
		profileReq.ProfileConfig.PurgeEntries = profileReq.ProfileConfig.TfSourceConfig.PurgeEntries
		profileReq.ProfileConfig.SharePersistence = profileReq.ProfileConfig.TfSourceConfig.SharePersistence

	} else if profileReq.ProfileConfig.TfServerConfig != nil && profileReq.ServiceType == serverProfile {
		profileReq.ProfileConfig.ProfileType = profileReq.ProfileConfig.TfServerConfig.ProfileType
		profileReq.ProfileConfig.SSLSuite = profileReq.ProfileConfig.TfServerConfig.SSLSuite
		profileReq.ProfileConfig.SessionCache = profileReq.ProfileConfig.TfServerConfig.SessionCache

	} else if profileReq.ProfileConfig.TfClientConfig != nil && profileReq.ServiceType == clientProfile {
		profileReq.ProfileConfig.PreferServerCipher = profileReq.ProfileConfig.TfClientConfig.PreferServerCipher
		profileReq.ProfileConfig.ProfileType = profileReq.ProfileConfig.TfClientConfig.ProfileType
		profileReq.ProfileConfig.SSLSuite = profileReq.ProfileConfig.TfClientConfig.SSLSuite
		profileReq.ProfileConfig.SessionCache = profileReq.ProfileConfig.TfClientConfig.SessionCache
		profileReq.ProfileConfig.SessionCacheEntryTimeout = profileReq.ProfileConfig.TfClientConfig.SessionCacheEntryTimeout
	}
	return nil
}
