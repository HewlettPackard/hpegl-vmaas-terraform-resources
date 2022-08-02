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
	updateReq := models.CreateLBVirtualServers{}
	if err := tftags.Get(d, &updateReq.CreateLBVirtualServersReq); err != nil {
		return err
	}

	// align createReq and fill json related fields
	if err := lb.virtualServerAlignRequest(ctx, meta, &updateReq.CreateLBVirtualServersReq); err != nil {
		return err
	}

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
	createReq := models.CreateLBVirtualServers{}
	if err := tftags.Get(d, &createReq.CreateLBVirtualServersReq); err != nil {
		return err
	}

	//align createReq and fill json related fields
	if err := lb.virtualServerAlignRequest(ctx, meta, &createReq.CreateLBVirtualServersReq); err != nil {
		return err
	}

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

func (lb *loadBalancerVirtualServer) virtualServerAlignRequest(ctx context.Context, meta interface{},
	createReq *models.CreateLBVirtualServersReq) error {

	if createReq.TcpApplicationProfileConfig != nil && createReq.VipProtocol == "tcp" {
		createReq.VirtualServerConfig.ApplicationProfile = createReq.TcpApplicationProfileConfig.ApplicationProfile

	} else if createReq.UdpApplicationProfileConfig != nil && createReq.VipProtocol == "udp" {
		createReq.VirtualServerConfig.ApplicationProfile = createReq.UdpApplicationProfileConfig.ApplicationProfile

	} else if createReq.HttpApplicationProfileConfig != nil && createReq.VipProtocol == "http" {
		createReq.VirtualServerConfig.ApplicationProfile = createReq.HttpApplicationProfileConfig.ApplicationProfile

	}

	if createReq.CookiePersistenceProfileConfig != nil && createReq.Persistence == "COOKIE" {
		createReq.VirtualServerConfig.PersistenceProfile = createReq.CookiePersistenceProfileConfig.PersistenceProfile
		createReq.VirtualServerConfig.Persistence = createReq.Persistence
	} else if createReq.SourceipPersistenceProfileConfig != nil && createReq.Persistence == "SOURCE_IP" {
		createReq.VirtualServerConfig.Persistence = createReq.Persistence
		createReq.VirtualServerConfig.PersistenceProfile = createReq.SourceipPersistenceProfileConfig.PersistenceProfile

	}

	if createReq.SSLClientConfig != nil && createReq.SSLCert != 0 {
		createReq.VirtualServerConfig.SSLClientProfile = createReq.SSLClientConfig.SSLClientProfile
	}
	if createReq.SSLServerConfig != nil && createReq.SSLServerCert != 0 {
		createReq.VirtualServerConfig.SSLServerProfile = createReq.SSLServerConfig.SSLServerProfile
	}

	return nil
}
