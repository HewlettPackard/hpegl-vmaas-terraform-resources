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

type passiveMonitords struct {
	lbClient *client.LoadBalancerAPIService
}

func newPassiveMonitorDS(passiveMonitorClient *client.LoadBalancerAPIService) *passiveMonitords {
	return &passiveMonitords{lbClient: passiveMonitorClient}
}

func (n *passiveMonitords) Read(ctx context.Context, d *utils.Data, meta interface{}) error {
	setMeta(meta, n.lbClient.Client)
	log.Printf("[DEBUG] Get Passive Monitors")
	monitorType := d.GetString("type")
	lbID := d.GetInt("lb_id")

	// Pre check
	if err := d.Error(); err != nil {
		return err
	}

	lb, err := n.lbClient.GetLBMonitors(ctx, lbID)
	if err != nil {
		return err
	}

	for i, n := range lb.GetLBMonitorsResp {
		if n.Type == monitorType {
			log.Print("[DEBUG]", lb.GetLBMonitorsResp[i].ID)

			return tftags.Set(d, lb.GetLBMonitorsResp[i])

		}
	}
	return fmt.Errorf(errExactMatch, "Passive Monitors")
}
