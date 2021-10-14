// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package resources

const (
	// datasource key
	DSNetwork          = "hpegl_vmaas_network"
	DSNetworkType      = "hpegl_vmaas_network_type"
	DSNetworkPool      = "hpegl_vmaas_network_pool"
	DSNetworkProxy     = "hpegl_vmaas_network_proxy"
	DSNetworkDomain    = "hpegl_vmaas_network_domain"
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
	DSCloudFolder      = "hpegl_vmaas_cloud_folder"
	DSRouter           = "hpegl_vmaas_router"
	// resource key
	ResInstance      = "hpegl_vmaas_instance"
	ResInstanceClone = "hpegl_vmaas_instance_clone"
	ResNetwork       = "hpegl_vmaas_network"
	ResRouter        = "hpegl_vmaas_router"

	// documentation related constants
	generalNamedesc = "Name of the %s as it appears on GLPC Portal. " +
		"If there is no %s with this name, a 'NOT FOUND' error will returned."
	generalDDesc = "Unique ID to identify a %s."
)
