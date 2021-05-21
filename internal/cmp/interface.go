// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hpe-hcss/vmaas-cmp-go-sdk/pkg/models"
)

// CreateInstance CRUD operations for instances
type Instance interface {
	// CreateInstance will create instance and return nil if no error
	CreateInstance(ctx context.Context, d *schema.ResourceData, instanceBody models.CreateInstanceBody) error
	// GetInstance will fetch instance details as per ID
	GetInstance(ctx context.Context, d *schema.ResourceData, id int) (models.GetInstanceResponse, error)
}
