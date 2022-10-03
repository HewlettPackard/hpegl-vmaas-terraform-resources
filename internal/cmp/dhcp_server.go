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

type dhcpServer struct {
	dhcpClient *client.DhcpServerAPIService
}

func newDhcpServer(dhcpServerClient *client.DhcpServerAPIService) *dhcpServer {
	return &dhcpServer{
		dhcpClient: dhcpServerClient,
	}
}

func (dhcp *dhcpServer) Read(ctx context.Context, d *utils.Data, meta interface{}) error {
	setMeta(meta, dhcp.dhcpClient.Client)
	var dhcpServerResp models.CreateNetworkDhcpServer
	if err := tftags.Get(d, &dhcpServerResp); err != nil {
		return err
	}
	getdhcpServerResp, err := dhcp.dhcpClient.GetSpecificDhcpServer(ctx, dhcpServerResp.ServerID,
		dhcpServerResp.ID)
	if err != nil {
		return err
	}

	return tftags.Set(d, getdhcpServerResp.GetSpecificNetworkDhcpServerResp)
}

func (dhcp *dhcpServer) Update(ctx context.Context, d *utils.Data, meta interface{}) error {
	id := d.GetID()
	var updateReq models.CreateNetworkDhcpServerRequest
	if err := tftags.Get(d, &updateReq.NetworkDhcpServer); err != nil {
		return err
	}

	retry := &utils.CustomRetry{
		InitialDelay: time.Second * 15,
		RetryDelay:   time.Second * 30,
	}
	_, err := retry.Retry(ctx, meta, func(ctx context.Context) (interface{}, error) {
		return dhcp.dhcpClient.UpdateDhcpServer(ctx, updateReq.NetworkDhcpServer.ServerID, id, updateReq)
	})
	if err != nil {
		return err
	}

	return tftags.Set(d, updateReq.NetworkDhcpServer)
}

func (dhcp *dhcpServer) Create(ctx context.Context, d *utils.Data, meta interface{}) error {
	setMeta(meta, dhcp.dhcpClient.Client)
	var createReq models.CreateNetworkDhcpServerRequest
	if err := tftags.Get(d, &createReq.NetworkDhcpServer); err != nil {
		return err
	}

	dhcpResp, err := dhcp.dhcpClient.CreateDhcpServer(ctx, createReq.NetworkDhcpServer.ServerID,
		createReq)
	if err != nil {
		return err
	}
	if !dhcpResp.Success {
		return fmt.Errorf(successErr, "creating dhcp")
	}
	createReq.NetworkDhcpServer.ID = dhcpResp.ID

	// wait until created
	retry := &utils.CustomRetry{
		InitialDelay: time.Second * 15,
		RetryDelay:   time.Second * 30,
	}
	_, err = retry.Retry(ctx, meta, func(ctx context.Context) (interface{}, error) {
		return dhcp.dhcpClient.GetSpecificDhcpServer(ctx,
			createReq.NetworkDhcpServer.ServerID, createReq.NetworkDhcpServer.ID)

	})
	if err != nil {
		return err
	}

	return tftags.Set(d, createReq.NetworkDhcpServer)
}

func (dhcp *dhcpServer) Delete(ctx context.Context, d *utils.Data, meta interface{}) error {
	setMeta(meta, dhcp.dhcpClient.Client)
	var dhcpServerResp models.CreateNetworkDhcpServer
	if err := tftags.Get(d, &dhcpServerResp); err != nil {
		return err
	}

	resp, err := dhcp.dhcpClient.DeleteDhcpServer(ctx, dhcpServerResp.ServerID, dhcpServerResp.ID)
	if err != nil {
		return err
	}

	if !resp.Success {
		return fmt.Errorf("got success = 'false' while deleting DHCP-SERVER")
	}

	return nil
}
