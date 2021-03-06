// // (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"
	"fmt"
	"log"

	"github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/client"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/utils"
)

type environment struct {
	eClient *client.EnvironmentAPIService
}

func newEnvironment(eClient *client.EnvironmentAPIService) *environment {
	return &environment{
		eClient: eClient,
	}
}

func (c *environment) Read(ctx context.Context, d *utils.Data, meta interface{}) error {
	setMeta(meta, c.eClient.Client)
	log.Printf("[DEBUG] Get Environments")

	name := d.GetString("name")
	if err := d.Error(); err != nil {
		return err
	}
	environment, err := c.eClient.GetAllEnvironment(ctx, map[string]string{
		nameKey: name,
	})
	if err != nil {
		return err
	}
	if len(environment.Environments) != 1 {
		return fmt.Errorf(errExactMatch, "environments")
	}
	d.SetString("code", environment.Environments[0].Code)
	d.SetID(environment.Environments[0].ID)

	// post check
	return d.Error()
}
