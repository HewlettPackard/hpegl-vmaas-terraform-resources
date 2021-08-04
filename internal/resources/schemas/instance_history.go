package schemas

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func GetInstanceHistorySchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Description: `History details for the instance`,
		Computed:    true,
		Elem: schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:     schema.TypeInt,
					Computed: true,
				},
				"accountId": {
					Type:     schema.TypeInt,
					Computed: true,
				},
				"uniqueId": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"processType": {
					Type:     schema.TypeSet,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"code": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"name": {
								Type:     schema.TypeString,
								Computed: true,
							},
						},
					},
				},
				"displayName": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"instanceId": {
					Type:     schema.TypeInt,
					Computed: true,
				},
				"status": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"reason": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"percent": {
					Type:     schema.TypeFloat,
					Computed: true,
				},
				"statusEta": {
					Type:     schema.TypeFloat,
					Computed: true,
				},
				"startDate": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"endDate": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"duration": {
					Type:     schema.TypeInt,
					Computed: true,
				},
				"dateCreated": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"lastUpdated": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"createdBy": {
					Type: schema.TypeSet,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"username": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"displayName": {
								Type:     schema.TypeString,
								Computed: true,
							},
						},
					},
				},
				"updatedBy": {
					Type: schema.TypeSet,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"username": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"displayName": {
								Type:     schema.TypeString,
								Computed: true,
							},
						},
					},
				},
				"events": {
					Type:     schema.TypeList,
					Computed: true,
					Elem:     &schema.Resource{},
				},
			},
		},
	}
}
