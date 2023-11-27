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

type poolMemberGroupds struct {
	lbClient *client.LoadBalancerAPIService
	rClient  *client.RouterAPIService
}

func newLBPoolMemberGroupDS(loadBalancerClient *client.LoadBalancerAPIService,
	routerClient *client.RouterAPIService) *poolMemberGroupds {
	return &poolMemberGroupds{
		lbClient: loadBalancerClient,
		rClient:  routerClient,
	}
}

func (n *poolMemberGroupds) Read(ctx context.Context, d *utils.Data, meta interface{}) error {
	setMeta(meta, n.lbClient.Client)
	log.Printf("[DEBUG] Get Pool Member Group")
	name := d.GetString("name")

	// Pre check
	if err := d.Error(); err != nil {
		return err
	}
	nsxType, err := GetNsxTypeFromCMP(ctx, n.rClient.Client)
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

	lb, err := n.lbClient.GetLBPoolMemberGroup(ctx, serverID)
	if err != nil {
		return err
	}

	for i, n := range lb.MemeberGroup {
		if n.Name == name {
			log.Print("[DEBUG]", lb.MemeberGroup[i].ExternalID)

			return tftags.Set(d, lb.MemeberGroup[i])
		}
	}

	return fmt.Errorf(errExactMatch, "Pool Member Group")
}
