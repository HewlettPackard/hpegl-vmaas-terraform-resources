// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package resources

const (
	// datasource key
	DSNetwork          = "hpegl_vmaas_network"
	DSLayout           = "hpegl_vmaas_layout"
	DSGroup            = "hpegl_vmaas_group"
	DSPlan             = "hpegl_vmaas_plan"
	DSCloud            = "hpegl_vmaas_cloud"
	DSResourcePool     = "hpegl_vmaas_resource_pool"
	DSDatastore        = "hpegl_vmaas_datastore"
	DSPowerSchedule    = "hpegl_vmaas_power_schedule"
	DSTemplate         = "hpegl_vmaas_template"
	DSEnvironment      = "hpegl_vmaas_environment"
	DSNetworkInterface = "hpegl_vmaas_network_interface"
	// resource key
	ResInstance = "hpegl_vmaas_instance"

	// documentation related constants
	generalNamedesc = "Name of the %s as it appears on GLPC Portal. " +
		"If no %s is found with this name 'NOT FOUND' error will returns."
	generalDDesc = "Unique ID to identify a %s."
)
