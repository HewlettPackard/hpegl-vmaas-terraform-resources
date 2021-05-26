// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"
	"errors"
	"strconv"

	"github.com/hpe-hcss/vmaas-cmp-go-sdk/pkg/client"
	"github.com/hpe-hcss/vmaas-terraform-resources/internal/logger"
	"github.com/hpe-hcss/vmaas-terraform-resources/internal/utils"
)

type group struct {
	gClient           *client.GroupsApiService
	serviceInstanceID string
}

func newGroup(gClient *client.GroupsApiService, serviceInstanceID string) *group {
	return &group{
		gClient:           gClient,
		serviceInstanceID: serviceInstanceID,
	}
}

func (g *group) Read(ctx context.Context, d *utils.Data) error {
	logger.Debug("Get Group")

	name := d.GetString("name")
	groups, err := g.gClient.GetAllGroups(ctx, g.serviceInstanceID, map[string]string{
		"name": name,
	})
	if err != nil {
		return err
	}

	if len(*groups.Groups) != 1 {
		return errors.New("error coudn't find exact group, please check the name")
	}

	d.SetID(strconv.Itoa((*groups.Groups)[0].Id))

	// post check
	if err := d.Error(); err != nil {
		return err
	}

	return nil
}
