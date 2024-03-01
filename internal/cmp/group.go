// (C) Copyright 2021-2024 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"
	"fmt"
	"log"

	"github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/client"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/utils"
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
	setMeta(meta, g.gClient.Client)
	log.Printf("[DEBUG] Get Group")

	name := d.GetString("name")
	// Pre check
	if err := d.Error(); err != nil {
		return err
	}
	param := map[string]string{maxKey: "-1"}
	groups, err := g.gClient.GetAllGroups(ctx, param)
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
