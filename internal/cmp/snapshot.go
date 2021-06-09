// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hpe-hcss/vmaas-cmp-go-sdk/pkg/client"
	"github.com/hpe-hcss/vmaas-cmp-go-sdk/pkg/models"
	"github.com/hpe-hcss/vmaas-terraform-resources/internal/logger"
	"github.com/hpe-hcss/vmaas-terraform-resources/internal/utils"
)

// instance implements functions related to cmp instances
type snapshot struct {
	// expose Instance API service to instances related operations
	iClient *client.InstancesApiService
}

func newSnapshot(iClient *client.InstancesApiService) *instance {
	return &instance{
		iClient: iClient,
	}
}

// Create snapshot
func (s *snapshot) Create(ctx context.Context, d *utils.Data) error {
	logger.Debug("Creating VMware snapshot of instance")
	instanceID := d.GetInt("instance_id")
	req := &models.SnapshotBody{
		Snapshot: &models.SnapshotBodySnapshot{
			Name: d.GetString("name"),
			Description: d.GetString("Description"),
			},
	}

	// Pre check
	if err := d.Error(); err != nil {
		return err
	}
	// create snapshot
	resp, err := utils.Retry(func() (interface{}, error) {
		return s.iClient.SnapshotAnInstance(ctx,instanceID,req)
	})
	if err != nil {
		return err
	}
	snapshotResp:= resp.(models.SuccessOrErrorMessage)
	if !snapshotResp.Success {
		return fmt.Errorf("%s", snapshotResp.Message)
	}

	// post check
	return d.Error()
}


// Read snapshot and set state values accordingly
func (s *snapshot) Read(ctx context.Context, d *utils.Data) error {
	instanceId := d.GetInt("instance_id")

//	logger.Debug("Get snapshot with ID %d", id)

	// Precheck
	if err := d.Error(); err != nil {
		return err
	}

	resp, err := utils.Retry(func() (interface{}, error) {
		return s.iClient.GetListOfSnapshotsForAnInstance(ctx,instanceId)
	})
	if err != nil {
		return err
	}
	snapshots := resp.(models.ListSnapshotResponse)
	d.SetID(strconv.Itoa(snapshots.Snapshots[0].ID))
	d.SetString("status", snapshots.Snapshots[0].Status)

	// post check
	return d.Error()
}
