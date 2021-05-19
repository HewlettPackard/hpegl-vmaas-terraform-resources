// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package cmp

import (
	client "github.com/hpe-hcss/vmaas-cmp-go-sdk/pkg/client"
	"github.com/hpe-hcss/vmaas-cmp-go-sdk/pkg/models"
)

// instance implements functions related to cmp instances
type instance struct {
	client *client.APIClient
}

// CreateInstance create instance
func (i *instance) CreateInstance(instanceBody models.CreateInstanceBody) error {
	return nil
}
