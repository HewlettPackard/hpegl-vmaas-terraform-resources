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

type lb_poolds struct {
	lbClient *client.LoadBalancerAPIService
}

func newLBPoolDS(poolClient *client.LoadBalancerAPIService) *lb_poolds {
	return &lb_poolds{lbClient: poolClient}
}

func (n *lb_poolds) Read(ctx context.Context, d *utils.Data, meta interface{}) error {
	setMeta(meta, n.lbClient.Client)
	log.Printf("[DEBUG] Get Pool")
	name := d.GetString("name")
	lbID := d.GetInt("lb_id")

	// Pre check
	if err := d.Error(); err != nil {
		return err
	}

	lb, err := n.lbClient.GetLBPools(ctx, lbID)
	if err != nil {
		return err
	}

	for i, n := range lb.GetLBPoolsResp {
		if n.Name == name {
			log.Print("[DEBUG]", lb.GetLBPoolsResp[i].ID)

			return tftags.Set(d, lb.GetLBPoolsResp[i])

		}
	}

	return fmt.Errorf(errExactMatch, "Pool")
}
