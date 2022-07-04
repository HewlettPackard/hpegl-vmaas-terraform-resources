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

type loadBalancerpoolds struct {
	lbClient *client.LoadBalancerAPIService
}

func newLoadBalancerPoolDS(loadBalancerClient *client.LoadBalancerAPIService) *loadBalancerpoolds {
	return &loadBalancerpoolds{lbClient: loadBalancerClient}
}

func (p *loadBalancerpoolds) Read(ctx context.Context, d *utils.Data, meta interface{}) error {
	setMeta(meta, p.lbClient.Client)
	log.Printf("[DEBUG] Get Load balancer")
	name := d.GetString("name")

	// Pre check
	if err := d.Error(); err != nil {
		return err
	}

	lb, err := p.lbClient.GetSpecificLoadBalancers(ctx, 186)
	if err != nil {
		return err
	}

	// var lbPools models.GetLBPools
	// for i, _ := range lb.GetNetworkLoadBalancerResp {
	// 	lbPools, err = p.lbClient.GetLBPools(ctx, lb.GetNetworkLoadBalancerResp[i].ID)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	//var lbPools models.GetLBPools
	//for i, _ := range lb.GetSpecificNetworkLoadBalancerResp {
	lbPools, err := p.lbClient.GetLBPools(ctx, lb.GetSpecificNetworkLoadBalancerResp.ID)
	if err != nil {
		return err
	}
	//	}

	for i, n := range lbPools.GetLBPoolsResp {
		if n.Name == name {
			log.Print("[DEBUG]", lbPools.GetLBPoolsResp[i].ID)

			return tftags.Set(d, lbPools.GetLBPoolsResp[i])
		}
	}

	return fmt.Errorf(errExactMatch, "LoadBalancer Pool")
}
