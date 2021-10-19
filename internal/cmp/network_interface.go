// // (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/client"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/utils"
)

type networkInterface struct {
	cClient *client.CloudsAPIService
	pClient *client.ProvisioningAPIService
}

func newNetworkInterface(cClient *client.CloudsAPIService, pClient *client.ProvisioningAPIService) *networkInterface {
	return &networkInterface{
		cClient: cClient,
		pClient: pClient,
	}
}

func (c *networkInterface) Read(ctx context.Context, d *utils.Data, meta interface{}) error {
	log.Printf("[DEBUG] Get Network interface")

	cloudID := d.GetInt("cloud_id")
	name := d.GetString("name")
	if err := d.Error(); err != nil {
		return err
	}

	// Get vmware provision-type id
	provision, err := c.pClient.GetAllProvisioningTypes(ctx, map[string]string{
		nameKey: vmware,
	})
	if err != nil {
		return err
	}
	if len(provision.ProvisionTypes) != 1 {
		return errors.New("could not find vmware provision type. Please contact administrator to resolve the issue")
	}

	networkInterface, err := c.cClient.GetAllCloudNetworks(ctx, cloudID, provision.ProvisionTypes[0].ID)
	if err != nil {
		return err
	}

	for _, n := range networkInterface.Data.NetworkTypes {
		if n.Name == name {
			d.Set("code", n.Code)
			d.SetID(n.ID)

			return d.Error()
		}
	}

	return fmt.Errorf(errExactMatch, "network-interface")
}
