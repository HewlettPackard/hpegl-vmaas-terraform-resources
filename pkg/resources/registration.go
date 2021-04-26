// (C) Copyright 2020 Hewlett Packard Enterprise Development LP

package resources

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hpe-hcss/hpegl-provider-lib/pkg/registration"

	"github.com/hpe-hcss/poc-caas-terraform-resources/internal/resources"
)

// Assert that Registration implements the ServiceRegistration interface
var _ registration.ServiceRegistration = (*Registration)(nil)

type Registration struct{}

func (r Registration) Name() string {
	return "CAAS Service"
}

func (r Registration) SupportedDataSources() map[string]*schema.Resource {
	return nil
}

func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"hpegl_caas_cluster_blueprint": resources.ClusterBlueprint(),
		"hpegl_caas_cluster":           resources.Cluster(),
	}
}
