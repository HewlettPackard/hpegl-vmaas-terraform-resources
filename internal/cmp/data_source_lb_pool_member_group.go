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
}

func newLBPoolMemberGroupDS(poolMemberGroupClient *client.LoadBalancerAPIService) *poolMemberGroupds {
	return &poolMemberGroupds{lbClient: poolMemberGroupClient}
}

func (n *poolMemberGroupds) Read(ctx context.Context, d *utils.Data, meta interface{}) error {
	setMeta(meta, n.lbClient.Client)
	log.Printf("[DEBUG] Get Pool Member Group")
	name := d.GetString("name")

	// Pre check
	if err := d.Error(); err != nil {
		return err
	}

	lb, err := n.lbClient.GetLBPoolMemberGroup(ctx)
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
