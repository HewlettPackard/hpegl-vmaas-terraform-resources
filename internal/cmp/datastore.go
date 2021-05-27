// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"
	"errors"
	"strconv"

	"github.com/hpe-hcss/vmaas-cmp-go-sdk/pkg/client"
	"github.com/hpe-hcss/vmaas-cmp-go-sdk/pkg/models"
	"github.com/hpe-hcss/vmaas-terraform-resources/internal/logger"
	"github.com/hpe-hcss/vmaas-terraform-resources/internal/utils"
)

type datastore struct {
	nClient           *client.CloudsApiService
	serviceInstanceID string
}

func newDatastore(nClient *client.CloudsApiService, serviceInstanceID string) *datastore {
	return &datastore{nClient: nClient, serviceInstanceID: serviceInstanceID}
}

func (n *datastore) Read(ctx context.Context, d *utils.Data) error {
	logger.Debug("Get Datastore")

	// name := d.GetString("name")
	cloudID := d.GetInt("cloud_id")
	name := d.GetString("name")

	// Pre check
	if err := d.Error(); err != nil {
		return err
	}
	resp, err := utils.Retry(func() (interface{}, error) {
		return n.nClient.GetAllCloudDataStores(ctx, n.serviceInstanceID, int(cloudID),
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
	d.SetID(strconv.Itoa(datastores.Datastores[0].ID))

	// post check
	return d.Error()
}
