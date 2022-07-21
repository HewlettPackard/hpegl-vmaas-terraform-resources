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

type sslServerProfileds struct {
	lbServer *client.LoadBalancerAPIService
}

func newSslServerProfileDS(sslServerProfileServer *client.LoadBalancerAPIService) *sslServerProfileds {
	return &sslServerProfileds{lbServer: sslServerProfileServer}
}

func (n *sslServerProfileds) Read(ctx context.Context, d *utils.Data, meta interface{}) error {
	setMeta(meta, n.lbServer.Client)
	log.Printf("[DEBUG] Get SSLServer Profile")
	name := d.GetString("name")
	lbID := d.GetInt("lb_id")

	// Pre check
	if err := d.Error(); err != nil {
		return err
	}

	lb, err := n.lbServer.GetLBProfiles(ctx, lbID)
	if err != nil {
		return err
	}

	for i, n := range lb.GetLBProfilesResp {
		if n.Name == name {
			log.Print("[DEBUG]", lb.GetLBProfilesResp[i].ID)

			return tftags.Set(d, lb.GetLBProfilesResp[i])

		}
	}

	return fmt.Errorf(errExactMatch, "SSLServer Profile")
}
