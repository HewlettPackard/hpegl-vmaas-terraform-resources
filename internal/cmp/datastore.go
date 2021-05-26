// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"
	"io/ioutil"

	"github.com/hpe-hcss/vmaas-cmp-go-sdk/pkg/client"
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
	res, err := n.nClient.GetAllCloudDataStores(ctx, n.serviceInstanceID, int(cloudID), nil)
	if err != nil {
		return err
	}
	body, _ := ioutil.ReadAll(res.Body)
	logger.Debug(string(body))
	// if len(datastores.Datastores) != 1 {
	// 	return errors.New("error coudn't find exact datastore, please check the name")
	// }
	// d.SetID(strconv.Itoa(datastores.Datastores[0].Id))

	// // post check
	// if err := d.Error(); err != nil {
	// 	return err
	// }

	return nil
}
