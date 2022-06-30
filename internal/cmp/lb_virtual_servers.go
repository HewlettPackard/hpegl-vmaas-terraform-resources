// (C) Copyright 2022 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"
	"fmt"

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

	_, err := lb.lbClient.GetSpecificLBVirtualServer(ctx, lbVSResp.LbID, lbVSResp.ID)
	if err != nil {
		return err
	}
	return tftags.Set(d, lbVSResp)
}

func (lb *loadBalancerVirtualServer) Update(ctx context.Context, d *utils.Data, meta interface{}) error {
	return nil
}

func (lb *loadBalancerVirtualServer) Create(ctx context.Context, d *utils.Data, meta interface{}) error {
	setMeta(meta, lb.lbClient.Client)
	var createReq models.CreateLBVirtualServers

	// createReq := models.CreateLBVirtualServers{
	// 	CreateLBVirtualServersReq: models.CreateLBVirtualServersReq{
	// 		Description:   d.GetString("description"),
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

	if err := tftags.Get(d, &createReq.CreateLBVirtualServersReq); err != nil {
		return err
	}

	profileData, err := lb.lbClient.GetLBProfiles(ctx, createReq.CreateLBVirtualServersReq.LbID)
	if err != nil {
		return err
	}

	for i, profile := range profileData.GetLBProfilesResp {
		if profile.LBProfileConfig.ProfileType == ApplicationProfile &&
			profile.ServiceType == ServiceTypeLBHttpProfile {
			createReq.CreateLBVirtualServersReq.VirtualServerConfig.ApplicationProfile =
				profileData.GetLBProfilesResp[i].ID
			break
		}
	}
	//createReq.CreateLBVirtualServersReq.Pool = poolID.GetLBPoolsResp[0].ID

	lbVirtualServersResp, err := lb.lbClient.CreateLBVirtualServers(ctx, createReq, createReq.CreateLBVirtualServersReq.LbID)
	if err != nil {
		return err
	}
	if !lbVirtualServersResp.Success {
		return fmt.Errorf(successErr, "creating loadBalancerVirtualServer Virtual Servers")
	}

	createReq.CreateLBVirtualServersReq.ID = lbVirtualServersResp.CreateLBVirtualServersResp.ID
	return tftags.Set(d, createReq.CreateLBVirtualServersReq)
}

func (lb *loadBalancerVirtualServer) Delete(ctx context.Context, d *utils.Data, meta interface{}) error {
	lbVirtualServerID := d.GetID()
	lbDetails, err := lb.lbClient.GetLoadBalancers(ctx)
	if err != nil {
		return err
	}
	_, err = lb.lbClient.DeleteLBVirtualServers(ctx, lbDetails.GetNetworkLoadBalancerResp[0].ID, lbVirtualServerID)
	if err != nil {
		return err
	}

	return nil
}
