package schemas

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func GetInstanceHistorySchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Description: `History details for the instance`,
		Computed:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id":         computedInt(),
				"account_id": computedInt(),
				"unique_id":  computedString(),
				"process_type": {
					Type:     schema.TypeSet,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"code": computedString(),
							"name": computedString(),
						},
					},
				},
				"display_name": computedString(),
				"instance_id":  computedInt(),
				"status":       computedString(),
				"reason":       computedString(),
				"percent": {
					Type:     schema.TypeFloat,
					Computed: true,
				},
				"status_eta": {
					Type:     schema.TypeFloat,
					Computed: true,
				},
				"start_date":   computedString(),
				"end_date":     computedString(),
				"duration":     computedInt(),
				"date_created": computedString(),
				"last_updated": computedString(),
				"created_by":   computedUpdatedBySchema(),
				"updated_by":   computedUpdatedBySchema(),
			},
		},
	}
}

func GetInstanceContainerSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Description: `Container's details for the instance which contains IP addresses, hostname and other stats`,
		Computed:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id":             computedInt(),
				"name":           computedString(),
				"ip":             computedString(),
				"external_fqdn":  computedString(),
				"container_type": computedSingleStringSet("name"),
				"server": {
					Type:     schema.TypeSet,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"id":    computedInt(),
							"owner": computedSingleStringSet("username"),
							"compute_server_type": {
								Type:     schema.TypeSet,
								Computed: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"name":            computedString(),
										"managed":         computedBool(),
										"external_delete": computedBool(),
									},
								},
							},
							"visibility":       computedString(),
							"ssh_host":         computedString(),
							"ssh_port":         computedInt(),
							"platform":         computedString(),
							"platform_version": computedString(),
							"date_created":     computedString(),
							"last_updated":     computedString(),
							"server_os":        computedSingleStringSet("name"),
						},
					},
				},
				"hostname":    computedString(),
				"max_storage": computedInt(),
				"max_memory":  computedInt(),
				"max_cores":   computedInt(),
			},
		},
	}
}

func computedString() *schema.Schema {
	return &schema.Schema{
		Computed: true,
		Type:     schema.TypeString,
	}
}

func computedInt() *schema.Schema {
	return &schema.Schema{
		Computed: true,
		Type:     schema.TypeInt,
	}
}

func computedBool() *schema.Schema {
	return &schema.Schema{
		Computed: true,
		Type:     schema.TypeBool,
	}
}

func computedSingleStringSet(key string) *schema.Schema {
	return &schema.Schema{
		Computed: true,
		Type:     schema.TypeSet,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				key: computedString(),
			},
		},
	}
}

func computedUpdatedBySchema() *schema.Schema {
	return &schema.Schema{
		Computed: true,
		Type:     schema.TypeSet,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"username":     computedString(),
				"display_name": computedString(),
			},
		},
	}
}
