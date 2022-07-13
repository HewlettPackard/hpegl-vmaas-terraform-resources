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

type loadBalancerVirtualServer struct {
	lbClient *client.LoadBalancerAPIService
}

func newLoadBalancerVirtualServer(loadBalancerClient *client.LoadBalancerAPIService) *loadBalancerVirtualServer {
	return &loadBalancerVirtualServer{
		lbClient: loadBalancerClient,
	}
}

func (lb *loadBalancerVirtualServer) Read(ctx context.Context, d *utils.Data, meta interface{}) error {
	setMeta(meta, lb.lbClient.Client)
	var lbVSResp models.CreateLBVirtualServersReq

	if err := tftags.Get(d, &lbVSResp); err != nil {
		return err
	}

	getPoolLoadBalancer, err := lb.lbClient.GetSpecificLBVirtualServer(ctx, lbVSResp.LbID, lbVSResp.ID)
	if err != nil {
		return err
	}
	return tftags.Set(d, getPoolLoadBalancer.GetSpecificLBVirtualServersResp)
}

func (lb *loadBalancerVirtualServer) Update(ctx context.Context, d *utils.Data, meta interface{}) error {

	id := d.GetID()
	var updateReq models.CreateLBVirtualServers
	if err := tftags.Get(d, &updateReq); err != nil {
		return err
	}

	// updateReq := models.CreateLBVirtualServers{
	// 	CreateLBVirtualServersReq: models.CreateLBVirtualServersReq{
	// 		Description:   d.GetString("description"),
	// 		LbID:          d.GetInt("lb_id"),
	// 		VipName:       d.GetString("vip_name"),
	// 		VipAddress:    d.GetString("vip_address"),
	// 		VipProtocol:   d.GetString("vip_protocol"),
	// 		VipPort:       d.GetString("vip_port"),
	// 		Pool:          d.GetInt("pool"),
	// 		SSLServerCert: d.GetInt("ssl_server_cert"),
	// 		SSLCert:       d.GetInt("ssl_cert"),
	// 		VirtualServerConfig: models.VirtualServerConfig{
	// 			Persistence:        d.GetString("persistence"),
	// 			PersistenceProfile: d.GetInt("persistence_profile"),
	// 			ApplicationProfile: d.GetInt("application_profile"),
	// 			SSLClientProfile:   d.GetString("ssl_client_profile"),
	// 			SSLServerProfile:   d.GetString("ssl_server_profile"),
	// 		},
	// 	},
	// }

	retry := &utils.CustomRetry{
		InitialDelay: time.Second * 15,
		RetryDelay:   time.Second * 30,
	}
	_, err := retry.Retry(ctx, meta, func(ctx context.Context) (interface{}, error) {
		return lb.lbClient.UpdateLBVirtualServers(ctx, updateReq,
			updateReq.CreateLBVirtualServersReq.LbID, id)
	})
	if err != nil {
		return err
	}

	return tftags.Set(d, updateReq.CreateLBVirtualServersReq)
}

func (lb *loadBalancerVirtualServer) Create(ctx context.Context, d *utils.Data, meta interface{}) error {
	setMeta(meta, lb.lbClient.Client)
	var createReq models.CreateLBVirtualServers
	if err := tftags.Get(d, &createReq); err != nil {
		return err
	}

	// createReq := models.CreateLBVirtualServers{
	// 	CreateLBVirtualServersReq: models.CreateLBVirtualServersReq{
	// 		Description:   d.GetString("description"),
	// 		LbID:          d.GetInt("lb_id"),
	// 		VipName:       d.GetString("vip_name"),
	// 		VipAddress:    d.GetString("vip_address"),
	// 		VipProtocol:   d.GetString("vip_protocol"),
	// 		VipPort:       d.GetString("vip_port"),
	// 		Pool:          d.GetInt("pool"),
	// 		SSLServerCert: d.GetInt("ssl_server_cert"),
	// 		SSLCert:       d.GetInt("ssl_cert"),
	// 		VirtualServerConfig: models.VirtualServerConfig{
	// 			Persistence:        d.GetString("persistence"),
	// 			PersistenceProfile: d.GetInt("persistence_profile"),
	// 			ApplicationProfile: d.GetInt("application_profile"),
	// 			SSLClientProfile:   d.GetString("ssl_client_profile"),
	// 			SSLServerProfile:   d.GetString("ssl_server_profile"),
	// 		},
	// 	},
	// }

	// profileData, err := lb.lbClient.GetLBProfiles(ctx, createReq.CreateLBVirtualServersReq.LbID)
	// if err != nil {
	// 	return err
	// }

	// for i, profile := range profileData.GetLBProfilesResp {
	// 	if profile.LBProfileConfig.ProfileType == ApplicationProfile &&
	// 		profile.ServiceType == ServiceTypeLBHttpProfile {
	// 		createReq.CreateLBVirtualServersReq.VirtualServerConfig.ApplicationProfile =
	// 			profileData.GetLBProfilesResp[i].ID
	// 		break
	// 	}
	// }

	lbVirtualServersResp, err := lb.lbClient.CreateLBVirtualServers(ctx, createReq, createReq.CreateLBVirtualServersReq.LbID)
	if err != nil {
		return err
	}
	if !lbVirtualServersResp.Success {
		return fmt.Errorf(successErr, "creating loadBalancerVirtualServer Virtual Servers")
	}

	createReq.CreateLBVirtualServersReq.ID = lbVirtualServersResp.CreateLBVirtualServersResp.ID

	// wait until created
	retry := &utils.CustomRetry{
		InitialDelay: time.Second * 15,
		RetryDelay:   time.Second * 30,
	}
	_, err = retry.Retry(ctx, meta, func(ctx context.Context) (interface{}, error) {
		return lb.lbClient.GetSpecificLBVirtualServer(ctx, createReq.CreateLBVirtualServersReq.LbID,
			lbVirtualServersResp.CreateLBVirtualServersResp.ID)
	})
	if err != nil {
		return err
	}

	return tftags.Set(d, createReq.CreateLBVirtualServersReq)
}

func (lb *loadBalancerVirtualServer) Delete(ctx context.Context, d *utils.Data, meta interface{}) error {
	setMeta(meta, lb.lbClient.Client)
	var tfLBVirtualServer models.CreateLBVirtualServersReq
	if err := tftags.Get(d, &tfLBVirtualServer); err != nil {
		return err
	}

	resp, err := lb.lbClient.DeleteLBVirtualServers(ctx, tfLBVirtualServer.LbID, tfLBVirtualServer.ID)
	if err != nil {
		return err
	}

	if !resp.Success {
		return fmt.Errorf("got success = 'false' while deleting LB-VIRTUAL-SERVER")
	}

	return nil
}
