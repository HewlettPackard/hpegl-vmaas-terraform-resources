// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Resource interface implements all the resource operations (CRUD)
// All client resources expected to inherit and implement following
// functions.
type Resource interface {
	// Read terraform operations. Context and resource data as params.
	// will return error
	DataSource
	// Create terraform operations. Context and resource data as params.
	// will return error
	Create(context.Context, *schema.ResourceData) error
	// Update terraform operations. Context and resource data as params.
	// will return error
	Update(context.Context, *schema.ResourceData) error
	// Delete terraform operations. Context and resource data as params.
	// will return error
	Delete(context.Context, *schema.ResourceData) error
}

// DataSource interface wraps read operations which is expected to
// implement by all data source clients
type DataSource interface {
	Read(context.Context, *schema.ResourceData) error
}
