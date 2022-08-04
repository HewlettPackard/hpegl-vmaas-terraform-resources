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

type lbSslVirtualServerCertds struct {
	lbClient *client.LoadBalancerAPIService
}

func newLBsslVirtualServerCertDS(sslVirtualServerCertClient *client.LoadBalancerAPIService) *lbSslVirtualServerCertds {
	return &lbSslVirtualServerCertds{lbClient: sslVirtualServerCertClient}
}

func (n *lbSslVirtualServerCertds) Read(ctx context.Context, d *utils.Data, meta interface{}) error {
	setMeta(meta, n.lbClient.Client)
	log.Printf("[DEBUG] Get SSL Certs")
	name := d.GetString("name")

	// Pre check
	if err := d.Error(); err != nil {
		return err
	}

	lb, err := n.lbClient.GetLBVirtualServerSSLCerts(ctx)
	if err != nil {
		return err
	}

	for i, n := range lb.Certificates {
		if n.Name == name {
			log.Print("[DEBUG]", lb.Certificates[i].ID)

			return tftags.Set(d, lb.Certificates[i])
		}
	}

	return fmt.Errorf(errExactMatch, "SSL Certificates")
}
