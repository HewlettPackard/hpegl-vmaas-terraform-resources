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

type loadBalancerds struct {
	lbClient *client.LoadBalancerAPIService
}

func newLoadBalancerDS(loadBalancerClient *client.LoadBalancerAPIService) *loadBalancerds {
	return &loadBalancerds{lbClient: loadBalancerClient}
}

func (n *loadBalancerds) Read(ctx context.Context, d *utils.Data, meta interface{}) error {
	setMeta(meta, n.lbClient.Client)
	log.Printf("[DEBUG] Get Load balancer")
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
			log.Print("[DEBUG]", lb.GetNetworkLoadBalancerResp[i].ID)

			return tftags.Set(d, lb.GetNetworkLoadBalancerResp[i])

		}
	}

	return fmt.Errorf(errExactMatch, "LoadBalancer")
}
