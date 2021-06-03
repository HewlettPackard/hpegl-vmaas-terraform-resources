// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package resources

const (
	// datasource key
	DSNetwork       = "hpegl_vmaas_network"
	DSLayout        = "hpegl_vmaas_layout"
	DSGroup         = "hpegl_vmaas_group"
	DSPlan          = "hpegl_vmaas_plan"
	DSCloud         = "hpegl_vmaas_cloud"
	DSResourcePool  = "hpegl_vmaas_resourcePool"
	DSDatastore     = "hpegl_vmaas_datastore"
	DSPowerSchedule = "hpegl_vmaas_powerSchedule"
	// resource key
	ResInstance = "hpegl_vmaas_instance"

	// documentation related constants
	generalNamedesc = "Name of the %s as it appears on GLPC Portal. " +
		"If no %s is found with this name standard not found error returns will return."
	generalDDesc = "Unique ID to identify a %s."
)
