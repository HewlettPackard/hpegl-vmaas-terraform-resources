// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"
	"fmt"

	"github.com/hpe-hcss/vmaas-cmp-go-sdk/pkg/client"
	"github.com/hpe-hcss/vmaas-cmp-go-sdk/pkg/models"
	"github.com/hpe-hcss/vmaas-terraform-resources/internal/logger"
	"github.com/hpe-hcss/vmaas-terraform-resources/internal/utils"
	"github.com/hpe-hcss/vmaas-terraform-resources/pkg/auth"
)

type group struct {
	gClient *client.GroupsAPIService
}

func newGroup(gClient *client.GroupsAPIService) *group {
	return &group{
		gClient: gClient,
	}
}

func (g *group) Read(ctx context.Context, d *utils.Data, meta interface{}) error {
	logger.Debug("Get Group")

	name := d.GetString("name")
	// Pre check
	if err := d.Error(); err != nil {
		return err
	}
	resp, err := utils.Retry(func() (interface{}, error) {
		auth.SetScmClientToken(&ctx, meta)
		return g.gClient.GetAllGroups(ctx, nil)
	})
	groups := resp.(models.Groups)
	if err != nil {
		return err
	}
	isMatched := false
	for i, g := range *groups.Groups {
		if g.Name == name {
			isMatched = true
			d.SetID((*groups.Groups)[i].ID)

			break
		}
	}
	if !isMatched {
		return fmt.Errorf(errExactMatch, "group")
	}

	// post check
	return d.Error()
}
