// (C) Copyright 2020 Hewlett Packard Enterprise Development LP

package resources

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hpe-hcss/hpegl-provider-lib/pkg/registration"

	"github.com/hpe-hcss/vmaas-terraform-resources/internal/resources"
)

// Assert that Registration implements the ServiceRegistration interface
var _ registration.ServiceRegistration = (*Registration)(nil)

type Registration struct{}

func (r Registration) Name() string {
	return "VMAAS Service"
}

func (r Registration) SupportedDataSources() map[string]*schema.Resource {
	return nil
}

func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"hpegl_vmaas_cluster_blueprint": resources.ClusterBlueprint(),
		"hpegl_vmaas_cluster":           resources.Cluster(),
	}
}
