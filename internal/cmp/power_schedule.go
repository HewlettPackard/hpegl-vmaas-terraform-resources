// // (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"
	"fmt"
	"log"

	"github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/client"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/utils"
)

type powerSchedule struct {
	powerScheduleClient *client.PowerSchedulesAPIService
}

func newPowerSchedule(powerScheduleClient *client.PowerSchedulesAPIService) *powerSchedule {
	return &powerSchedule{
		powerScheduleClient: powerScheduleClient,
	}
}

func (c *powerSchedule) Read(ctx context.Context, d *utils.Data, meta interface{}) error {
	log.Printf("[DEBUG] Get Power Schedule")

	name := d.GetString("name")
	if err := d.Error(); err != nil {
		return err
	}
	powerSchedule, err := c.powerScheduleClient.GetAllPowerSchedules(ctx, map[string]string{
		nameKey: name,
	})
	if err != nil {
		return err
	}
	if len(powerSchedule.Schedules) != 1 {
		return fmt.Errorf(errExactMatch, "powerSchedules")
	}
	d.SetID(powerSchedule.Schedules[0].ID)

	// post check
	return d.Error()
}
