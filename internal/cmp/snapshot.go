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

// snapshot implements functions related to cmp instance snapshot
type snapshot struct {
	// expose Instance API service to instances related operations
	sClient *client.InstancesApiService
}

func newSnapshot(sClient *client.InstancesApiService) *snapshot {
	return &snapshot{
		sClient: sClient,
	}
}

// Create snapshot
func (s *snapshot) Create(ctx context.Context, d *utils.Data) error {
	logger.Debug("Creating VMware snapshot of instance")
	instanceID := d.GetInt("instance_id")
	req := &models.SnapshotBody{
		Snapshot: &models.SnapshotBodySnapshot{
			Name:        d.GetString("name"),
			Description: d.GetString("description"),
		},
	}

	// Pre check
	if err := d.Error(); err != nil {
		return err
	}
	// create snapshot
	resp, err := utils.Retry(func() (interface{}, error) {
		return s.sClient.SnapshotAnInstance(ctx, instanceID, req)
	})
	if err != nil {
		return err
	}
	snapshotResp := resp.(models.Instances)
	if !snapshotResp.Success {
		return fmt.Errorf("Failed to create snapshot, Please try again",
		 "If issue persists contact your administrator")
	}

	// post check
	return d.Error()
}

// Read snapshot and set state values accordingly
func (s *snapshot) Read(ctx context.Context, d *utils.Data) error {
	instanceID := d.GetInt("instance_id")

	// Precheck
	if err := d.Error(); err != nil {
		return err
	}

	resp, err := utils.Retry(func() (interface{}, error) {
		return s.sClient.GetListOfSnapshotsForAnInstance(ctx, instanceID)
	})
	if err != nil {
		return err
	}
	snapshots := resp.(models.ListSnapshotResponse)
	d.SetID(strconv.Itoa(snapshots.Snapshots[0].ID))
	d.SetString("status", snapshots.Snapshots[0].Status)
	d.SetString("timestamp", snapshots.Snapshots[0].DateCreated.String())

	// post check
	return d.Error()
}

func (s *snapshot) Delete(ctx context.Context, d *utils.Data) error {
	return nil
}

func (s *snapshot) Update(ctx context.Context, d *utils.Data) error {
	return nil
}
