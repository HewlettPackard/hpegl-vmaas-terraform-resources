// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"
	"errors"
	"log"

	"github.com/hpe-hcss/vmaas-cmp-go-sdk/pkg/client"
	"github.com/hpe-hcss/vmaas-cmp-go-sdk/pkg/models"
	"github.com/hpe-hcss/vmaas-terraform-resources/internal/utils"
)

type datastore struct {
	nClient *client.CloudsAPIService
}

func newDatastore(nClient *client.CloudsAPIService) *datastore {
	return &datastore{nClient: nClient}
}

func (n *datastore) Read(ctx context.Context, d *utils.Data, meta interface{}) error {
	log.Printf("[INFO] Get Datastore")

	// name := d.GetString("name")
	cloudID := d.GetInt("cloud_id")
	name := d.GetString("name")

	// Pre check
	if err := d.Error(); err != nil {
		return err
	}
	resp, err := utils.Retry(ctx, meta, func(ctx context.Context) (interface{}, error) {
		return n.nClient.GetAllCloudDataStores(ctx, cloudID,
			map[string]string{"name": name},
		)
	})
	if err != nil {
		return err
	}
	datastores := resp.(models.DataStoresResp)
	if len(datastores.Datastores) != 1 {
		return errors.New("error coudn't find exact datastore, please check the name")
	}
	d.SetID(datastores.Datastores[0].ID)

	// post check
	return d.Error()
}
