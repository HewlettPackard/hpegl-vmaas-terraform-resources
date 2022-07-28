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
	if err := tftags.Get(d, &updateReq.CreateLBVirtualServersReq); err != nil {
		return err
	}

	// align createReq and fill json related fields
	if err := lb.virtualServerAlignRequest(ctx, meta, &updateReq); err != nil {
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
	var createReq models.CreateLBVirtualServers
	if err := tftags.Get(d, &createReq.CreateLBVirtualServersReq); err != nil {
		return err
	}

	// align createReq and fill json related fields
	if err := lb.virtualServerAlignRequest(ctx, meta, &createReq); err != nil {
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

func (lb *loadBalancerVirtualServer) virtualServerAlignRequest(ctx context.Context, meta interface{}, createReq *models.CreateLBVirtualServers) error {

	if createReq.CreateLBVirtualServersReq.TcpApplicationProfileConfig != nil &&
		createReq.CreateLBVirtualServersReq.VipProtocol == "tcp" {
		createReq.CreateLBVirtualServersReq.VirtualServerConfig.ApplicationProfile =
			createReq.CreateLBVirtualServersReq.TcpApplicationProfileConfig.ApplicationProfile

	}
	if createReq.CreateLBVirtualServersReq.UdpApplicationProfileConfig != nil &&
		createReq.CreateLBVirtualServersReq.VipProtocol == "udp" {
		createReq.CreateLBVirtualServersReq.VirtualServerConfig.ApplicationProfile =
			createReq.CreateLBVirtualServersReq.UdpApplicationProfileConfig.ApplicationProfile

	}
	if createReq.CreateLBVirtualServersReq.HttpApplicationProfileConfig != nil &&
		createReq.CreateLBVirtualServersReq.VipProtocol == "http" {
		createReq.CreateLBVirtualServersReq.VirtualServerConfig.ApplicationProfile =
			createReq.CreateLBVirtualServersReq.HttpApplicationProfileConfig.ApplicationProfile

	}

	if createReq.CreateLBVirtualServersReq.CookiePersistenceProfileConfig != nil &&
		createReq.CreateLBVirtualServersReq.Persistence == "COOKIE" {

		createReq.CreateLBVirtualServersReq.VirtualServerConfig.PersistenceProfile =
			createReq.CreateLBVirtualServersReq.CookiePersistenceProfileConfig.PersistenceProfile
	}
	if createReq.CreateLBVirtualServersReq.SourceipPersistenceProfileConfig != nil &&
		createReq.CreateLBVirtualServersReq.Persistence == "SOURCE_IP" {
		createReq.CreateLBVirtualServersReq.VirtualServerConfig.PersistenceProfile =
			createReq.CreateLBVirtualServersReq.SourceipPersistenceProfileConfig.PersistenceProfile

	}

	if createReq.CreateLBVirtualServersReq.SSLClientConfig != nil {
		createReq.CreateLBVirtualServersReq.VirtualServerConfig.SSLClientProfile =
			createReq.CreateLBVirtualServersReq.SSLClientConfig.SSLClientProfile
	}
	if createReq.CreateLBVirtualServersReq.SSLServerConfig != nil {
		createReq.CreateLBVirtualServersReq.VirtualServerConfig.SSLServerProfile =
			createReq.CreateLBVirtualServersReq.SSLServerConfig.SSLServerProfile
	}

	return nil
}
