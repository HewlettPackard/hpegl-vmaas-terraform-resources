package schemas

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

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
