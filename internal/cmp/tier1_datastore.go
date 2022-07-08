// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"
	"fmt"
	"log"

	"github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/client"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/utils"
	"github.com/tshihad/tftags"
)

type lbTier1ds struct {
	lbClient *client.LoadBalancerAPIService
}

func newlbTier1dsDS(loadBalancerClient *client.LoadBalancerAPIService) *lbTier1ds {
	return &lbTier1ds{lbClient: loadBalancerClient}
}

func (n *lbTier1ds) Read(ctx context.Context, d *utils.Data, meta interface{}) error {
	setMeta(meta, n.lbClient.Client)
	log.Printf("[DEBUG] Get Load balancer tier1")
	name := d.GetString("name")

	// Pre check
	if err := d.Error(); err != nil {
		return err
	}

	lb, err := n.lbClient.GetLoadBalancers(ctx)
	if err != nil {
		return err
	}

	for i, n := range lb.GetNetworkLoadBalancerResp {
		if n.Name == name {
			log.Print("[DEBUG]", lb.GetNetworkLoadBalancerResp[i].Config.Tier1)

			return tftags.Set(d, lb.GetNetworkLoadBalancerResp[i].Config.Tier1)

		}
	}

	return fmt.Errorf(errExactMatch, "LoadBalancer Tier1")
}
