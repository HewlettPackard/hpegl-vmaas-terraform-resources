// (C) Copyright 2022 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"
	"fmt"
	"log"

	"github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/client"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/utils"
	"github.com/tshihad/tftags"
)

type dhcpServerds struct {
	dhcpClient *client.DhcpServerAPIService
	rClient    *client.RouterAPIService
}

func newDHCPServerDS(dhcpServerClient *client.DhcpServerAPIService,
	routerClient *client.RouterAPIService) *dhcpServerds {
	return &dhcpServerds{
		dhcpClient: dhcpServerClient,
		rClient:    routerClient,
	}
}

func (n *dhcpServerds) Read(ctx context.Context, d *utils.Data, meta interface{}) error {
	setMeta(meta, n.dhcpClient.Client)
	log.Printf("[DEBUG] Get DHCP Server")
	name := d.GetString("name")

	// Pre check
	if err := d.Error(); err != nil {
		return err
	}
	nsxType, err := GetNsxTypeFromCMP(n.rClient.Client)
	if err != nil {
		return err
	}
	setMeta(meta, n.rClient.Client)
	// Get network server ID for nsx-t
	serverResp, err := n.rClient.GetNetworkServices(ctx, nil)
	if err != nil {
		return err
	}

	var serverID int
	for i, n := range serverResp.NetworkServices {
		if n.TypeName == nsxType {
			serverID = serverResp.NetworkServices[i].ID

			break
		}
	}

	if serverID == 0 {
		return fmt.Errorf(errExactMatch, "network server")
	}

	dh, err := n.dhcpClient.GetDhcpServers(ctx, serverID)
	if err != nil {
		return err
	}

	for i, n := range dh.GetNetworkDhcpServerResp {
		if n.Name == name {
			log.Print("[DEBUG]", dh.GetNetworkDhcpServerResp[i].ProviderID)

			return tftags.Set(d, dh.GetNetworkDhcpServerResp[i])
		}
	}

	return fmt.Errorf(errExactMatch, "DHCP Server")
}
